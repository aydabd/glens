package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"glens/tools/glens/internal/parser"
)

// GenerateReport creates a comprehensive report from specification and results
func GenerateReport(spec *parser.OpenAPISpec, endpointResults []EndpointResult) *Report {
	log.Info().
		Int("endpoints", len(endpointResults)).
		Msg("Generating comprehensive report")

	startTime := time.Now()

	report := &Report{
		Specification:   *spec,
		EndpointResults: endpointResults,
		GeneratedAt:     time.Now(),
		Metadata:        make(map[string]interface{}),
	}

	// Generate summary
	report.Summary = generateSummary(spec, endpointResults)

	// Generate model comparison
	report.ModelComparison = generateModelComparison(endpointResults)

	// Calculate overall execution time
	report.ExecutionTime = time.Since(startTime)

	// Add metadata
	report.Metadata["report_version"] = "1.0.0"
	report.Metadata["generator"] = "glens"
	report.Metadata["total_endpoints"] = len(spec.Endpoints)
	report.Metadata["processed_endpoints"] = len(endpointResults)

	log.Info().
		Dur("generation_time", report.ExecutionTime).
		Float64("overall_health_score", report.Summary.OverallHealthScore).
		Msg("Report generation completed")

	return report
}

// generateSummary creates the summary section of the report
func generateSummary(spec *parser.OpenAPISpec, results []EndpointResult) Summary {
	summary := Summary{
		TotalEndpoints:     len(spec.Endpoints),
		EndpointsProcessed: len(results),
		AIModelsUsed:       make([]string, 0),
		Frameworks:         make([]string, 0),
	}

	totalTests := 0
	passedTests := 0
	failedTests := 0
	skippedTests := 0
	issuesCreated := 0
	modelsMap := make(map[string]bool)
	frameworksMap := make(map[string]bool)

	var executionTimes []time.Duration
	var generationTimes []time.Duration

	for i := range results {
		result := &results[i]
		if result.IssueNumber > 0 {
			issuesCreated++
		}

		for modelName := range result.Tests {
			testResult := result.Tests[modelName]
			modelsMap[modelName] = true
			frameworksMap[testResult.Framework] = true

			totalTests++

			if testResult.ExecutionResult != nil {
				if testResult.ExecutionResult.Passed {
					passedTests++
				} else if testResult.ExecutionResult.Failed {
					failedTests++
				}
				if testResult.ExecutionResult.Skipped {
					skippedTests++
				}

				executionTimes = append(executionTimes, testResult.ExecutionResult.Duration)
			}

			if testResult.Metrics.Performance.GenerationTime > 0 {
				generationTimes = append(generationTimes, testResult.Metrics.Performance.GenerationTime)
			}
		}
	}

	// Convert maps to slices
	for model := range modelsMap {
		summary.AIModelsUsed = append(summary.AIModelsUsed, model)
	}
	for framework := range frameworksMap {
		summary.Frameworks = append(summary.Frameworks, framework)
	}

	summary.TotalTests = totalTests
	summary.PassedTests = passedTests
	summary.FailedTests = failedTests
	summary.SkippedTests = skippedTests
	summary.TotalIssuesCreated = issuesCreated

	// Calculate execution summary
	summary.ExecutionSummary = calculateExecutionSummary(executionTimes, generationTimes, passedTests, totalTests)

	// Calculate overall health score
	summary.OverallHealthScore = calculateOverallHealthScore(&summary)

	return summary
}

// calculateExecutionSummary calculates timing and performance statistics
func calculateExecutionSummary(executionTimes, generationTimes []time.Duration, passedTests, totalTests int) ExecutionSummary {
	summary := ExecutionSummary{}

	if len(executionTimes) > 0 {
		total := time.Duration(0)
		fastest := executionTimes[0]
		slowest := executionTimes[0]

		for _, duration := range executionTimes {
			total += duration
			if duration < fastest {
				fastest = duration
			}
			if duration > slowest {
				slowest = duration
			}
		}

		summary.TotalDuration = total
		summary.AverageTestTime = total / time.Duration(len(executionTimes))
		summary.FastestTest = fastest
		summary.SlowestTest = slowest
	}

	if len(generationTimes) > 0 {
		total := time.Duration(0)
		for _, duration := range generationTimes {
			total += duration
		}
		summary.GenerationTime = total
	}

	if totalTests > 0 {
		summary.SuccessRate = float64(passedTests) / float64(totalTests)
	}

	return summary
}

