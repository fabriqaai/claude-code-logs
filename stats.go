package main

import (
	"sort"
	"time"
)

// StatsData holds all computed analytics for the stats dashboard
type StatsData struct {
	// Summary totals
	TotalProjects int     `json:"totalProjects"`
	TotalSessions int     `json:"totalSessions"`
	TotalMessages int     `json:"totalMessages"`
	TotalTokens   int     `json:"totalTokens"`
	TotalCost     float64 `json:"totalCost"`

	// Time series for charts (last 30 days)
	MessagesPerDay []TimePoint `json:"messagesPerDay"`
	TokensPerDay   []TimePoint `json:"tokensPerDay"`

	// Project breakdown
	ProjectStats []ProjectStat `json:"projectStats"`

	// Session statistics
	AvgSessionLengthMins  float64 `json:"avgSessionLengthMins"`
	AvgMessagesPerSession float64 `json:"avgMessagesPerSession"`

	// Computed timestamp
	ComputedAt time.Time `json:"computedAt"`
}

// TimePoint represents a single data point in a time series
type TimePoint struct {
	Date  string  `json:"date"`  // YYYY-MM-DD format
	Value int     `json:"value"`
	Cost  float64 `json:"cost"`
}

// ProjectStat holds statistics for a single project
type ProjectStat struct {
	Path                  string      `json:"path"`
	Slug                  string      `json:"slug"`
	Sessions              int         `json:"sessions"`
	Messages              int         `json:"messages"`
	Tokens                int         `json:"tokens"`
	Cost                  float64     `json:"cost"`
	LastUsed              time.Time   `json:"lastUsed"`
	MessagesPerDay        []TimePoint `json:"messagesPerDay"`
	TokensPerDay          []TimePoint `json:"tokensPerDay"`
	AvgSessionLengthMins  float64     `json:"avgSessionLengthMins"`
	AvgMessagesPerSession float64     `json:"avgMessagesPerSession"`
}

// ComputeStats calculates analytics from loaded projects
func ComputeStats(projects []Project) *StatsData {
	stats := &StatsData{
		ComputedAt:     time.Now(),
		MessagesPerDay: []TimePoint{},
		TokensPerDay:   []TimePoint{},
		ProjectStats:   []ProjectStat{},
	}

	if len(projects) == 0 {
		return stats
	}

	// Maps for daily aggregation (global)
	messagesByDay := make(map[string]int)
	tokensByDay := make(map[string]int)
	inputTokensByDay := make(map[string]int64)
	outputTokensByDay := make(map[string]int64)

	// Maps for per-project daily aggregation
	projectMessagesByDay := make(map[string]map[string]int)
	projectTokensByDay := make(map[string]map[string]int)
	projectInputTokensByDay := make(map[string]map[string]int64)
	projectOutputTokensByDay := make(map[string]map[string]int64)

	// Track session lengths for averaging
	var totalSessionMins float64
	var sessionCount int

	for _, project := range projects {
		slug := ProjectSlug(project.Path)
		projectStat := ProjectStat{
			Path:     project.Path,
			Slug:     slug,
			Sessions: len(project.Sessions),
		}

		// Initialize per-project maps
		if projectMessagesByDay[slug] == nil {
			projectMessagesByDay[slug] = make(map[string]int)
			projectTokensByDay[slug] = make(map[string]int)
			projectInputTokensByDay[slug] = make(map[string]int64)
			projectOutputTokensByDay[slug] = make(map[string]int64)
		}

		// Track per-project session stats
		var projectSessionMins float64
		var projectSessionCount int

		for _, session := range project.Sessions {
			sessionCount++
			projectSessionCount++
			stats.TotalSessions++

			// Calculate session length (first to last message)
			if len(session.Messages) > 0 {
				sessionLength := session.UpdatedAt.Sub(session.CreatedAt).Minutes()
				if sessionLength > 0 {
					totalSessionMins += sessionLength
					projectSessionMins += sessionLength
				}
			}

			// Track last used time for project
			if session.UpdatedAt.After(projectStat.LastUsed) {
				projectStat.LastUsed = session.UpdatedAt
			}

			for _, msg := range session.Messages {
				stats.TotalMessages++
				projectStat.Messages++

				// Estimate tokens from content
				tokens := estimateTokens(msg)
				stats.TotalTokens += tokens
				projectStat.Tokens += tokens

				// Track input vs output tokens based on role
				// user messages = input tokens, assistant messages = output tokens
				var inputTokens, outputTokens int64
				if msg.Role == "user" {
					inputTokens = int64(tokens)
				} else {
					outputTokens = int64(tokens)
				}

				// Calculate cost for this message
				msgCost := CalculateCost(inputTokens, outputTokens)
				stats.TotalCost += msgCost
				projectStat.Cost += msgCost

				// Aggregate by day (global)
				day := msg.Timestamp.Format("2006-01-02")
				messagesByDay[day]++
				tokensByDay[day] += tokens
				inputTokensByDay[day] += inputTokens
				outputTokensByDay[day] += outputTokens

				// Aggregate by day (per-project)
				projectMessagesByDay[slug][day]++
				projectTokensByDay[slug][day] += tokens
				projectInputTokensByDay[slug][day] += inputTokens
				projectOutputTokensByDay[slug][day] += outputTokens
			}
		}

		// Build per-project time series
		projectStat.MessagesPerDay = buildTimeSeries(projectMessagesByDay[slug], 30)
		projectStat.TokensPerDay = buildTimeSeriesWithCost(
			projectTokensByDay[slug],
			projectInputTokensByDay[slug],
			projectOutputTokensByDay[slug],
			30,
		)

		// Calculate per-project averages
		if projectSessionCount > 0 {
			projectStat.AvgSessionLengthMins = projectSessionMins / float64(projectSessionCount)
			projectStat.AvgMessagesPerSession = float64(projectStat.Messages) / float64(projectSessionCount)
		}

		stats.ProjectStats = append(stats.ProjectStats, projectStat)
	}

	stats.TotalProjects = len(projects)

	// Calculate averages
	if sessionCount > 0 {
		stats.AvgSessionLengthMins = totalSessionMins / float64(sessionCount)
		stats.AvgMessagesPerSession = float64(stats.TotalMessages) / float64(sessionCount)
	}

	// Convert daily maps to sorted time series (last 30 days)
	stats.MessagesPerDay = buildTimeSeries(messagesByDay, 30)
	stats.TokensPerDay = buildTimeSeriesWithCost(tokensByDay, inputTokensByDay, outputTokensByDay, 30)

	// Sort projects by last used (most recent first)
	sort.Slice(stats.ProjectStats, func(i, j int) bool {
		return stats.ProjectStats[i].LastUsed.After(stats.ProjectStats[j].LastUsed)
	})

	return stats
}

