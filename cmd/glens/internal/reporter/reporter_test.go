package reporter

import (
	"testing"
	"time"
)

func TestCalculateExecutionSummary_SuccessRate(t *testing.T) {
	exec := []time.Duration{time.Second}
	gen := []time.Duration{}

	tests := []struct {
		name     string
		passed   int
		total    int
		wantRate float64
	}{
		{"no tests", 0, 0, 0.0},
		{"all passed", 3, 3, 1.0},
		{"partial", 1, 4, 0.25},
		{"all failed", 0, 5, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateExecutionSummary(exec, gen, tt.passed, tt.total)
			if got.SuccessRate != tt.wantRate {
				t.Errorf("SuccessRate = %v, want %v", got.SuccessRate, tt.wantRate)
			}
		})
	}
}
