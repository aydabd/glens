package generator

import (
	"time"

	"glens/internal/parser"
)

// TestGenerator handles test code generation and execution
type TestGenerator struct {
	framework string
	timeout   time.Duration
}

// ExecutionResult contains the results of test execution
type ExecutionResult struct {
	Passed       bool          `json:"passed"`
	Failed       bool          `json:"failed"`
	Skipped      bool          `json:"skipped"`
	Duration     time.Duration `json:"duration"`
	TestCount    int           `json:"test_count"`
	FailureCount int           `json:"failure_count"`
	ErrorCount   int           `json:"error_count"`
	Output       string        `json:"output"`
	Errors       []TestError   `json:"errors,omitempty"`
	Coverage     *Coverage     `json:"coverage,omitempty"`
	Performance  *Performance  `json:"performance,omitempty"`
}

// TestError represents a test execution error
type TestError struct {
	TestName string `json:"test_name"`
	Message  string `json:"message"`
	Stack    string `json:"stack,omitempty"`
	Type     string `json:"type"` // failure, error, panic
}

// Coverage represents test coverage information
type Coverage struct {
	LinesTotal   int     `json:"lines_total"`
	LinesCovered int     `json:"lines_covered"`
	Percentage   float64 `json:"percentage"`
	Functions    int     `json:"functions"`
	Branches     int     `json:"branches"`
}

// Performance represents performance metrics
type Performance struct {
	MinDuration    time.Duration `json:"min_duration"`
	MaxDuration    time.Duration `json:"max_duration"`
	AvgDuration    time.Duration `json:"avg_duration"`
	MedianDuration time.Duration `json:"median_duration"`
	TotalDuration  time.Duration `json:"total_duration"`
	RequestsCount  int           `json:"requests_count"`
	MemoryUsage    int64         `json:"memory_usage,omitempty"`
}

// TestFile represents a generated test file
type TestFile struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Content     string            `json:"content"`
	Framework   string            `json:"framework"`
	Endpoint    parser.Endpoint   `json:"endpoint"`
	GeneratedAt time.Time         `json:"generated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TestSuite represents a collection of test files
type TestSuite struct {
	Name        string     `json:"name"`
	Files       []TestFile `json:"files"`
	TotalTests  int        `json:"total_tests"`
	Framework   string     `json:"framework"`
	GeneratedAt time.Time  `json:"generated_at"`
}

// Framework represents supported test frameworks
type Framework string

// Framework constants define the available test frameworks
const (
	// FrameworkTestify represents the testify testing framework
	FrameworkTestify Framework = "testify"
	// FrameworkGinkgo represents the Ginkgo BDD testing framework
	FrameworkGinkgo Framework = "ginkgo"
	// FrameworkStandard represents the standard Go testing framework
	FrameworkStandard Framework = "standard"
)

// TestCategory represents different types of tests
type TestCategory string

// Test category constants define the different types of tests that can be generated
const (
	CategoryHappyPath   TestCategory = "happy-path"
	CategoryErrorHandle TestCategory = "error-handling"
	CategoryBoundary    TestCategory = "boundary"
	CategorySecurity    TestCategory = "security"
	CategoryPerformance TestCategory = "performance"
	CategoryIntegration TestCategory = "integration"
)
