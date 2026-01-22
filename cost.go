package main

// Anthropic Claude pricing constants (per million tokens)
// Based on: https://platform.claude.com/docs/en/about-claude/pricing
const (
	// Sonnet 4 pricing (default)
	SonnetInputPricePerMTok  = 3.0  // $3.00 per million input tokens
	SonnetOutputPricePerMTok = 15.0 // $15.00 per million output tokens
)

// CalculateCost computes the estimated API cost based on token counts.
// Uses Sonnet 4 pricing by default.
// Returns cost in USD.
func CalculateCost(inputTokens, outputTokens int64) float64 {
	inputCost := float64(inputTokens) / 1_000_000 * SonnetInputPricePerMTok
	outputCost := float64(outputTokens) / 1_000_000 * SonnetOutputPricePerMTok
	return inputCost + outputCost
}