// calculateOverallHealthScore calculates a composite health score
func calculateOverallHealthScore(summary *Summary) float64 {
	if summary.TotalTests == 0 {
		return 0.0
	}

	// Calculate success rate
	successRate := float64(summary.PassedTests) / float64(summary.TotalTests)

	// Calculate coverage (endpoints processed vs total)
	coverageRate := float64(summary.EndpointsProcessed) / float64(summary.TotalEndpoints)

	// Weighted score (70% success rate, 30% coverage)
	healthScore := (successRate * 0.7) + (coverageRate * 0.3)

	return healthScore * 100 // Return as percentage
}

// generateModelComparison creates the model comparison section
func generateModelComparison(results []EndpointResult) ModelComparison {
	comparison := ModelComparison{
		Models: make([]ModelResult, 0),
		ComparisonMatrix: ComparisonMatrix{
			QualityComparison:     make(map[string]float64),
			CoverageComparison:    make(map[string]float64),
			PerformanceComparison: make(map[string]float64),
			SecurityComparison:    make(map[string]float64),
			ReliabilityComparison: make(map[string]float64),
		},
		Recommendations: make([]Recommendation, 0),
		Rankings:        make([]ModelRanking, 0),
	}

	// Aggregate results by model
	modelStats := make(map[string]*ModelResult)

	for i := range results {
		result := &results[i]
		for modelName := range result.Tests {
			testResult := result.Tests[modelName]
			if _, exists := modelStats[modelName]; !exists {
				modelStats[modelName] = &ModelResult{
					ModelName:  modelName,
					Strengths:  make([]string, 0),
					Weaknesses: make([]string, 0),
				}
			}

			stats := modelStats[modelName]
			stats.TestsGenerated++

			if testResult.ExecutionResult != nil {
				if testResult.ExecutionResult.Passed {
					stats.TestsPassed++
				} else {
					stats.TestsFailed++
				}

				stats.AvgExecutionTime += testResult.ExecutionResult.Duration
			}

			stats.AvgQualityScore += testResult.QualityScore
			stats.AvgCoverageScore += testResult.Metrics.TestCoverage.CoveragePercentage
			stats.TotalTokensUsed += testResult.Metrics.Performance.TokensUsed
		}
	}

	// Calculate averages and finalize stats
	for modelName, stats := range modelStats {
		if stats.TestsGenerated > 0 {
			stats.AvgQualityScore /= float64(stats.TestsGenerated)
			stats.AvgCoverageScore /= float64(stats.TestsGenerated)
			stats.AvgExecutionTime /= time.Duration(stats.TestsGenerated)
			stats.SuccessRate = float64(stats.TestsPassed) / float64(stats.TestsGenerated)
		}

		// Identify strengths and weaknesses
		stats.Strengths, stats.Weaknesses = identifyModelCharacteristics(stats)

		comparison.Models = append(comparison.Models, *stats)

		// Populate comparison matrix
		comparison.ComparisonMatrix.QualityComparison[modelName] = stats.AvgQualityScore
		comparison.ComparisonMatrix.CoverageComparison[modelName] = stats.AvgCoverageScore
		comparison.ComparisonMatrix.PerformanceComparison[modelName] = float64(stats.AvgExecutionTime.Milliseconds())
		comparison.ComparisonMatrix.ReliabilityComparison[modelName] = stats.SuccessRate
	}

	// Generate rankings
	comparison.Rankings = generateRankings(comparison.Models)

	// Determine best performer
	comparison.BestPerformer = determineBestPerformer(comparison.Models)

	// Generate recommendations
	comparison.Recommendations = generateRecommendations(comparison.Models)

	return comparison
}

