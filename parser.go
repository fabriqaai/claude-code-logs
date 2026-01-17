package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultClaudeProjectsPath returns the default path to Claude projects
func DefaultClaudeProjectsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "projects"), nil
}

// DecodeProjectPath converts an encoded folder name to a human-readable path
// Example: "-Users-name-project" -> "/Users/name/project"
func DecodeProjectPath(folderName string) string {
	if folderName == "" {
		return ""
	}

	// Handle special case of just "-" (root marker)
	if folderName == "-" {
		return "/"
	}

	// Remove leading dash and split by dash
	if !strings.HasPrefix(folderName, "-") {
		return folderName // Not an encoded path
	}

	// Remove leading dash
	path := folderName[1:]

	// Replace dashes with slashes
	path = strings.ReplaceAll(path, "-", "/")

	// Prepend slash
	return "/" + path
}

// DiscoverProjects scans the Claude projects directory and returns all projects
func DiscoverProjects(projectsPath string) ([]Project, error) {
	// Check if directory exists
	info, err := os.Stat(projectsPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("claude projects directory not found: %s", projectsPath)
	}
	if err != nil {
		return nil, fmt.Errorf("accessing projects directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("projects path is not a directory: %s", projectsPath)
	}

	// Read directory entries
	entries, err := os.ReadDir(projectsPath)
	if err != nil {
		return nil, fmt.Errorf("reading projects directory: %w", err)
	}

	var projects []Project
	for _, entry := range entries {
		// Skip non-directories and hidden files
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip entries that don't look like encoded paths
		if !strings.HasPrefix(entry.Name(), "-") {
			continue
		}

		project := Project{
			FolderName: entry.Name(),
			Path:       DecodeProjectPath(entry.Name()),
			Sessions:   nil, // Sessions are loaded lazily
		}
		projects = append(projects, project)
	}

	// Sort projects by path for consistent ordering
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Path < projects[j].Path
	})

	return projects, nil
}

// ListSessions lists all sessions for a project
func ListSessions(projectsPath string, project *Project) ([]Session, error) {
	projectDir := filepath.Join(projectsPath, project.FolderName)

	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, fmt.Errorf("reading project directory %s: %w", project.Path, err)
	}

	var sessions []Session
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		// Extract session ID from filename (UUID.jsonl)
		sessionID := strings.TrimSuffix(entry.Name(), ".jsonl")

		// Parse session to get metadata
		sessionPath := filepath.Join(projectDir, entry.Name())
		session, err := ParseSession(sessionPath, sessionID)
		if err != nil {
			// Log warning but continue with other sessions
			fmt.Fprintf(os.Stderr, "Warning: failed to parse session %s: %v\n", sessionID, err)
			continue
		}

		sessions = append(sessions, *session)
	}

	// Sort sessions by creation date (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}

// ParseSession parses a JSONL session file into a Session struct
func ParseSession(filePath string, sessionID string) (*Session, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening session file: %w", err)
	}
	defer file.Close()

	session := &Session{
		ID:         sessionID,
		Messages:   []Message{},
		SourcePath: filePath,
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer size for very long lines (some Claude sessions have 20MB+ lines)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 50*1024*1024) // 50MB max line size

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		// Parse the JSONL entry
		var entry jsonlEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping malformed JSON at line %d in %s: %v\n",
				lineNum, sessionID, err)
			continue
		}

		// Handle different entry types
		switch entry.Type {
		case "summary":
			session.Summary = entry.Summary

		case "user", "assistant":
			msg, err := parseMessage(&entry)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse message at line %d in %s: %v\n",
					lineNum, sessionID, err)
				continue
			}
			session.Messages = append(session.Messages, *msg)

			// Capture CWD from first entry that has it (actual project path)
			if session.CWD == "" && entry.CWD != "" {
				session.CWD = entry.CWD
			}

			// Update session timestamps
			if session.CreatedAt.IsZero() || msg.Timestamp.Before(session.CreatedAt) {
				session.CreatedAt = msg.Timestamp
			}
			if msg.Timestamp.After(session.UpdatedAt) {
				session.UpdatedAt = msg.Timestamp
			}

		case "file-history-snapshot":
			// Skip these entries
			continue

		default:
			// Unknown type, skip silently
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading session file: %w", err)
	}

	// If no summary, use first user message as fallback
	if session.Summary == "" && len(session.Messages) > 0 {
		for _, msg := range session.Messages {
			if msg.Role == "user" {
				for _, block := range msg.Content {
					if block.Type == "text" && block.Text != "" {
						// Truncate to first 100 chars
						summary := block.Text
						if len(summary) > 100 {
							summary = summary[:97] + "..."
						}
						session.Summary = summary
						break
					}
				}
				break
			}
		}
	}

	return session, nil
}

