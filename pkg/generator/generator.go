package generator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"glens/pkg/parser"
)

// NewTestGenerator creates a new test generator
func NewTestGenerator(framework string) *TestGenerator {
	return &TestGenerator{
		framework: framework,
		timeout:   2 * time.Minute,
	}
}

// ExecuteTest executes the generated test code and returns results
func (g *TestGenerator) ExecuteTest(ctx context.Context, testCode string, endpoint *parser.Endpoint) (*ExecutionResult, error) {
	startTime := time.Now()

	log.Debug().
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Str("framework", g.framework).
		Msg("Executing generated test")

	// Create temporary directory for test execution
	tmpDir, err := os.MkdirTemp("", "glens-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tmpDir); removeErr != nil {
			log.Debug().Err(removeErr).Msg("failed to remove temporary directory")
		}
	}()

	// Write test code to file
	testFileName := g.generateTestFileName(endpoint)
	testFilePath := filepath.Join(tmpDir, testFileName)

	if err := os.WriteFile(testFilePath, []byte(testCode), 0o600); err != nil {
		return nil, fmt.Errorf("failed to write test file: %w", err)
	}

	// Create go.mod for the test
	if err := g.createTestModule(tmpDir); err != nil {
		return nil, fmt.Errorf("failed to create test module: %w", err)
	}

	// Run the test
	result, err := g.runTest(ctx, tmpDir, testFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to run test: %w", err)
	}

	result.Duration = time.Since(startTime)

	log.Info().
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Bool("passed", result.Passed).
		Dur("duration", result.Duration).
		Int("test_count", result.TestCount).
		Msg("Test execution completed")

	return result, nil
}

// generateTestFileName creates a standardized test file name
func (g *TestGenerator) generateTestFileName(endpoint *parser.Endpoint) string {
	// Clean path for filename
	path := strings.ReplaceAll(endpoint.Path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	path = strings.Trim(path, "_")

	if path == "" {
		path = "root"
	}

	method := strings.ToLower(endpoint.Method)
	return fmt.Sprintf("%s_%s_test.go", method, path)
}

// createTestModule creates a go.mod file for the test
func (g *TestGenerator) createTestModule(dir string) error {
	goModContent := `module glens-temp

go 1.21

require (
	github.com/stretchr/testify v1.8.4
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.29.0
)
`

	goModPath := filepath.Join(dir, "go.mod")
	return os.WriteFile(goModPath, []byte(goModContent), 0o600)
}

// runTest executes the test using go test command
func (g *TestGenerator) runTest(ctx context.Context, dir, fileName string) (*ExecutionResult, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	// Run go mod tidy first
	tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = dir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		log.Debug().
			Str("output", string(output)).
			Err(err).
			Msg("go mod tidy failed, continuing anyway")
	}

	// Build test command based on framework
	args := g.buildTestCommand(fileName)

	// Validate args to ensure they're safe (gosec G204 mitigation)
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid test command arguments")
	}
	// Validate that first argument is a safe command
	allowedCommands := map[string]bool{
		"test": true,
		"run":  true,
	}
	if !allowedCommands[args[0]] {
		return nil, fmt.Errorf("invalid command: %s", args[0])
	}

	cmd := exec.CommandContext(ctx, "go", args...) //nolint:gosec // args are validated and come from controlled buildTestCommand function
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	result := &ExecutionResult{
		Output: outputStr,
	}

	// Parse test results based on framework
	switch g.framework {
	case "testify", "standard":
		g.parseGoTestOutput(result, outputStr, err)
	case "ginkgo":
		g.parseGinkgoOutput(result, outputStr, err)
	default:
		g.parseGoTestOutput(result, outputStr, err)
	}

	return result, nil
}

// buildTestCommand builds the appropriate test command for the framework
func (g *TestGenerator) buildTestCommand(fileName string) []string {
	switch g.framework {
	case "ginkgo":
		return []string{"run", "github.com/onsi/ginkgo/v2/ginkgo", "-v", "--json-report=results.json"}
	default:
		return []string{"test", "-v", "-json", "./" + strings.TrimSuffix(fileName, ".go")}
	}
}

// parseGoTestOutput parses standard go test output
func (g *TestGenerator) parseGoTestOutput(result *ExecutionResult, output string, cmdErr error) {
	lines := strings.Split(output, "\n")

	testCount := 0
	failureCount := 0
	errorCount := 0
	var errors []TestError

	// Regex patterns for parsing test output
	testRunPattern := regexp.MustCompile(`^=== RUN\s+(\S+)`)
	testFailPattern := regexp.MustCompile(`^--- FAIL:\s+(\S+)\s+\(([0-9.]+)s\)`)
	testSkipPattern := regexp.MustCompile(`^--- SKIP:\s+(\S+)\s+\(([0-9.]+)s\)`)

	for i, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case testRunPattern.MatchString(line):
			testCount++
		case testFailPattern.MatchString(line):
			failureCount++
			matches := testFailPattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				testName := matches[1]

				// Look for error message in following lines
				errorMsg := ""
				for j := i + 1; j < len(lines) && j < i+10; j++ {
					if strings.HasPrefix(strings.TrimSpace(lines[j]), "---") {
						break
					}
					if strings.TrimSpace(lines[j]) != "" {
						errorMsg += lines[j] + "\n"
					}
				}

				errors = append(errors, TestError{
					TestName: testName,
					Message:  strings.TrimSpace(errorMsg),
					Type:     "failure",
				})
			}
		case testSkipPattern.MatchString(line):
			// Handle skipped tests
			result.Skipped = true
		}
	}

	// Determine overall result
	result.TestCount = testCount
	result.FailureCount = failureCount
	result.ErrorCount = errorCount
	result.Errors = errors
	result.Passed = (failureCount+errorCount) == 0 && testCount > 0
	result.Failed = (failureCount + errorCount) > 0

	// If command failed but no specific test failures found, treat as error
	if cmdErr != nil && !result.Failed && !result.Passed {
		result.Failed = true
		result.ErrorCount = 1
		result.Errors = append(result.Errors, TestError{
			TestName: "compilation",
			Message:  output,
			Type:     "error",
		})
	}
}

