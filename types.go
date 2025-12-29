package main

import (
	"encoding/json"
	"time"
)

// Project represents a Claude Code project folder
type Project struct {
	Path       string    // Decoded path: /Users/name/project
	FolderName string    // Encoded folder name: -Users-name-project
	Sessions   []Session // Sessions in this project
}

// Session represents a single chat session
type Session struct {
	ID         string    // UUID from filename
	Summary    string    // Session title from summary entry
	Messages   []Message // All messages in the session
	CreatedAt  time.Time // First message timestamp
	UpdatedAt  time.Time // Last message timestamp
	SourcePath string    // Full path to source JSONL file
}

// Message represents a single message in a session
type Message struct {
	UUID       string         // Message UUID
	ParentUUID string         // Parent message UUID for threading
	Role       string         // "user" or "assistant"
	Content    []ContentBlock // Content blocks
	Timestamp  time.Time      // Message timestamp
}

// ContentBlock represents a block of content within a message
type ContentBlock struct {
	Type       string // "text", "tool_use", "tool_result"
	Text       string // For text blocks
	ToolName   string // For tool_use blocks
	ToolInput  string // JSON string of tool input
	ToolUseID  string // For tool_result blocks
	ToolOutput string // For tool_result blocks
}

// jsonlEntry represents a raw JSONL line (used for initial parsing)
type jsonlEntry struct {
	Type      string          `json:"type"`
	Summary   string          `json:"summary,omitempty"`
	LeafUUID  string          `json:"leafUuid,omitempty"`
	UUID      string          `json:"uuid,omitempty"`
	ParentUUID *string        `json:"parentUuid,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	Message   json.RawMessage `json:"message,omitempty"`
}

// messageContent represents the message field in a JSONL entry
type messageContent struct {
	Role    string            `json:"role"`
	Content json.RawMessage   `json:"content"`
}

// rawContentBlock represents a content block as it appears in JSONL
type rawContentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
}
