package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML metadata at the top of a Markdown file
type Frontmatter struct {
	Source     string `yaml:"source"`      // Original JSONL filename
	SourceHash string `yaml:"source_hash"` // SHA256 of JSONL content
	Project    string `yaml:"project"`     // Actual project path (from CWD)
	Title      string `yaml:"title"`       // Session summary
	Created    string `yaml:"created"`     // ISO 8601 timestamp
}

// Marshal serializes the frontmatter to YAML with delimiters
func (f Frontmatter) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(0)
	if err := encoder.Encode(f); err != nil {
		return nil, fmt.Errorf("encoding frontmatter: %w", err)
	}
	encoder.Close()

	buf.WriteString("---\n")
	return buf.Bytes(), nil
}

// ParseFrontmatter extracts frontmatter from markdown content
// Returns the frontmatter, the remaining content, and any error
func ParseFrontmatter(content []byte) (Frontmatter, []byte, error) {
	var fm Frontmatter

	str := string(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(str, "---\n") {
		return fm, content, nil // No frontmatter
	}

	// Find the closing delimiter
	rest := str[4:] // Skip opening "---\n"
	endIdx := strings.Index(rest, "\n---\n")
	if endIdx == -1 {
		// Try with just "---" at end of file
		endIdx = strings.Index(rest, "\n---")
		if endIdx == -1 {
			return fm, content, fmt.Errorf("unclosed frontmatter")
		}
	}

	// Extract YAML content
	yamlContent := rest[:endIdx]

	// Parse YAML
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return fm, content, fmt.Errorf("parsing frontmatter YAML: %w", err)
	}

	// Return remaining content (skip frontmatter + closing delimiter + newline)
	remainingStart := 4 + endIdx + 5 // "---\n" + yaml + "\n---\n"
	if remainingStart > len(str) {
		remainingStart = len(str)
	}
	remaining := []byte(str[remainingStart:])

	return fm, remaining, nil
}

// ComputeFileHash computes the SHA256 hash of a file
func ComputeFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file for hash: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("computing hash: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// NewFrontmatter creates a new Frontmatter from a Session
func NewFrontmatter(session *Session, sourceHash string) Frontmatter {
	return Frontmatter{
		Source:     filepath.Base(session.SourcePath),
		SourceHash: sourceHash,
		Project:    session.CWD,
		Title:      session.Summary,
		Created:    session.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
