package main

import (
	"testing"
	"time"
)

func TestComputeStats_EmptyProjects(t *testing.T) {
	stats := ComputeStats([]Project{})

	if stats.TotalProjects != 0 {
		t.Errorf("Expected TotalProjects=0, got %d", stats.TotalProjects)
	}
	if stats.TotalSessions != 0 {
		t.Errorf("Expected TotalSessions=0, got %d", stats.TotalSessions)
	}
	if stats.TotalMessages != 0 {
		t.Errorf("Expected TotalMessages=0, got %d", stats.TotalMessages)
	}
	if stats.TotalTokens != 0 {
		t.Errorf("Expected TotalTokens=0, got %d", stats.TotalTokens)
	}
}

func TestComputeStats_SingleProject(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/test/project",
			FolderName: "-test-project",
			Sessions: []Session{
				{
					ID:        "session-1",
					Summary:   "Test session",
					CreatedAt: now.Add(-1 * time.Hour),
					UpdatedAt: now,
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now.Add(-1 * time.Hour),
							Content: []ContentBlock{
								{Type: "text", Text: "Hello, this is a test message"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Hello! How can I help you today?"},
							},
						},
					},
				},
			},
		},
	}

	stats := ComputeStats(projects)

	if stats.TotalProjects != 1 {
		t.Errorf("Expected TotalProjects=1, got %d", stats.TotalProjects)
	}
	if stats.TotalSessions != 1 {
		t.Errorf("Expected TotalSessions=1, got %d", stats.TotalSessions)
	}
	if stats.TotalMessages != 2 {
		t.Errorf("Expected TotalMessages=2, got %d", stats.TotalMessages)
	}
	if stats.TotalTokens == 0 {
		t.Error("Expected TotalTokens > 0")
	}
	if len(stats.ProjectStats) != 1 {
		t.Errorf("Expected 1 ProjectStat, got %d", len(stats.ProjectStats))
	}
	if stats.ProjectStats[0].Path != "/test/project" {
		t.Errorf("Expected ProjectStats[0].Path='/test/project', got %s", stats.ProjectStats[0].Path)
	}
}

