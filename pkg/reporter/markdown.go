package reporter

import (
	"fmt"
	"strings"
	"time"

	"glens/pkg/parser"
)

// fixMarkdownListSpacing ensures lists in markdown text have proper blank lines
func fixMarkdownListSpacing(text string) string {
	lines := strings.Split(text, "\n")
	var result []string

	for i, line := range lines {
		// Check if current line is a list item
		trimmed := strings.TrimSpace(line)
		isListItem := strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ")

		// Check if previous line exists and is not a list item or blank
		if i > 0 && isListItem {
			prevLine := strings.TrimSpace(lines[i-1])
			prevIsListItem := strings.HasPrefix(prevLine, "- ") || strings.HasPrefix(prevLine, "* ")
			prevIsBlank := prevLine == ""

			// Add blank line before list if previous line is not blank and not a list item
			if !prevIsBlank && !prevIsListItem && prevLine != "" {
				result = append(result, "")
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// generateMarkdownReport creates a markdown formatted report
func generateMarkdownReport(report *Report) (string, error) {
	var md strings.Builder

	// Header
	fmt.Fprintf(&md, "# OpenAPI Integration Test Report\n\n")
	fmt.Fprintf(&md, "**Generated:** %s\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&md, "**Execution Time:** %s\n", report.ExecutionTime)
	fmt.Fprintf(&md, "**API:** %s v%s\n\n", report.Specification.Info.Title, report.Specification.Info.Version)

	// Executive Summary
	fmt.Fprintf(&md, "## 📊 Executive Summary\n\n")
	writeExecutiveSummary(&md, &report.Summary)

	// API Specification Overview
	fmt.Fprintf(&md, "## 📋 API Specification\n\n")
	writeSpecificationOverview(&md, &report.Specification)

	// Model Performance Comparison
	fmt.Fprintf(&md, "## 🤖 AI Model Performance Comparison\n\n")
	writeModelComparison(&md, &report.ModelComparison)

	// Detailed Endpoint Results
	fmt.Fprintf(&md, "## 🎯 Endpoint Test Results\n\n")
	writeEndpointResults(&md, report.EndpointResults)

	// Recommendations
	if len(report.ModelComparison.Recommendations) > 0 {
		fmt.Fprintf(&md, "## 💡 Recommendations\n\n")
		writeRecommendations(&md, report.ModelComparison.Recommendations)
	}

	// Appendices
	fmt.Fprintf(&md, "## 📎 Appendices\n\n")
	writeAppendices(&md, report)

	return md.String(), nil
}

// writeExecutiveSummary writes the executive summary section
func writeExecutiveSummary(md *strings.Builder, summary *Summary) {
	fmt.Fprintf(md, "| Metric | Value |\n")
	fmt.Fprintf(md, "|--------|-------|\n")
	fmt.Fprintf(md, "| **Total Endpoints** | %d |\n", summary.TotalEndpoints)
	fmt.Fprintf(md, "| **Endpoints Processed** | %d |\n", summary.EndpointsProcessed)
	fmt.Fprintf(md, "| **Total Tests Generated** | %d |\n", summary.TotalTests)
	fmt.Fprintf(md, "| **Tests Passed** | %d ✅ |\n", summary.PassedTests)
	fmt.Fprintf(md, "| **Tests Failed** | %d ❌ |\n", summary.FailedTests)
	fmt.Fprintf(md, "| **Tests Skipped** | %d ⏭️ |\n", summary.SkippedTests)
	fmt.Fprintf(md, "| **GitHub Issues Created** | %d |\n", summary.TotalIssuesCreated)
	fmt.Fprintf(md, "| **AI Models Used** | %s |\n", strings.Join(summary.AIModelsUsed, ", "))
	fmt.Fprintf(md, "| **Overall Health Score** | %.1f%% |\n", summary.OverallHealthScore)

	// Health Score Badge
	healthEmoji := "🟢"
	if summary.OverallHealthScore < 70 {
		healthEmoji = "🟡"
	}
	if summary.OverallHealthScore < 50 {
		healthEmoji = "🔴"
	}

	fmt.Fprintf(md, "\n### Overall Health Status\n\n")
	fmt.Fprintf(md, "%s **%.1f%%** - ", healthEmoji, summary.OverallHealthScore)

	switch {
	case summary.OverallHealthScore >= 80:
		fmt.Fprintf(md, "Excellent API test coverage and quality")
	case summary.OverallHealthScore >= 70:
		fmt.Fprintf(md, "Good API test coverage with room for improvement")
	case summary.OverallHealthScore >= 50:
		fmt.Fprintf(md, "Moderate API test coverage - requires attention")
	default:
		fmt.Fprintf(md, "Poor API test coverage - immediate action required")
	}

	fmt.Fprintf(md, "\n\n### Performance Summary\n\n")
	fmt.Fprintf(md, "| Metric | Value |\n")
	fmt.Fprintf(md, "|--------|-------|\n")
	fmt.Fprintf(md, "| **Total Execution Time** | %s |\n", summary.ExecutionSummary.TotalDuration)
	fmt.Fprintf(md, "| **Average Test Time** | %s |\n", summary.ExecutionSummary.AverageTestTime)
	fmt.Fprintf(md, "| **Fastest Test** | %s |\n", summary.ExecutionSummary.FastestTest)
	fmt.Fprintf(md, "| **Slowest Test** | %s |\n", summary.ExecutionSummary.SlowestTest)
	fmt.Fprintf(md, "| **Success Rate** | %.1f%% |\n", summary.ExecutionSummary.SuccessRate*100)

	fmt.Fprintf(md, "\n")
}

// writeSpecificationOverview writes the API specification overview
func writeSpecificationOverview(md *strings.Builder, spec *parser.OpenAPISpec) {
	fmt.Fprintf(md, "**Title:** %s\n", spec.Info.Title)
	fmt.Fprintf(md, "**Version:** %s\n", spec.Info.Version)
	fmt.Fprintf(md, "**OpenAPI Version:** %s\n", spec.Version)

	if spec.Info.Description != "" {
		// Fix list spacing in description to ensure markdown linting compliance
		fixedDesc := fixMarkdownListSpacing(spec.Info.Description)
		fmt.Fprintf(md, "**Description:** %s\n\n", fixedDesc)
	}

	fmt.Fprintf(md, "**Total Endpoints:** %d\n", len(spec.Endpoints))

	// Server information
	if len(spec.Servers) > 0 {
		fmt.Fprintf(md, "\n### Servers\n")
		for _, server := range spec.Servers {
			fmt.Fprintf(md, "\n- **%s**", server.URL)
			if server.Description != "" {
				if server.Description != "" {
					fmt.Fprintf(md, " - %s", server.Description)
				}
			}
			fmt.Fprintf(md, "\n")
		}
	}

	// Endpoint breakdown by method
	methodCounts := make(map[string]int)
	for i := range spec.Endpoints {
		methodCounts[spec.Endpoints[i].Method]++
	}

	if len(methodCounts) > 0 {
		fmt.Fprintf(md, "\n### Endpoint Breakdown\n\n")
		fmt.Fprintf(md, "| HTTP Method | Count |\n")
		fmt.Fprintf(md, "|-------------|-------|\n")
		for method, count := range methodCounts {
			fmt.Fprintf(md, "| %s | %d |\n", method, count)
		}
	}

	fmt.Fprintf(md, "\n")
}

// writeModelComparison writes the AI model comparison section
func writeModelComparison(md *strings.Builder, comparison *ModelComparison) {
	if len(comparison.Models) == 0 {
		fmt.Fprintf(md, "No model results available.\n\n")
		return
	}

	fmt.Fprintf(md, "**Best Performer:** %s 🏆\n\n", comparison.BestPerformer)

	// Overall comparison table
	fmt.Fprintf(md, "### Model Performance Overview\n\n")
	fmt.Fprintf(md, "| Model | Tests Generated | Success Rate | Avg Quality | Avg Coverage | Avg Execution Time |\n")
	fmt.Fprintf(md, "|-------|----------------|--------------|-------------|--------------|-------------------|\n")

	for i := range comparison.Models {
		model := &comparison.Models[i]
		fmt.Fprintf(md, "| **%s** | %d | %.1f%% | %.1f | %.1f%% | %s |\n",
			model.ModelName,
			model.TestsGenerated,
			model.SuccessRate*100,
			model.AvgQualityScore,
			model.AvgCoverageScore,
			model.AvgExecutionTime)
	}

	// Rankings
	if len(comparison.Rankings) > 0 {
		fmt.Fprintf(md, "\n### Performance Rankings\n\n")
		for _, ranking := range comparison.Rankings {
			fmt.Fprintf(md, "#### %s\n\n", ranking.Criteria)
			fmt.Fprintf(md, "| Rank | Model | Score |\n")
			fmt.Fprintf(md, "|------|-------|-------|\n")

			for _, entry := range ranking.Rankings {
				medal := ""
				switch entry.Rank {
				case 1:
					medal = "🥇"
				case 2:
					medal = "🥈"
				case 3:
					medal = "🥉"
				}
				// Build rank cell without extra space when no medal.
				rankCell := fmt.Sprintf("%d", entry.Rank)
				if medal != "" {
					rankCell += " " + medal
				}
				fmt.Fprintf(md, "| %s | %s | %.1f |\n", rankCell, entry.Model, entry.Score)
			}
			fmt.Fprintf(md, "\n")
		}
	}

	// Detailed model analysis
	fmt.Fprintf(md, "### Detailed Model Analysis\n\n")
	fmt.Fprintf(md, "Total Models Evaluated: %d\n\n", len(comparison.Models))
	for i := range comparison.Models {
		model := &comparison.Models[i]
		fmt.Fprintf(md, "#### %s\n\n", model.ModelName)

		// Strengths and weaknesses
		if len(model.Strengths) > 0 {
			fmt.Fprintf(md, "**Strengths:**\n\n")
			for _, strength := range model.Strengths {
				fmt.Fprintf(md, "- ✅ %s\n", strength)
			}
			fmt.Fprintf(md, "\n")
		}

		if len(model.Weaknesses) > 0 {
			fmt.Fprintf(md, "**Weaknesses:**\n\n")
			for _, weakness := range model.Weaknesses {
				fmt.Fprintf(md, "- ⚠️ %s\n", weakness)
			}
			fmt.Fprintf(md, "\n")
		}

		// Detailed metrics
		fmt.Fprintf(md, "**Metrics:**\n\n")
		fmt.Fprintf(md, "- Tests Generated: %d\n", model.TestsGenerated)
		fmt.Fprintf(md, "- Tests Passed: %d\n", model.TestsPassed)
		fmt.Fprintf(md, "- Tests Failed: %d\n", model.TestsFailed)
		fmt.Fprintf(md, "- Success Rate: %.1f%%\n", model.SuccessRate*100)
		fmt.Fprintf(md, "- Average Quality Score: %.1f\n", model.AvgQualityScore)
		fmt.Fprintf(md, "- Average Coverage: %.1f%%\n", model.AvgCoverageScore)
		fmt.Fprintf(md, "- Average Execution Time: %s\n", model.AvgExecutionTime)
		fmt.Fprintf(md, "- Total Tokens Used: %d\n", model.TotalTokensUsed)

		fmt.Fprintf(md, "\n")
	}
}

// writeEndpointResults writes the detailed endpoint results
func writeEndpointResults(md *strings.Builder, results []EndpointResult) {
	if len(results) == 0 {
		fmt.Fprintf(md, "No endpoint results available.\n\n")
		return
	}

	fmt.Fprintf(md, "### Summary\n\n")
	fmt.Fprintf(md, "| Endpoint | Status | Issue | Tests | Passed | Failed | Overall Score |\n")
	fmt.Fprintf(md, "|----------|--------|-------|-------|--------|--------|--------------|\n")

	for i := range results {
		result := &results[i]
		statusEmoji := getStatusEmoji(result.Status)
		issueLink := ""
		if result.IssueNumber > 0 {
			issueLink = fmt.Sprintf("#%d", result.IssueNumber)
		}

		testCount := len(result.Tests)
		passedCount := 0
		failedCount := 0

		for modelName := range result.Tests {
			if result.Tests[modelName].ExecutionResult != nil {
				if result.Tests[modelName].ExecutionResult.Passed {
					passedCount++
				} else {
					failedCount++
				}
			}
		}

		fmt.Fprintf(md, "| `%s %s` | %s %s | %s | %d | %d | %d | %.1f |\n",
			result.Endpoint.Method,
			result.Endpoint.Path,
			statusEmoji,
			result.Status,
			issueLink,
			testCount,
			passedCount,
			failedCount,
			result.OverallScore)
	}

	// Detailed results for each endpoint
	fmt.Fprintf(md, "\n### Detailed Results\n\n")
	for i := range results {
		result := &results[i]
		fmt.Fprintf(md, "#### %d. %s %s\n\n", i+1, result.Endpoint.Method, result.Endpoint.Path)

		if result.Endpoint.Summary != "" {
			fmt.Fprintf(md, "**Summary:** %s\n\n", result.Endpoint.Summary)
		}

		if result.IssueNumber > 0 {
			fmt.Fprintf(md, "**GitHub Issue:** #%d\n\n", result.IssueNumber)
		}

		fmt.Fprintf(md, "**Test Results by Model:**\n\n")
		for modelName := range result.Tests {
			test := result.Tests[modelName]
			fmt.Fprintf(md, "##### Model: %s\n\n", modelName)

			if test.ExecutionResult != nil {
				status := "✅ Passed"
				if test.ExecutionResult.Failed {
					status = "❌ Failed"
				} else if test.ExecutionResult.Skipped {
					status = "⏭️ Skipped"
				}

				fmt.Fprintf(md, "- **Status:** %s\n", status)
				fmt.Fprintf(md, "- **Duration:** %s\n", test.ExecutionResult.Duration)
				fmt.Fprintf(md, "- **Test Count:** %d\n", test.ExecutionResult.TestCount)

				if len(test.ExecutionResult.Errors) > 0 {
					fmt.Fprintf(md, "- **Errors:**\n")
					for _, err := range test.ExecutionResult.Errors {
						if err.Message != "" {
							fmt.Fprintf(md, "  - %s: %s\n", err.TestName, err.Message)
						} else {
							fmt.Fprintf(md, "  - %s\n", err.TestName)
						}
					}
				}
			} else if test.ExecutionError != "" {
				fmt.Fprintf(md, "- **Status:** ❌ Execution Error\n")
				fmt.Fprintf(md, "- **Error:** %s\n", test.ExecutionError)
			}

			fmt.Fprintf(md, "- **Quality Score:** %.1f\n", test.QualityScore)
			fmt.Fprintf(md, "- **Framework:** %s\n", test.Framework)
			fmt.Fprintf(md, "- **Generated At:** %s\n", test.GeneratedAt.Format(time.RFC3339))

			fmt.Fprintf(md, "\n")
		}

		fmt.Fprintf(md, "---\n\n")
	}
}

// writeRecommendations writes the recommendations section
func writeRecommendations(md *strings.Builder, recommendations []Recommendation) {
	for _, rec := range recommendations {
		priorityEmoji := "📌"
		switch rec.Priority {
		case "high":
			priorityEmoji = "🔴"
		case "medium":
			priorityEmoji = "🟡"
		case "low":
			priorityEmoji = "🟢"
		}

		fmt.Fprintf(md, "### %s %s\n\n", priorityEmoji, rec.Title)
		fmt.Fprintf(md, "**Category:** %s\n\n", rec.Category)
		fmt.Fprintf(md, "**Priority:** %s\n\n", strings.ToUpper(rec.Priority))
		fmt.Fprintf(md, "**Description:** %s\n\n", rec.Description)

		if len(rec.ActionItems) > 0 {
			fmt.Fprintf(md, "**Action Items:**\n")
			for _, item := range rec.ActionItems {
				fmt.Fprintf(md, "- [ ] %s\n", item)
			}
			fmt.Fprintf(md, "\n")
		}
	}
}

// writeAppendices writes the appendices section
func writeAppendices(md *strings.Builder, report *Report) {
	fmt.Fprintf(md, "### A. Metadata\n\n")
	fmt.Fprintf(md, "| Key | Value |\n")
	fmt.Fprintf(md, "|-----|-------|\n")
	for key, value := range report.Metadata {
		fmt.Fprintf(md, "| %s | %v |\n", key, value)
	}

	fmt.Fprintf(md, "\n### B. Test Execution Environment\n\n")
	fmt.Fprintf(md, "- **Test Framework:** Go with testify\n")
	fmt.Fprintf(md, "- **Execution Mode:** Sequential\n")
	fmt.Fprintf(md, "- **Timeout:** 2 minutes per test\n")
	fmt.Fprintf(md, "- **Report Generated:** %s\n\n", report.GeneratedAt.Format(time.RFC3339))

	fmt.Fprintf(md, "---\n\n")
	fmt.Fprintf(md, "This report was automatically generated by Glens\n")
}

// getStatusEmoji returns an emoji for the endpoint status
func getStatusEmoji(status EndpointStatus) string {
	switch status {
	case StatusCompleted:
		return "✅"
	case StatusFailed:
		return "❌"
	case StatusProcessing:
		return "⏳"
	case StatusSkipped:
		return "⏭️"
	default:
		return "⏸️"
	}
}
