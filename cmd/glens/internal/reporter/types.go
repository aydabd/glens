package reporter

import (
	"time"

	"glens/tools/glens/internal/generator"
	"glens/tools/glens/internal/parser"
)

// Report represents the final comprehensive report
type Report struct {
	Summary         Summary                `json:"summary"`
	Specification   parser.OpenAPISpec     `json:"specification"`
	EndpointResults []EndpointResult       `json:"endpoint_results"`
	ModelComparison ModelComparison        `json:"model_comparison"`
	GeneratedAt     time.Time              `json:"generated_at"`
	ExecutionTime   time.Duration          `json:"execution_time"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// Summary contains high-level statistics
type Summary struct {
	TotalEndpoints     int              `json:"total_endpoints"`
	EndpointsProcessed int              `json:"endpoints_processed"`
	TotalTests         int              `json:"total_tests"`
	PassedTests        int              `json:"passed_tests"`
	FailedTests        int              `json:"failed_tests"`
	SkippedTests       int              `json:"skipped_tests"`
	TotalIssuesCreated int              `json:"total_issues_created"`
	AIModelsUsed       []string         `json:"ai_models_used"`
	Frameworks         []string         `json:"frameworks"`
	ExecutionSummary   ExecutionSummary `json:"execution_summary"`
	OverallHealthScore float64          `json:"overall_health_score"`
}

// ExecutionSummary contains timing and performance data
type ExecutionSummary struct {
	TotalDuration   time.Duration `json:"total_duration"`
	AverageTestTime time.Duration `json:"average_test_time"`
	FastestTest     time.Duration `json:"fastest_test"`
	SlowestTest     time.Duration `json:"slowest_test"`
	GenerationTime  time.Duration `json:"generation_time"`
	ExecutionTime   time.Duration `json:"execution_time"`
	SuccessRate     float64       `json:"success_rate"`
}

// EndpointResult contains results for a specific endpoint
type EndpointResult struct {
	Endpoint     parser.Endpoint       `json:"endpoint"`
	IssueNumber  int                   `json:"issue_number,omitempty"`
	Tests        map[string]TestResult `json:"tests"` // key: AI model name
	OverallScore float64               `json:"overall_score"`
	Status       EndpointStatus        `json:"status"`
	ProcessedAt  time.Time             `json:"processed_at"`
}

// TestResult contains results for a specific AI model's test
type TestResult struct {
	AIModel         string                     `json:"ai_model"`
	Prompt          string                     `json:"prompt"`
	TestCode        string                     `json:"test_code"`
	Framework       string                     `json:"framework"`
	ExecutionResult *generator.ExecutionResult `json:"execution_result,omitempty"`
	ExecutionError  string                     `json:"execution_error,omitempty"`
	GeneratedAt     time.Time                  `json:"generated_at"`
	Metrics         TestMetrics                `json:"metrics"`
	QualityScore    float64                    `json:"quality_score"`
}

// TestMetrics contains detailed test metrics
type TestMetrics struct {
	CodeQuality      CodeQuality        `json:"code_quality"`
	TestCoverage     TestCoverage       `json:"test_coverage"`
	Performance      PerformanceMetrics `json:"performance"`
	SecurityCoverage SecurityCoverage   `json:"security_coverage"`
}

// CodeQuality measures the quality of generated test code
type CodeQuality struct {
	LinesOfCode       int      `json:"lines_of_code"`
	TestFunctionCount int      `json:"test_function_count"`
	AssertionCount    int      `json:"assertion_count"`
	CommentLines      int      `json:"comment_lines"`
	ComplexityScore   float64  `json:"complexity_score"`
	ReadabilityScore  float64  `json:"readability_score"`
	CategoriesCovered []string `json:"categories_covered"`
}

// TestCoverage measures how well the test covers the endpoint
type TestCoverage struct {
	HTTPMethodsCovered   []string `json:"http_methods_covered"`
	StatusCodesCovered   []string `json:"status_codes_covered"`
	ParametersCovered    int      `json:"parameters_covered"`
	ParametersTotal      int      `json:"parameters_total"`
	ResponseTypesCovered int      `json:"response_types_covered"`
	EdgeCasesCovered     int      `json:"edge_cases_covered"`
	CoveragePercentage   float64  `json:"coverage_percentage"`
}

// PerformanceMetrics contains performance-related metrics
type PerformanceMetrics struct {
	GenerationTime  time.Duration `json:"generation_time"`
	ExecutionTime   time.Duration `json:"execution_time"`
	TokensUsed      int           `json:"tokens_used"`
	APICallsCount   int           `json:"api_calls_count"`
	MemoryUsage     int64         `json:"memory_usage,omitempty"`
	ResponseTimesMs []float64     `json:"response_times_ms,omitempty"`
}

// SecurityCoverage measures security test coverage
type SecurityCoverage struct {
	AuthenticationTests  bool     `json:"authentication_tests"`
	AuthorizationTests   bool     `json:"authorization_tests"`
	InputValidationTests bool     `json:"input_validation_tests"`
	SQLInjectionTests    bool     `json:"sql_injection_tests"`
	XSSTests             bool     `json:"xss_tests"`
	SecurityScore        float64  `json:"security_score"`
	VulnerabilitiesFound []string `json:"vulnerabilities_found,omitempty"`
}

// ModelComparison compares results across different AI models
type ModelComparison struct {
	Models           []ModelResult    `json:"models"`
	ComparisonMatrix ComparisonMatrix `json:"comparison_matrix"`
	Recommendations  []Recommendation `json:"recommendations"`
	BestPerformer    string           `json:"best_performer"`
	Rankings         []ModelRanking   `json:"rankings"`
}

// ModelResult contains aggregated results for a specific AI model
type ModelResult struct {
	ModelName        string        `json:"model_name"`
	TestsGenerated   int           `json:"tests_generated"`
	TestsPassed      int           `json:"tests_passed"`
	TestsFailed      int           `json:"tests_failed"`
	AvgQualityScore  float64       `json:"avg_quality_score"`
	AvgCoverageScore float64       `json:"avg_coverage_score"`
	AvgExecutionTime time.Duration `json:"avg_execution_time"`
	TotalTokensUsed  int           `json:"total_tokens_used"`
	SuccessRate      float64       `json:"success_rate"`
	Strengths        []string      `json:"strengths"`
	Weaknesses       []string      `json:"weaknesses"`
}

// ComparisonMatrix provides side-by-side comparison data
type ComparisonMatrix struct {
	QualityComparison     map[string]float64 `json:"quality_comparison"`
	CoverageComparison    map[string]float64 `json:"coverage_comparison"`
	PerformanceComparison map[string]float64 `json:"performance_comparison"`
	SecurityComparison    map[string]float64 `json:"security_comparison"`
	ReliabilityComparison map[string]float64 `json:"reliability_comparison"`
}

// Recommendation provides actionable insights
type Recommendation struct {
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // high, medium, low
	ActionItems []string `json:"action_items"`
}

// ModelRanking ranks models by different criteria
type ModelRanking struct {
	Criteria string         `json:"criteria"`
	Rankings []RankingEntry `json:"rankings"`
}

// RankingEntry represents a single ranking entry
type RankingEntry struct {
	Rank     int     `json:"rank"`
	Model    string  `json:"model"`
	Score    float64 `json:"score"`
	Comments string  `json:"comments,omitempty"`
}

// EndpointStatus represents the processing status of an endpoint
type EndpointStatus string

// Endpoint status constants define the different states an endpoint can be in during processing
const (
	// StatusPending indicates the endpoint is waiting to be processed
	StatusPending EndpointStatus = "pending"
	// StatusProcessing indicates the endpoint is currently being processed
	StatusProcessing EndpointStatus = "processing"
	// StatusCompleted indicates the endpoint has been successfully processed
	StatusCompleted EndpointStatus = "completed"
	// StatusFailed indicates the endpoint processing failed
	StatusFailed EndpointStatus = "failed"
	// StatusSkipped indicates the endpoint was skipped during processing
	StatusSkipped EndpointStatus = "skipped"
)

// ReportFormat represents the output format for reports
type ReportFormat string

// Report format constants define the different output formats available for reports
const (
	// FormatMarkdown generates reports in Markdown format
	FormatMarkdown ReportFormat = "markdown"
	// FormatJSON generates reports in JSON format
	FormatJSON ReportFormat = "json"
	// FormatHTML generates reports in HTML format
	FormatHTML ReportFormat = "html"
	// FormatPDF generates reports in PDF format
	FormatPDF ReportFormat = "pdf"
)
