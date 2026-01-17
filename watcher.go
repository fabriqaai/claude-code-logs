package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchConfig configures the file watcher
type WatchConfig struct {
	SourceDir        string        // Claude projects directory
	OutputDir        string        // HTML output directory
	PollInterval     time.Duration // Interval for scanning new directories
	DebounceDelay    time.Duration // Delay before regenerating after changes
	SelectedProjects []string      // Project folder names to watch (nil = all projects)
}

// DefaultWatchConfig returns the default watcher configuration
func DefaultWatchConfig() WatchConfig {
	return WatchConfig{
		PollInterval:  30 * time.Second,
		DebounceDelay: 2 * time.Second,
	}
}

// Watcher monitors for file changes and triggers regeneration
type Watcher struct {
	config    WatchConfig
	fsWatcher *fsnotify.Watcher

	// Debouncing state
	mu            sync.Mutex
	pendingTimers map[string]*time.Timer // project folder -> debounce timer

	// Callback for regeneration
	onRegenerate func(projectFolder string) error
}

// NewWatcher creates a new file watcher
func NewWatcher(config WatchConfig) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating fsnotify watcher: %w", err)
	}

	return &Watcher{
		config:        config,
		fsWatcher:     fsWatcher,
		pendingTimers: make(map[string]*time.Timer),
	}, nil
}

// Close closes the watcher and releases resources
func (w *Watcher) Close() error {
	w.mu.Lock()
	// Cancel all pending timers
	for _, timer := range w.pendingTimers {
		timer.Stop()
	}
	w.pendingTimers = make(map[string]*time.Timer)
	w.mu.Unlock()

	return w.fsWatcher.Close()
}

// Watch starts watching for changes (blocking until context is cancelled)
func (w *Watcher) Watch(ctx context.Context) error {
	// Add watches for source directory and all project subdirectories
	if err := w.addWatches(); err != nil {
		return fmt.Errorf("adding watches: %w", err)
	}

	if w.config.SelectedProjects != nil {
		fmt.Printf("Watching %d selected projects for changes\n", len(w.config.SelectedProjects))
	} else {
		fmt.Printf("Watching for changes in %s\n", w.config.SourceDir)
	}
	logVerbose("Poll interval: %v, Debounce delay: %v", w.config.PollInterval, w.config.DebounceDelay)

	// Start directory scanner for new project folders
	scanTicker := time.NewTicker(w.config.PollInterval)
	defer scanTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nShutting down watcher...")
			return nil

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return nil
			}
			w.handleEvent(event)

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)

		case <-scanTicker.C:
			// Periodically scan for new project directories
			w.scanForNewDirectories()
		}
	}
}

// addWatches adds watches for the source directory and all project subdirectories
func (w *Watcher) addWatches() error {
	// Watch the root projects directory (for new projects)
	if err := w.fsWatcher.Add(w.config.SourceDir); err != nil {
		return fmt.Errorf("watching source directory: %w", err)
	}

	// Watch each project subdirectory
	entries, err := os.ReadDir(w.config.SourceDir)
	if err != nil {
		return fmt.Errorf("reading source directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip if we have a filter and this project isn't in it
		if !w.isProjectSelected(entry.Name()) {
			continue
		}

		projectPath := filepath.Join(w.config.SourceDir, entry.Name())
		if err := w.fsWatcher.Add(projectPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to watch %s: %v\n", projectPath, err)
			continue
		}
	}

	return nil
}

// isProjectSelected returns true if the project should be watched
// Returns true if no filter is set (SelectedProjects is nil) or if the project is in the filter
func (w *Watcher) isProjectSelected(projectFolder string) bool {
	if w.config.SelectedProjects == nil {
		return true // No filter, watch all
	}
	for _, p := range w.config.SelectedProjects {
		if p == projectFolder {
			return true
		}
	}
	return false
}

// scanForNewDirectories checks for new project directories and adds watches
func (w *Watcher) scanForNewDirectories() {
	// Skip scanning for new directories if we have a specific project filter
	if w.config.SelectedProjects != nil {
		return
	}

	entries, err := os.ReadDir(w.config.SourceDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to scan for new directories: %v\n", err)
		return
	}

	watchList := w.fsWatcher.WatchList()
	watchSet := make(map[string]bool)
	for _, path := range watchList {
		watchSet[path] = true
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projectPath := filepath.Join(w.config.SourceDir, entry.Name())
		if !watchSet[projectPath] {
			if err := w.fsWatcher.Add(projectPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to watch new directory %s: %v\n", projectPath, err)
				continue
			}
			fmt.Printf("Now watching new project: %s\n", entry.Name())

			// Trigger regeneration for the new project
			w.scheduleRegeneration(entry.Name())
		}
	}
}

// handleEvent processes a single fsnotify event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only care about .jsonl files
	if !strings.HasSuffix(event.Name, ".jsonl") {
		// Check if it's a new directory being created (only if no filter)
		if event.Op&fsnotify.Create != 0 && w.config.SelectedProjects == nil {
			info, err := os.Stat(event.Name)
			if err == nil && info.IsDir() {
				// New project directory - add watch
				if err := w.fsWatcher.Add(event.Name); err == nil {
					projectFolder := filepath.Base(event.Name)
					fmt.Printf("Now watching new project: %s\n", projectFolder)
					w.scheduleRegeneration(projectFolder)
				}
			}
		}
		return
	}

	// Extract project folder from path
	// Path format: /path/to/projects/-Project-Name/session.jsonl
	dir := filepath.Dir(event.Name)
	projectFolder := filepath.Base(dir)

	// Skip if project is not in our selection
	if !w.isProjectSelected(projectFolder) {
		return
	}

	// Handle different event types
	switch {
	case event.Op&fsnotify.Write != 0:
		logVerbose("Modified: %s", filepath.Base(event.Name))
		w.scheduleRegeneration(projectFolder)

	case event.Op&fsnotify.Create != 0:
		logVerbose("Created: %s", filepath.Base(event.Name))
		w.scheduleRegeneration(projectFolder)

	case event.Op&fsnotify.Remove != 0:
		logVerbose("Removed: %s", filepath.Base(event.Name))
		w.scheduleRegeneration(projectFolder)

	case event.Op&fsnotify.Rename != 0:
		logVerbose("Renamed: %s", filepath.Base(event.Name))
		w.scheduleRegeneration(projectFolder)
	}
}