// identifyModelCharacteristics identifies strengths and weaknesses of each model
func identifyModelCharacteristics(model *ModelResult) (strengths, weaknesses []string) {
	strengths = make([]string, 0)
	weaknesses = make([]string, 0)

	// Quality assessment
	if model.AvgQualityScore > 80 {
		strengths = append(strengths, "High code quality")
	} else if model.AvgQualityScore < 60 {
		weaknesses = append(weaknesses, "Low code quality")
	}

	// Coverage assessment
	if model.AvgCoverageScore > 85 {
		strengths = append(strengths, "Excellent test coverage")
	} else if model.AvgCoverageScore < 70 {
		weaknesses = append(weaknesses, "Limited test coverage")
	}

	// Performance assessment
	if model.AvgExecutionTime < 5*time.Second {
		strengths = append(strengths, "Fast test execution")
	} else if model.AvgExecutionTime > 15*time.Second {
		weaknesses = append(weaknesses, "Slow test execution")
	}

	// Reliability assessment
	if model.SuccessRate > 0.9 {
		strengths = append(strengths, "High reliability")
	} else if model.SuccessRate < 0.7 {
		weaknesses = append(weaknesses, "Low reliability")
	}

	// Token efficiency
	avgTokensPerTest := float64(model.TotalTokensUsed) / float64(model.TestsGenerated)
	if avgTokensPerTest < 2000 {
		strengths = append(strengths, "Token efficient")
	} else if avgTokensPerTest > 4000 {
		weaknesses = append(weaknesses, "High token usage")
	}

	return strengths, weaknesses
}

// generateRankings creates rankings for different criteria
func generateRankings(models []ModelResult) []ModelRanking {
	rankings := make([]ModelRanking, 0)

	// Quality ranking
	qualityRanking := ModelRanking{
		Criteria: "Code Quality",
		Rankings: make([]RankingEntry, 0),
	}

	// Sort models by quality score
	sortedByQuality := make([]ModelResult, len(models))
	copy(sortedByQuality, models)
	sort.Slice(sortedByQuality, func(i, j int) bool {
		return sortedByQuality[i].AvgQualityScore > sortedByQuality[j].AvgQualityScore
	})

	for i := range sortedByQuality {
		model := &sortedByQuality[i]
		qualityRanking.Rankings = append(qualityRanking.Rankings, RankingEntry{
			Rank:  i + 1,
			Model: model.ModelName,
			Score: model.AvgQualityScore,
		})
	}
	rankings = append(rankings, qualityRanking)

	// Coverage ranking
	coverageRanking := ModelRanking{
		Criteria: "Test Coverage",
		Rankings: make([]RankingEntry, 0),
	}

	sortedByCoverage := make([]ModelResult, len(models))
	copy(sortedByCoverage, models)
	sort.Slice(sortedByCoverage, func(i, j int) bool {
		return sortedByCoverage[i].AvgCoverageScore > sortedByCoverage[j].AvgCoverageScore
	})

	for i := range sortedByCoverage {
		model := &sortedByCoverage[i]
		coverageRanking.Rankings = append(coverageRanking.Rankings, RankingEntry{
			Rank:  i + 1,
			Model: model.ModelName,
			Score: model.AvgCoverageScore,
		})
	}
	rankings = append(rankings, coverageRanking)

	// Reliability ranking
	reliabilityRanking := ModelRanking{
		Criteria: "Reliability",
		Rankings: make([]RankingEntry, 0),
	}

	sortedByReliability := make([]ModelResult, len(models))
	copy(sortedByReliability, models)
	sort.Slice(sortedByReliability, func(i, j int) bool {
		return sortedByReliability[i].SuccessRate > sortedByReliability[j].SuccessRate
	})

	for i := range sortedByReliability {
		model := &sortedByReliability[i]
		reliabilityRanking.Rankings = append(reliabilityRanking.Rankings, RankingEntry{
			Rank:  i + 1,
			Model: model.ModelName,
			Score: model.SuccessRate * 100,
		})
	}
	rankings = append(rankings, reliabilityRanking)

	return rankings
}

// determineBestPerformer identifies the overall best performing model
func determineBestPerformer(models []ModelResult) string {
	if len(models) == 0 {
		return ""
	}

	bestModel := models[0]
	bestScore := calculateCompositeScore(&bestModel)

	for i := 1; i < len(models); i++ {
		model := &models[i]
		score := calculateCompositeScore(model)
		if score > bestScore {
			bestScore = score
			bestModel = *model
		}
	}

	return bestModel.ModelName
}

