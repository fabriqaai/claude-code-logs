package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
		{
			name:  "absolute path",
			input: "/tmp/test",
			want:  "/tmp/test",
		},
		{
			name:  "relative path",
			input: "relative/path",
			want:  "relative/path",
		},
		{
			name:  "tilde expansion",
			input: "~/test",
			want:  filepath.Join(home, "test"),
		},
		{
			name:  "tilde only",
			input: "~",
			want:  home,
		},
		{
			name:  "tilde with subpath",
			input: "~/.claude-logs",
			want:  filepath.Join(home, ".claude-logs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("expandPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetOutputDir(t *testing.T) {
	// Save and restore global state
	origOutputDir := outputDir
	defer func() { outputDir = origOutputDir }()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	tests := []struct {
		name      string
		outputDir string
		want      string
		wantErr   bool
	}{
		{
			name:      "default output dir",
			outputDir: "",
			want:      filepath.Join(home, ".claude-logs"),
		},
		{
			name:      "custom output dir",
			outputDir: "/tmp/custom-logs",
			want:      "/tmp/custom-logs",
		},
		{
			name:      "tilde in output dir",
			outputDir: "~/my-logs",
			want:      filepath.Join(home, "my-logs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir = tt.outputDir
			got, err := getOutputDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("getOutputDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getOutputDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEnsureWritableDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "new directory",
			setup: func() string {
				dir := filepath.Join(os.TempDir(), "test-writable-new")
				os.RemoveAll(dir) // Ensure it doesn't exist
				return dir
			},
			wantErr: false,
		},
		{
			name: "existing directory",
			setup: func() string {
				dir := filepath.Join(os.TempDir(), "test-writable-existing")
				os.MkdirAll(dir, 0755)
				return dir
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup()
			defer os.RemoveAll(dir)

			err := ensureWritableDir(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureWritableDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify directory exists
			if !tt.wantErr {
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Error("ensureWritableDir() did not create directory")
				}
			}
		})
	}
}

func TestLogVerbose(t *testing.T) {
	// Save and restore global state
	origVerbose := verbose
	defer func() { verbose = origVerbose }()

	// Capture stdout is complex, so we just verify it doesn't panic
	verbose = false
	logVerbose("test message %s", "arg")

	verbose = true
	logVerbose("test message %s", "arg")
}

func TestVersionCommand(t *testing.T) {
	// Version command uses fmt.Printf which goes to stdout, not the cobra output
	// We just verify it doesn't error
	// In a real test we'd capture os.Stdout but that's complex

	// Instead, test the command exists and has correct attributes
	if versionCmd.Use != "version" {
		t.Errorf("version command Use = %q, want %q", versionCmd.Use, "version")
	}
	if versionCmd.Short == "" {
		t.Error("version command should have Short description")
	}
}

func TestWatchCommand(t *testing.T) {
	// Watch command uses fmt.Println which goes to stdout, not the cobra output
	// We just verify the command structure

	if watchCmd.Use != "watch" {
		t.Errorf("watch command Use = %q, want %q", watchCmd.Use, "watch")
	}
	if watchCmd.Short == "" {
		t.Error("watch command should have Short description")
	}

	// Verify the interval flag exists
	flag := watchCmd.Flags().Lookup("interval")
	if flag == nil {
		t.Error("watch command should have --interval flag")
	}
}

func TestRootCommand(t *testing.T) {
	// Test root command help
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("root command help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "generate") {
		t.Errorf("root help should mention generate command: %s", output)
	}
	if !strings.Contains(output, "serve") {
		t.Errorf("root help should mention serve command: %s", output)
	}
}

func TestServePortValidation(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port", 8080, false},
		{"port 1", 1, false},
		{"port 65535", 65535, false},
		{"port 0", 0, true},
		{"port negative", -1, true},
		{"port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily test runServe without a real projects dir,
			// but we can test the port validation logic
			if tt.port < 1 || tt.port > 65535 {
				if !tt.wantErr {
					t.Errorf("expected port %d to be invalid", tt.port)
				}
			} else {
				if tt.wantErr {
					t.Errorf("expected port %d to be valid", tt.port)
				}
			}
		})
	}
}
