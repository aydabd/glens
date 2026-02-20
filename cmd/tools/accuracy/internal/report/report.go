// Package report builds the accuracy markdown report.
package report

import (
	"fmt"
	"strings"
	"time"

	"glens/tools/accuracy/internal/analyze"
)

// Build generates a markdown accuracy report from the given results.
func Build(results []analyze.Result) string {
	var sb strings.Builder
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

	sb.WriteString("# Glens Accuracy Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", timestamp))

	total := len(results)
	passed := 0
	totalEndpoints := 0
	for _, r := range results {
		if r.Err == nil {
			passed++
			totalEndpoints += r.Endpoints
		}
	}

	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Specs Tested | %d |\n", total))
	sb.WriteString(fmt.Sprintf("| Successful | %d |\n", passed))
	sb.WriteString(fmt.Sprintf("| Failed | %d |\n", total-passed))
	sb.WriteString(fmt.Sprintf("| Total Endpoints | %d |\n", totalEndpoints))
	if total > 0 {
		sb.WriteString(fmt.Sprintf("| Success Rate | %d%% |\n", passed*100/total))
	}
	sb.WriteString("\n")

	sb.WriteString("## Results\n\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("### %s\n\n", r.Name))
		sb.WriteString(fmt.Sprintf("**Spec:** `%s`\n\n", r.SpecPath))
		sb.WriteString(fmt.Sprintf("**Duration:** %s\n\n", r.Elapsed.Round(time.Millisecond)))
		if r.Err != nil {
			sb.WriteString("**Status:** ❌ Failed\n\n")
			sb.WriteString(fmt.Sprintf("**Error:**\n```\n%v\n```\n\n", r.Err))
		} else {
			sb.WriteString("**Status:** ✅ Success\n\n")
			if r.Title != "" {
				sb.WriteString(fmt.Sprintf("**Title:** %s\n\n", r.Title))
			}
			sb.WriteString(fmt.Sprintf("**Endpoints Found:** %d\n\n", r.Endpoints))
		}
		sb.WriteString("---\n\n")
	}
	return sb.String()
}