// parseGinkgoOutput parses Ginkgo test output
func (g *TestGenerator) parseGinkgoOutput(result *ExecutionResult, output string, cmdErr error) {
	// For now, use similar parsing to go test
	// In a full implementation, you would parse Ginkgo's JSON output
	g.parseGoTestOutput(result, output, cmdErr)

	// Ginkgo-specific patterns could be added here
	if strings.Contains(output, "Ran ") && strings.Contains(output, " of ") {
		// Parse Ginkgo summary line
		// Example: "Ran 5 of 5 Specs in 0.123 seconds"
		summaryPattern := regexp.MustCompile(`Ran (\d+) of (\d+) Specs`)
		if matches := summaryPattern.FindStringSubmatch(output); len(matches) >= 3 {
			if count, err := strconv.Atoi(matches[1]); err == nil {
				result.TestCount = count
			}
		}
	}
}

// GenerateTestFile creates a complete test file for an endpoint
func (g *TestGenerator) GenerateTestFile(endpoint *parser.Endpoint, testCode string) *TestFile {
	fileName := g.generateTestFileName(endpoint)

	return &TestFile{
		Name:        fileName,
		Path:        fileName,
		Content:     testCode,
		Framework:   g.framework,
		Endpoint:    *endpoint,
		GeneratedAt: time.Now(),
		Metadata: map[string]string{
			"generator_version": "1.0.0",
			"go_version":        "1.21",
		},
	}
}

// CreateTestSuite creates a test suite from multiple test files
func (g *TestGenerator) CreateTestSuite(name string, files []TestFile) *TestSuite {
	totalTests := 0
	for i := range files {
		file := &files[i]
		// Count test functions in the file content
		testFuncPattern := regexp.MustCompile(`func\s+Test\w+\s*\(`)
		matches := testFuncPattern.FindAllString(file.Content, -1)
		totalTests += len(matches)
	}

	return &TestSuite{
		Name:        name,
		Files:       files,
		TotalTests:  totalTests,
		Framework:   g.framework,
		GeneratedAt: time.Now(),
	}
}

// ValidateTestCode performs basic validation on generated test code
func (g *TestGenerator) ValidateTestCode(testCode string) error {
	// Check for basic Go syntax requirements
	if !strings.Contains(testCode, "package ") {
		return fmt.Errorf("test code missing package declaration")
	}

	if !strings.Contains(testCode, "func Test") {
		return fmt.Errorf("test code missing test functions")
	}

	// Framework-specific validation
	switch g.framework {
	case "testify":
		if !strings.Contains(testCode, "github.com/stretchr/testify") {
			log.Warn().Msg("testify framework specified but imports not found")
		}
	case "ginkgo":
		if !strings.Contains(testCode, "github.com/onsi/ginkgo") {
			log.Warn().Msg("ginkgo framework specified but imports not found")
		}
	}

	return nil
}

// ExtractTestMetrics extracts metrics from test code
func (g *TestGenerator) ExtractTestMetrics(testCode string) map[string]interface{} {
	metrics := make(map[string]interface{})

	// Count test functions
	testFuncPattern := regexp.MustCompile(`func\s+Test\w+\s*\(`)
	testFunctions := testFuncPattern.FindAllString(testCode, -1)
	metrics["test_function_count"] = len(testFunctions)

	// Count assertions (approximate)
	assertionPatterns := []string{
		"assert\\.",
		"require\\.",
		"Expect\\(",
		"gomega\\.",
	}

	totalAssertions := 0
	for _, pattern := range assertionPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(testCode, -1)
		totalAssertions += len(matches)
	}
	metrics["assertion_count"] = totalAssertions

	// Estimate lines of code
	lines := strings.Split(testCode, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "//") {
			nonEmptyLines++
		}
	}
	metrics["lines_of_code"] = nonEmptyLines

	// Check for different test categories
	categories := []string{
		"happy", "success", "valid",
		"error", "fail", "invalid",
		"boundary", "edge", "limit",
		"security", "auth", "permission",
		"performance", "timeout", "latency",
	}

	foundCategories := make([]string, 0)
	lowerCode := strings.ToLower(testCode)
	for _, category := range categories {
		if strings.Contains(lowerCode, category) {
			foundCategories = append(foundCategories, category)
		}
	}
	metrics["test_categories"] = foundCategories

	return metrics
}