// calculateCompositeScore calculates a weighted composite score for ranking
func calculateCompositeScore(model *ModelResult) float64 {
	// Weighted scoring: 30% quality, 25% coverage, 25% reliability, 20% performance
	qualityWeight := 0.30
	coverageWeight := 0.25
	reliabilityWeight := 0.25
	performanceWeight := 0.20

	// Normalize performance score (lower execution time is better)
	performanceScore := 100.0
	if model.AvgExecutionTime > 0 {
		// Convert to seconds and invert (max 100 for under 1 second)
		seconds := model.AvgExecutionTime.Seconds()
		performanceScore = 100.0 / (1.0 + seconds)
	}

	compositeScore := (model.AvgQualityScore * qualityWeight) +
		(model.AvgCoverageScore * coverageWeight) +
		(model.SuccessRate * 100 * reliabilityWeight) +
		(performanceScore * performanceWeight)

	return compositeScore
}

// generateRecommendations creates actionable recommendations
func generateRecommendations(models []ModelResult) []Recommendation {
	recommendations := make([]Recommendation, 0)

	// Analyze overall performance
	if len(models) > 1 {
		// Find best performer
		best := models[0]
		bestScore := calculateCompositeScore(&best)

		for i := 1; i < len(models); i++ {
			model := &models[i]
			score := calculateCompositeScore(model)
			if score > bestScore {
				bestScore = score
				best = *model
			}
		}

		recommendations = append(recommendations, Recommendation{
			Category:    "Model Selection",
			Title:       "Primary Model Recommendation",
			Description: fmt.Sprintf("Use %s as the primary model for test generation based on overall performance", best.ModelName),
			Priority:    "high",
			ActionItems: []string{
				fmt.Sprintf("Configure %s as the default model", best.ModelName),
				"Monitor performance metrics regularly",
				"Consider cost implications of model choice",
			},
		})

		// Quality recommendations
		avgQuality := 0.0
		for i := range models {
			avgQuality += models[i].AvgQualityScore
		}
		avgQuality /= float64(len(models))

		if avgQuality < 75 {
			recommendations = append(recommendations, Recommendation{
				Category:    "Code Quality",
				Title:       "Improve Test Code Quality",
				Description: "Overall test code quality is below acceptable threshold",
				Priority:    "medium",
				ActionItems: []string{
					"Review and refine AI prompts for better code generation",
					"Implement code quality checks in the pipeline",
					"Consider post-processing to improve generated code",
				},
			})
		}
	}

	// Performance recommendations
	for i := range models {
		model := &models[i]
		if model.AvgExecutionTime > 30*time.Second {
			recommendations = append(recommendations, Recommendation{
				Category:    "Performance",
				Title:       fmt.Sprintf("Optimize %s Performance", model.ModelName),
				Description: "Test execution time is higher than expected",
				Priority:    "medium",
				ActionItems: []string{
					"Review test complexity and reduce if possible",
					"Implement parallel test execution",
					"Consider timeout optimizations",
				},
			})
		}
	}

	return recommendations
}

// WriteReport writes the report to a file in the specified format
func WriteReport(report *Report, filePath string) error {
	log.Info().
		Str("file_path", filePath).
		Msg("Writing report to file")

	// Determine format from file extension
	format := FormatJSON
	if strings.HasSuffix(strings.ToLower(filePath), ".md") {
		format = FormatMarkdown
	} else if strings.HasSuffix(strings.ToLower(filePath), ".html") {
		format = FormatHTML
	}

	var content string
	var err error

	switch format {
	case FormatMarkdown:
		content, err = generateMarkdownReport(report)
	case FormatHTML:
		content, err = generateHTMLReport(report)
	default:
		// JSON format
		jsonData, jsonErr := json.MarshalIndent(report, "", "  ")
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal report to JSON: %w", jsonErr)
		}
		content = string(jsonData)
	}

	if err != nil {
		return fmt.Errorf("failed to generate report content: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	log.Info().
		Str("file_path", filePath).
		Str("format", string(format)).
		Int("size_bytes", len(content)).
		Msg("Report written successfully")

	return nil
}

// EnsureReportDirectory ensures the directory for the report file exists
func EnsureReportDirectory(filePath string) error {
	// Extract directory from file path
	var dir string
	if lastSlash := strings.LastIndex(filePath, "/"); lastSlash != -1 {
		dir = filePath[:lastSlash]
	} else {
		// No directory in path, current directory is fine
		return nil
	}

	// Create directory if it doesn't exist (0o750 for security)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	log.Debug().
		Str("directory", dir).
		Msg("Report directory ensured")

	return nil
}
