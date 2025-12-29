package main

import (
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// SearchIndex provides full-text search over chat messages
type SearchIndex struct {
	// Inverted index: term -> message indices
	index map[string][]int

	// All indexed messages
	messages []IndexedMessage
}

// IndexedMessage represents a searchable message
type IndexedMessage struct {
	Project      string
	ProjectSlug  string
	SessionID    string
	SessionTitle string
	MessageID    string
	Role         string
	Content      string
	Timestamp    time.Time
}

// SearchResult represents a search result for a session
type SearchResult struct {
	Project      string        `json:"project"`
	ProjectSlug  string        `json:"projectSlug"`
	SessionID    string        `json:"sessionId"`
	SessionTitle string        `json:"sessionTitle"`
	Matches      []MatchResult `json:"matches"`
}

// MatchResult represents a matching message
type MatchResult struct {
	MessageID string    `json:"messageId"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// SearchResponse is the API response format
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Query   string         `json:"query"`
}

// SearchRequest is the API request format
type SearchRequest struct {
	Query   string `json:"query"`
	Project string `json:"project,omitempty"`
	Session string `json:"session,omitempty"`
}

// NewSearchIndex creates a new search index from projects
func NewSearchIndex(projects []Project) *SearchIndex {
	idx := &SearchIndex{
		index:    make(map[string][]int),
		messages: []IndexedMessage{},
	}

	// Index all messages
	for _, project := range projects {
		projectSlug := ProjectSlug(project.Path)
		for _, session := range project.Sessions {
			for _, msg := range session.Messages {
				// Only index text content
				content := extractTextContent(msg)
				if content == "" {
					continue
				}

				msgIndex := len(idx.messages)
				idx.messages = append(idx.messages, IndexedMessage{
					Project:      project.Path,
					ProjectSlug:  projectSlug,
					SessionID:    session.ID,
					SessionTitle: session.Summary,
					MessageID:    msg.UUID,
					Role:         msg.Role,
					Content:      content,
					Timestamp:    msg.Timestamp,
				})

				// Tokenize and index
				terms := tokenize(content)
				for _, term := range terms {
					idx.index[term] = append(idx.index[term], msgIndex)
				}
			}
		}
	}

	return idx
}

// Search executes a search query with optional filters
func (idx *SearchIndex) Search(query, projectFilter, sessionFilter string) []SearchResult {
	if query == "" {
		return []SearchResult{}
	}

	// Parse query to extract phrases and terms
	parsed := parseQuery(query)
	if len(parsed.Terms) == 0 && len(parsed.Phrases) == 0 {
		return []SearchResult{}
	}

	// Find messages containing all query terms (AND logic)
	matchingIndices := idx.findMatchingMessages(parsed.Terms)

	// Apply filters, phrase matching, and group by session
	sessionMatches := make(map[string][]int) // sessionKey -> message indices
	for _, msgIdx := range matchingIndices {
		msg := idx.messages[msgIdx]

		// Apply project filter
		if projectFilter != "" && msg.Project != projectFilter && msg.ProjectSlug != projectFilter {
			continue
		}

		// Apply session filter
		if sessionFilter != "" && msg.SessionID != sessionFilter {
			continue
		}

		// Apply phrase filter - check that all phrases appear in the content
		if len(parsed.Phrases) > 0 {
			lowerContent := strings.ToLower(msg.Content)
			allPhrasesMatch := true
			for _, phrase := range parsed.Phrases {
				if !strings.Contains(lowerContent, phrase) {
					allPhrasesMatch = false
					break
				}
			}
			if !allPhrasesMatch {
				continue
			}
		}

		sessionKey := msg.ProjectSlug + "|" + msg.SessionID
		sessionMatches[sessionKey] = append(sessionMatches[sessionKey], msgIdx)
	}

	// Build results
	results := []SearchResult{}
	for _, indices := range sessionMatches {
		if len(indices) == 0 {
			continue
		}

		firstMsg := idx.messages[indices[0]]
		result := SearchResult{
			Project:      firstMsg.Project,
			ProjectSlug:  firstMsg.ProjectSlug,
			SessionID:    firstMsg.SessionID,
			SessionTitle: firstMsg.SessionTitle,
			Matches:      []MatchResult{},
		}

		for _, msgIdx := range indices {
			msg := idx.messages[msgIdx]
			highlighted := highlightMatchesWithPhrases(msg.Content, parsed.Terms, parsed.Phrases)
			result.Matches = append(result.Matches, MatchResult{
				MessageID: msg.MessageID,
				Role:      msg.Role,
				Content:   highlighted,
				Timestamp: msg.Timestamp,
			})
		}

		// Sort matches by timestamp
		sort.Slice(result.Matches, func(i, j int) bool {
			return result.Matches[i].Timestamp.Before(result.Matches[j].Timestamp)
		})

		results = append(results, result)
	}

	// Sort results by most recent match
	sort.Slice(results, func(i, j int) bool {
		if len(results[i].Matches) == 0 || len(results[j].Matches) == 0 {
			return false
		}
		return results[i].Matches[len(results[i].Matches)-1].Timestamp.After(
			results[j].Matches[len(results[j].Matches)-1].Timestamp,
		)
	})

	return results
}

// findMatchingMessages finds message indices containing all query terms
func (idx *SearchIndex) findMatchingMessages(terms []string) []int {
	if len(terms) == 0 {
		return []int{}
	}

	// Start with first term's matches
	firstTerm := terms[0]
	candidates := make(map[int]bool)
	for _, msgIdx := range idx.index[firstTerm] {
		candidates[msgIdx] = true
	}

	// Intersect with remaining terms
	for _, term := range terms[1:] {
		termMatches := make(map[int]bool)
		for _, msgIdx := range idx.index[term] {
			termMatches[msgIdx] = true
		}

		// Keep only candidates that also match this term
		for msgIdx := range candidates {
			if !termMatches[msgIdx] {
				delete(candidates, msgIdx)
			}
		}
	}

	// Convert to slice
	result := make([]int, 0, len(candidates))
	for msgIdx := range candidates {
		result = append(result, msgIdx)
	}

	return result
}

// extractTextContent extracts text from message content blocks
func extractTextContent(msg Message) string {
	var parts []string
	for _, block := range msg.Content {
		if block.Type == "text" && block.Text != "" {
			parts = append(parts, block.Text)
		}
	}
	return strings.Join(parts, " ")
}

// ParsedQuery represents a parsed search query with terms and phrases
type ParsedQuery struct {
	Terms   []string // Individual words (for index lookup)
	Phrases []string // Exact phrases to match (from quoted strings)
}

// parseQuery extracts quoted phrases and individual terms from a query
func parseQuery(query string) ParsedQuery {
	var result ParsedQuery

	// Extract quoted phrases
	quoteRegex := regexp.MustCompile(`"([^"]+)"`)
	matches := quoteRegex.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		if len(match) > 1 && strings.TrimSpace(match[1]) != "" {
			result.Phrases = append(result.Phrases, strings.ToLower(strings.TrimSpace(match[1])))
		}
	}

	// Remove quoted parts from query to get remaining terms
	remaining := quoteRegex.ReplaceAllString(query, " ")
	result.Terms = tokenize(remaining)

	// Also add tokenized versions of phrases for index lookup
	for _, phrase := range result.Phrases {
		phraseTerms := tokenize(phrase)
		for _, term := range phraseTerms {
			// Add if not already present
			found := false
			for _, existing := range result.Terms {
				if existing == term {
					found = true
					break
				}
			}
			if !found {
				result.Terms = append(result.Terms, term)
			}
		}
	}

	return result
}

