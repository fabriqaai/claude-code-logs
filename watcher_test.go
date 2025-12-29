package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestDefaultWatchConfig(t *testing.T) {
	config := DefaultWatchConfig()

	if config.PollInterval != 30*time.Second {
		t.Errorf("expected PollInterval 30s, got %v", config.PollInterval)
	}

	if config.DebounceDelay != 2*time.Second {
		t.Errorf("expected DebounceDelay 2s, got %v", config.DebounceDelay)
	}
}

func TestNewWatcher(t *testing.T) {
	config := WatchConfig{
		SourceDir:     t.TempDir(),
		OutputDir:     t.TempDir(),
		PollInterval:  1 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	if watcher.config.SourceDir != config.SourceDir {
		t.Errorf("expected SourceDir %s, got %s", config.SourceDir, watcher.config.SourceDir)
	}
}

func TestWatcherClose(t *testing.T) {
	config := WatchConfig{
		SourceDir:     t.TempDir(),
		OutputDir:     t.TempDir(),
		PollInterval:  1 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}

	// Close should not error
	if err := watcher.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestWatcherDebouncing(t *testing.T) {
	// Create temp directories
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a project directory
	projectDir := filepath.Join(sourceDir, "-test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second, // Long interval, we won't wait for it
		DebounceDelay: 200 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Track regeneration calls
	var mu sync.Mutex
	regenerations := make(map[string]int)

	watcher.SetRegenerateCallback(func(projectFolder string) error {
		mu.Lock()
		regenerations[projectFolder]++
		mu.Unlock()
		return nil
	})

	// Schedule multiple rapid regenerations for same project
	for i := 0; i < 5; i++ {
		watcher.scheduleRegeneration("-test-project")
		time.Sleep(50 * time.Millisecond) // Less than debounce delay
	}

	// Wait for debounce to fire
	time.Sleep(400 * time.Millisecond)

	mu.Lock()
	count := regenerations["-test-project"]
	mu.Unlock()

	// Should only have fired once due to debouncing
	if count != 1 {
		t.Errorf("expected 1 regeneration (debounced), got %d", count)
	}
}

func TestWatcherMultipleProjects(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// Create two project directories
	for _, name := range []string{"-project-a", "-project-b"} {
		if err := os.MkdirAll(filepath.Join(sourceDir, name), 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}
	}

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Track regeneration calls
	var mu sync.Mutex
	regenerations := make(map[string]int)

	watcher.SetRegenerateCallback(func(projectFolder string) error {
		mu.Lock()
		regenerations[projectFolder]++
		mu.Unlock()
		return nil
	})

	// Schedule regenerations for different projects
	watcher.scheduleRegeneration("-project-a")
	watcher.scheduleRegeneration("-project-b")

	// Wait for debounce to fire
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	countA := regenerations["-project-a"]
	countB := regenerations["-project-b"]
	mu.Unlock()

	// Each project should regenerate independently
	if countA != 1 {
		t.Errorf("expected 1 regeneration for project-a, got %d", countA)
	}
	if countB != 1 {
		t.Errorf("expected 1 regeneration for project-b, got %d", countB)
	}
}

func TestWatcherContextCancellation(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  1 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Run watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- watcher.Watch(ctx)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for watcher to stop
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Watch() returned error on cancel: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Watch() did not stop within timeout after context cancellation")
	}
}

func TestWatcherFileDetection(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a project directory
	projectDir := filepath.Join(sourceDir, "-test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Track regeneration calls
	var mu sync.Mutex
	regenerated := false

	watcher.SetRegenerateCallback(func(projectFolder string) error {
		mu.Lock()
		regenerated = true
		mu.Unlock()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watching in background
	go func() {
		watcher.Watch(ctx)
	}()

	// Give watcher time to start
	time.Sleep(200 * time.Millisecond)

	// Create a new .jsonl file
	testFile := filepath.Join(projectDir, "test-session.jsonl")
	if err := os.WriteFile(testFile, []byte(`{"type":"summary","summary":"test"}`), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Wait for detection and debounce
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	wasRegenerated := regenerated
	mu.Unlock()

	if !wasRegenerated {
		t.Error("expected regeneration after file creation, but none occurred")
	}
}

func TestWatcherIgnoresNonJSONL(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a project directory
	projectDir := filepath.Join(sourceDir, "-test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Track regeneration calls
	var mu sync.Mutex
	regenerated := false

	watcher.SetRegenerateCallback(func(projectFolder string) error {
		mu.Lock()
		regenerated = true
		mu.Unlock()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watching in background
	go func() {
		watcher.Watch(ctx)
	}()

	// Give watcher time to start
	time.Sleep(200 * time.Millisecond)

	// Create a non-jsonl file (should be ignored)
	testFile := filepath.Join(projectDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Wait a bit to ensure no regeneration happens
	time.Sleep(400 * time.Millisecond)

	mu.Lock()
	wasRegenerated := regenerated
	mu.Unlock()

	if wasRegenerated {
		t.Error("expected no regeneration for non-.jsonl file, but one occurred")
	}
}

func TestWatcherFileModification(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a project directory with existing file
	projectDir := filepath.Join(sourceDir, "-test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	testFile := filepath.Join(projectDir, "existing-session.jsonl")
	if err := os.WriteFile(testFile, []byte(`{"type":"summary","summary":"original"}`), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	watcher, err := NewWatcher(config)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Track regeneration calls
	var mu sync.Mutex
	regenerationCount := 0

	watcher.SetRegenerateCallback(func(projectFolder string) error {
		mu.Lock()
		regenerationCount++
		mu.Unlock()
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watching in background
	go func() {
		watcher.Watch(ctx)
	}()

	// Give watcher time to start
	time.Sleep(200 * time.Millisecond)

	// Modify the existing file
	if err := os.WriteFile(testFile, []byte(`{"type":"summary","summary":"modified"}`), 0644); err != nil {
		t.Fatalf("failed to modify test file: %v", err)
	}

	// Wait for detection and debounce
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	count := regenerationCount
	mu.Unlock()

	if count < 1 {
		t.Error("expected regeneration after file modification, but none occurred")
	}
}

func TestWatchInBackground(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outputDir,
		PollInterval:  10 * time.Second,
		DebounceDelay: 100 * time.Millisecond,
	}

	cancel, err := WatchInBackground(config)
	if err != nil {
		t.Fatalf("WatchInBackground() failed: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel should not block
	done := make(chan struct{})
	go func() {
		cancel()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("cancel() blocked for too long")
	}
}
