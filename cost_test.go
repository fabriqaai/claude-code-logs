package main

import (
	"math"
	"testing"
)

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name         string
		inputTokens  int64
		outputTokens int64
		expectedCost float64
	}{
		{
			name:         "zero tokens",
			inputTokens:  0,
			outputTokens: 0,
			expectedCost: 0.0,
		},
		{
			name:         "1 million input tokens only",
			inputTokens:  1_000_000,
			outputTokens: 0,
			expectedCost: 3.0, // $3/MTok
		},
		{
			name:         "1 million output tokens only",
			inputTokens:  0,
			outputTokens: 1_000_000,
			expectedCost: 15.0, // $15/MTok
		},
		{
			name:         "1 million each",
			inputTokens:  1_000_000,
			outputTokens: 1_000_000,
			expectedCost: 18.0, // $3 + $15
		},
		{
			name:         "realistic small usage (10k input, 2k output)",
			inputTokens:  10_000,
			outputTokens: 2_000,
			expectedCost: 0.06, // $0.03 + $0.03
		},
		{
			name:         "typical session (50k input, 10k output)",
			inputTokens:  50_000,
			outputTokens: 10_000,
			expectedCost: 0.30, // $0.15 + $0.15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateCost(tt.inputTokens, tt.outputTokens)
			// Use tolerance for floating point comparison
			if math.Abs(got-tt.expectedCost) > 0.0001 {
				t.Errorf("CalculateCost(%d, %d) = %v, want %v",
					tt.inputTokens, tt.outputTokens, got, tt.expectedCost)
			}
		})
	}
}

func TestPricingConstants(t *testing.T) {
	// Verify pricing constants match Anthropic's published rates
	if SonnetInputPricePerMTok != 3.0 {
		t.Errorf("SonnetInputPricePerMTok = %v, want 3.0", SonnetInputPricePerMTok)
	}
	if SonnetOutputPricePerMTok != 15.0 {
		t.Errorf("SonnetOutputPricePerMTok = %v, want 15.0", SonnetOutputPricePerMTok)
	}
}