// tokenize splits text into lowercase terms
func tokenize(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Split on non-alphanumeric characters
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Filter empty and very short words
	var terms []string
	seen := make(map[string]bool)
	for _, word := range words {
		if len(word) < 2 {
			continue
		}
		if seen[word] {
			continue
		}
		seen[word] = true
		terms = append(terms, word)
	}

	return terms
}

// highlightMatches wraps matching terms in <mark> tags and extracts context
func highlightMatches(content string, terms []string) string {
	return highlightMatchesWithPhrases(content, terms, nil)
}

// highlightMatchesWithPhrases wraps matching terms and phrases in <mark> tags and extracts context
func highlightMatchesWithPhrases(content string, terms []string, phrases []string) string {
	// Limit content length for display
	const maxLength = 500
	const contextChars = 100

	// Find first match position (check phrases first, then terms)
	lowerContent := strings.ToLower(content)
	firstMatchPos := -1

	// Check phrases first (they're more specific)
	for _, phrase := range phrases {
		pos := strings.Index(lowerContent, phrase)
		if pos != -1 && (firstMatchPos == -1 || pos < firstMatchPos) {
			firstMatchPos = pos
		}
	}

	// Then check individual terms
	for _, term := range terms {
		pos := strings.Index(lowerContent, term)
		if pos != -1 && (firstMatchPos == -1 || pos < firstMatchPos) {
			firstMatchPos = pos
		}
	}

	// Extract context around first match
	var excerpt string
	if len(content) <= maxLength {
		excerpt = content
	} else if firstMatchPos == -1 {
		excerpt = content[:maxLength] + "..."
	} else {
		start := firstMatchPos - contextChars
		if start < 0 {
			start = 0
		}
		end := firstMatchPos + contextChars
		if end > len(content) {
			end = len(content)
		}

		excerpt = ""
		if start > 0 {
			excerpt = "..."
		}
		excerpt += content[start:end]
		if end < len(content) {
			excerpt += "..."
		}
	}

	// Highlight matches (case-insensitive)
	// Highlight phrases first (longer matches take precedence)
	result := excerpt
	for _, phrase := range phrases {
		pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(phrase))
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			return "<mark>" + match + "</mark>"
		})
	}

	// Then highlight individual terms (skip if already inside a mark tag)
	for _, term := range terms {
		// Create case-insensitive regex
		pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(term))
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			return "<mark>" + match + "</mark>"
		})
	}

	return result
}

// MessageCount returns the number of indexed messages
func (idx *SearchIndex) MessageCount() int {
	return len(idx.messages)
}

// TermCount returns the number of unique terms
func (idx *SearchIndex) TermCount() int {
	return len(idx.index)
}