func TestComputeStats_MultipleProjects(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/project/one",
			FolderName: "-project-one",
			Sessions: []Session{
				{
					ID:        "session-1",
					CreatedAt: now.Add(-2 * time.Hour),
					UpdatedAt: now.Add(-1 * time.Hour),
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: now.Add(-2 * time.Hour), Content: []ContentBlock{{Type: "text", Text: "Message 1"}}},
					},
				},
			},
		},
		{
			Path:       "/project/two",
			FolderName: "-project-two",
			Sessions: []Session{
				{
					ID:        "session-2",
					CreatedAt: now.Add(-30 * time.Minute),
					UpdatedAt: now,
					Messages: []Message{
						{UUID: "msg-2", Role: "user", Timestamp: now.Add(-30 * time.Minute), Content: []ContentBlock{{Type: "text", Text: "Message 2"}}},
						{UUID: "msg-3", Role: "assistant", Timestamp: now, Content: []ContentBlock{{Type: "text", Text: "Message 3"}}},
					},
				},
			},
		},
	}

	stats := ComputeStats(projects)

	if stats.TotalProjects != 2 {
		t.Errorf("Expected TotalProjects=2, got %d", stats.TotalProjects)
	}
	if stats.TotalSessions != 2 {
		t.Errorf("Expected TotalSessions=2, got %d", stats.TotalSessions)
	}
	if stats.TotalMessages != 3 {
		t.Errorf("Expected TotalMessages=3, got %d", stats.TotalMessages)
	}

	// Project two should be first (more recent)
	if stats.ProjectStats[0].Path != "/project/two" {
		t.Errorf("Expected most recent project first, got %s", stats.ProjectStats[0].Path)
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		msg      Message
		expected int
	}{
		{
			name: "text content",
			msg: Message{
				Content: []ContentBlock{
					{Type: "text", Text: "Hello world"}, // 11 chars -> 2 tokens
				},
			},
			expected: 2,
		},
		{
			name: "tool use content",
			msg: Message{
				Content: []ContentBlock{
					{Type: "tool_use", ToolInput: `{"file": "test.txt"}`}, // 20 chars -> 5 tokens
				},
			},
			expected: 5,
		},
		{
			name: "tool result content",
			msg: Message{
				Content: []ContentBlock{
					{Type: "tool_result", ToolOutput: "File contents here"}, // 18 chars -> 4 tokens
				},
			},
			expected: 4,
		},
		{
			name: "mixed content",
			msg: Message{
				Content: []ContentBlock{
					{Type: "text", Text: "Let me read that"},       // 16 chars
					{Type: "tool_use", ToolInput: `{"path":"x"}`}, // 12 chars
				},
			},
			expected: 7, // 28 / 4 = 7
		},
		{
			name: "empty content",
			msg: Message{
				Content: []ContentBlock{},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := estimateTokens(tt.msg)
			if result != tt.expected {
				t.Errorf("estimateTokens() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestBuildTimeSeries(t *testing.T) {
	data := map[string]int{
		time.Now().Format("2006-01-02"):                    10,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"): 5,
	}

	result := buildTimeSeries(data, 7)

	if len(result) != 7 {
		t.Errorf("Expected 7 data points, got %d", len(result))
	}

	// Check that today's value is present
	today := time.Now().Format("2006-01-02")
	found := false
	for _, tp := range result {
		if tp.Date == today && tp.Value == 10 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected today's value (10) in time series")
	}
}

func TestBuildTimeSeries_Empty(t *testing.T) {
	result := buildTimeSeries(map[string]int{}, 7)

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty input, got %d items", len(result))
	}
}

func TestFilterStatsByTimeRange(t *testing.T) {
	now := time.Now()
	stats := &StatsData{
		TotalProjects: 1,
		TotalSessions: 2,
		TotalMessages: 100,
		TotalTokens:   1000,
		MessagesPerDay: []TimePoint{
			{Date: now.AddDate(0, 0, -10).Format("2006-01-02"), Value: 20},
			{Date: now.AddDate(0, 0, -3).Format("2006-01-02"), Value: 30},
			{Date: now.Format("2006-01-02"), Value: 50},
		},
		TokensPerDay: []TimePoint{
			{Date: now.AddDate(0, 0, -10).Format("2006-01-02"), Value: 200},
			{Date: now.AddDate(0, 0, -3).Format("2006-01-02"), Value: 300},
			{Date: now.Format("2006-01-02"), Value: 500},
		},
	}

	// Filter to last 7 days (should exclude -10 day)
	filtered := FilterStatsByTimeRange(stats, "week")

	if len(filtered.MessagesPerDay) != 2 {
		t.Errorf("Expected 2 data points after week filter, got %d", len(filtered.MessagesPerDay))
	}
	if filtered.TotalMessages != 80 { // 30 + 50
		t.Errorf("Expected TotalMessages=80, got %d", filtered.TotalMessages)
	}

	// "all" should return original stats
	allStats := FilterStatsByTimeRange(stats, "all")
	if allStats != stats {
		t.Error("Expected 'all' filter to return original stats")
	}
}

func TestAvgSessionLength(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path: "/test",
			Sessions: []Session{
				{
					ID:        "s1",
					CreatedAt: now.Add(-60 * time.Minute),
					UpdatedAt: now.Add(-30 * time.Minute), // 30 min session
					Messages: []Message{
						{UUID: "m1", Timestamp: now.Add(-60 * time.Minute), Content: []ContentBlock{{Type: "text", Text: "start"}}},
						{UUID: "m2", Timestamp: now.Add(-30 * time.Minute), Content: []ContentBlock{{Type: "text", Text: "end"}}},
					},
				},
				{
					ID:        "s2",
					CreatedAt: now.Add(-20 * time.Minute),
					UpdatedAt: now.Add(-10 * time.Minute), // 10 min session
					Messages: []Message{
						{UUID: "m3", Timestamp: now.Add(-20 * time.Minute), Content: []ContentBlock{{Type: "text", Text: "start"}}},
						{UUID: "m4", Timestamp: now.Add(-10 * time.Minute), Content: []ContentBlock{{Type: "text", Text: "end"}}},
					},
				},
			},
		},
	}

	stats := ComputeStats(projects)

	// Average should be (30 + 10) / 2 = 20 minutes
	if stats.AvgSessionLengthMins != 20 {
		t.Errorf("Expected AvgSessionLengthMins=20, got %f", stats.AvgSessionLengthMins)
	}

	// Average messages per session should be 4/2 = 2
	if stats.AvgMessagesPerSession != 2 {
		t.Errorf("Expected AvgMessagesPerSession=2, got %f", stats.AvgMessagesPerSession)
	}
}