// estimateTokens estimates token count from message content
// Uses rough approximation: 1 token ≈ 4 characters
func estimateTokens(msg Message) int {
	var totalChars int
	for _, block := range msg.Content {
		switch block.Type {
		case "text":
			totalChars += len(block.Text)
		case "tool_use":
			totalChars += len(block.ToolInput)
		case "tool_result":
			totalChars += len(block.ToolOutput)
		}
	}
	// Rough approximation: 4 chars per token
	return totalChars / 4
}

// buildTimeSeries creates a sorted time series from a day->value map
// Returns data for the last N days, filling in zeros for missing days
func buildTimeSeries(data map[string]int, days int) []TimePoint {
	if len(data) == 0 {
		return []TimePoint{}
	}

	// Find the date range
	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1)

	result := make([]TimePoint, 0, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		value := data[dateStr] // Will be 0 if not present
		result = append(result, TimePoint{
			Date:  dateStr,
			Value: value,
		})
	}

	return result
}

// buildTimeSeriesWithCost creates a sorted time series with cost data
// Returns data for the last N days, filling in zeros for missing days
func buildTimeSeriesWithCost(tokenData map[string]int, inputTokens, outputTokens map[string]int64, days int) []TimePoint {
	if len(tokenData) == 0 {
		return []TimePoint{}
	}

	// Find the date range
	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1)

	result := make([]TimePoint, 0, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		value := tokenData[dateStr]
		inputTok := inputTokens[dateStr]
		outputTok := outputTokens[dateStr]
		cost := CalculateCost(inputTok, outputTok)
		result = append(result, TimePoint{
			Date:  dateStr,
			Value: value,
			Cost:  cost,
		})
	}

	return result
}

// FilterStatsByTimeRange filters stats data to a specific time range
func FilterStatsByTimeRange(stats *StatsData, rangeType string) *StatsData {
	var daysBack int
	switch rangeType {
	case "today":
		daysBack = 1
	case "week":
		daysBack = 7
	case "month":
		daysBack = 30
	default: // "all"
		return stats
	}

	cutoff := time.Now().AddDate(0, 0, -daysBack)
	cutoffStr := cutoff.Format("2006-01-02")

	// Filter time series
	filteredMessages := filterTimePoints(stats.MessagesPerDay, cutoffStr)
	filteredTokens := filterTimePoints(stats.TokensPerDay, cutoffStr)

	// Recalculate totals from filtered data
	var totalMessages, totalTokens int
	var totalCost float64
	for _, tp := range filteredMessages {
		totalMessages += tp.Value
	}
	for _, tp := range filteredTokens {
		totalTokens += tp.Value
		totalCost += tp.Cost
	}

	return &StatsData{
		TotalProjects:         stats.TotalProjects,
		TotalSessions:         stats.TotalSessions,
		TotalMessages:         totalMessages,
		TotalTokens:           totalTokens,
		TotalCost:             totalCost,
		MessagesPerDay:        filteredMessages,
		TokensPerDay:          filteredTokens,
		ProjectStats:          stats.ProjectStats,
		AvgSessionLengthMins:  stats.AvgSessionLengthMins,
		AvgMessagesPerSession: stats.AvgMessagesPerSession,
		ComputedAt:            stats.ComputedAt,
	}
}

// filterTimePoints filters time points to dates >= cutoff
func filterTimePoints(points []TimePoint, cutoffDate string) []TimePoint {
	result := make([]TimePoint, 0)
	for _, tp := range points {
		if tp.Date >= cutoffDate {
			result = append(result, tp)
		}
	}
	return result
}