// parseMessage converts a JSONL entry to a Message struct
func parseMessage(entry *jsonlEntry) (*Message, error) {
	msg := &Message{
		UUID:    entry.UUID,
		Content: []ContentBlock{},
	}

	// Handle parentUuid (can be null)
	if entry.ParentUUID != nil {
		msg.ParentUUID = *entry.ParentUUID
	}

	// Parse timestamp
	if entry.Timestamp != "" {
		t, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			// Try alternate format
			t, err = time.Parse("2006-01-02T15:04:05.000Z", entry.Timestamp)
			if err != nil {
				msg.Timestamp = time.Now() // Fallback
			} else {
				msg.Timestamp = t
			}
		} else {
			msg.Timestamp = t
		}
	}

	// Parse message content
	if len(entry.Message) == 0 {
		return msg, nil
	}

	var msgContent messageContent
	if err := json.Unmarshal(entry.Message, &msgContent); err != nil {
		return nil, fmt.Errorf("parsing message content: %w", err)
	}

	msg.Role = msgContent.Role

	// Parse content blocks
	// Content can be a string or an array of blocks
	if len(msgContent.Content) == 0 {
		return msg, nil
	}

	// Try to parse as array of blocks first
	var blocks []rawContentBlock
	if err := json.Unmarshal(msgContent.Content, &blocks); err != nil {
		// Try to parse as a simple string
		var textContent string
		if err := json.Unmarshal(msgContent.Content, &textContent); err != nil {
			return nil, fmt.Errorf("parsing content blocks: %w", err)
		}
		msg.Content = []ContentBlock{{Type: "text", Text: textContent}}
		return msg, nil
	}

	// Parse each content block
	for _, block := range blocks {
		cb := ContentBlock{Type: block.Type}

		switch block.Type {
		case "text":
			cb.Text = block.Text

		case "tool_use":
			cb.ToolName = block.Name
			cb.ToolUseID = block.ID
			if len(block.Input) > 0 {
				cb.ToolInput = string(block.Input)
			}

		case "tool_result":
			cb.ToolUseID = block.ToolUseID
			// Content can be a string or complex object
			if len(block.Content) > 0 {
				var textContent string
				if err := json.Unmarshal(block.Content, &textContent); err != nil {
					// Store as JSON string if not a simple string
					cb.ToolOutput = string(block.Content)
				} else {
					cb.ToolOutput = textContent
				}
			}

		default:
			// Unknown block type, store type for debugging
			cb.Text = fmt.Sprintf("[unknown block type: %s]", block.Type)
		}

		msg.Content = append(msg.Content, cb)
	}

	return msg, nil
}

// LoadProjectWithSessions loads a project and all its sessions
func LoadProjectWithSessions(projectsPath string, project *Project) error {
	sessions, err := ListSessions(projectsPath, project)
	if err != nil {
		return err
	}
	project.Sessions = sessions

	// Use actual CWD from sessions if available (fixes path decoding ambiguity)
	// Sessions are sorted newest first, so use the first session's CWD
	for _, s := range sessions {
		if s.CWD != "" {
			project.Path = s.CWD
			break
		}
	}

	return nil
}

// LoadAllProjects discovers and loads all projects with their sessions
func LoadAllProjects(projectsPath string) ([]Project, error) {
	projects, err := DiscoverProjects(projectsPath)
	if err != nil {
		return nil, err
	}

	for i := range projects {
		if err := LoadProjectWithSessions(projectsPath, &projects[i]); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load sessions for %s: %v\n",
				projects[i].Path, err)
			// Continue with other projects
		}
	}

	// Sort projects by last update date (most recent first)
	sort.Slice(projects, func(i, j int) bool {
		return getProjectLastUpdate(&projects[i]).After(getProjectLastUpdate(&projects[j]))
	})

	return projects, nil
}

// getProjectLastUpdate returns the most recent UpdatedAt time from a project's sessions
func getProjectLastUpdate(project *Project) time.Time {
	var lastUpdate time.Time
	for _, session := range project.Sessions {
		if session.UpdatedAt.After(lastUpdate) {
			lastUpdate = session.UpdatedAt
		}
	}
	return lastUpdate
}