// scheduleRegeneration schedules a regeneration for a project with debouncing
func (w *Watcher) scheduleRegeneration(projectFolder string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel existing timer for this project
	if timer, exists := w.pendingTimers[projectFolder]; exists {
		timer.Stop()
	}

	// Create new timer
	w.pendingTimers[projectFolder] = time.AfterFunc(w.config.DebounceDelay, func() {
		w.regenerateProject(projectFolder)
	})
}

// regenerateProject regenerates HTML for a single project
func (w *Watcher) regenerateProject(projectFolder string) {
	w.mu.Lock()
	delete(w.pendingTimers, projectFolder)
	w.mu.Unlock()

	fmt.Printf("Regenerating: %s\n", DecodeProjectPath(projectFolder))

	if w.onRegenerate != nil {
		if err := w.onRegenerate(projectFolder); err != nil {
			fmt.Fprintf(os.Stderr, "Error regenerating %s: %v\n", projectFolder, err)
		}
	}
}

// SetRegenerateCallback sets the callback function for regeneration
func (w *Watcher) SetRegenerateCallback(fn func(projectFolder string) error) {
	w.onRegenerate = fn
}

// StartWatcher starts watching with default regeneration behavior
func StartWatcher(ctx context.Context, config WatchConfig) error {
	watcher, err := NewWatcher(config)
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Set up regeneration callback
	watcher.SetRegenerateCallback(func(projectFolder string) error {
		return regenerateProject(config.SourceDir, config.OutputDir, projectFolder)
	})

	return watcher.Watch(ctx)
}

// regenerateProject reloads a project and regenerates its Markdown files
func regenerateProject(sourceDir, outputDir, projectFolder string) error {
	_ = projectFolder // Currently regenerates all projects; mtime check handles efficiency

	// Load all projects (needed for index files)
	allProjects, err := LoadAllProjects(sourceDir)
	if err != nil {
		return fmt.Errorf("loading all projects: %w", err)
	}

	// Generate markdown (with force=false for incremental updates, no static HTML)
	result, err := GenerateAllMarkdown(allProjects, outputDir, sourceDir, false, false)
	if err != nil {
		return fmt.Errorf("generating markdown: %w", err)
	}

	fmt.Printf("Regenerated: %d generated, %d skipped\n", result.Generated, result.Skipped)
	return nil
}

// WatchInBackground starts the watcher in a background goroutine
// Returns a cancel function to stop the watcher
func WatchInBackground(config WatchConfig) (context.CancelFunc, error) {
	watcher, err := NewWatcher(config)
	if err != nil {
		return nil, err
	}

	// Set up regeneration callback
	watcher.SetRegenerateCallback(func(projectFolder string) error {
		return regenerateProject(config.SourceDir, config.OutputDir, projectFolder)
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := watcher.Watch(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
		watcher.Close()
	}()

	return cancel, nil
}
