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
	fmt.Fprintf(&sb, "**Generated:** %s\n\n", timestamp)

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
	fmt.Fprintf(&sb, "| Specs Tested | %d |\n", total)
	fmt.Fprintf(&sb, "| Successful | %d |\n", passed)
	fmt.Fprintf(&sb, "| Failed | %d |\n", total-passed)
	fmt.Fprintf(&sb, "| Total Endpoints | %d |\n", totalEndpoints)
	if total > 0 {
		fmt.Fprintf(&sb, "| Success Rate | %d%% |\n", passed*100/total)
	}
	sb.WriteString("\n")

	sb.WriteString("## Results\n\n")
	for _, r := range results {
		fmt.Fprintf(&sb, "### %s\n\n", r.Name)
		fmt.Fprintf(&sb, "**Spec:** `%s`\n\n", r.SpecPath)
		fmt.Fprintf(&sb, "**Duration:** %s\n\n", r.Elapsed.Round(time.Millisecond))
		if r.Err != nil {
			sb.WriteString("**Status:** ❌ Failed\n\n")
			fmt.Fprintf(&sb, "**Error:**\n```\n%v\n```\n\n", r.Err)
		} else {
			sb.WriteString("**Status:** ✅ Success\n\n")
			if r.Title != "" {
				fmt.Fprintf(&sb, "**Title:** %s\n\n", r.Title)
			}
			fmt.Fprintf(&sb, "**Endpoints Found:** %d\n\n", r.Endpoints)
		}
		sb.WriteString("---\n\n")
	}
	return sb.String()
}
