package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"glens/internal/parser"
)

// EnhancedMockClient is an improved mock AI client with modern features
type EnhancedMockClient struct {
	modelName       string
	patterns        map[string]TestPattern
	enableSecurity  bool
	enableEdgeCases bool
	enableStreaming bool
}

// TestPattern defines a test generation pattern
type TestPattern struct {
	Name        string
	Description string
	Template    func(*parser.Endpoint) string
	Scenarios   []string
}

// NewEnhancedMockClient creates a new enhanced mock AI client
func NewEnhancedMockClient(modelName string) *EnhancedMockClient {
	if modelName == "" {
		modelName = "enhanced-mock"
	}

	client := &EnhancedMockClient{
		modelName:       modelName,
		enableSecurity:  true,
		enableEdgeCases: true,
		enableStreaming: false,
		patterns:        make(map[string]TestPattern),
	}

	// Initialize patterns
	client.initializePatterns()

	return client
}

// GenerateTest generates an enhanced mock test
func (c *EnhancedMockClient) GenerateTest(_ context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	startTime := time.Now()

	// Select appropriate pattern based on endpoint
	pattern := c.selectPattern(endpoint)

	// Generate test code using pattern
	testCode := c.generateEnhancedTestCode(endpoint, pattern)

	// Calculate quality metrics
	metrics := c.calculateQualityMetrics(testCode, endpoint)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         c.buildPrompt(endpoint),
		ModelUsed:      c.modelName,
		Framework:      "testify",
		TestCategories: c.identifyCategories(endpoint),
		GeneratedAt:    time.Now().Format(time.RFC3339),
		GenerationTime: time.Since(startTime).String(),
		Metadata: map[string]string{
			"mock":             "true",
			"enhanced":         "true",
			"pattern":          pattern.Name,
			"security_enabled": fmt.Sprintf("%t", c.enableSecurity),
			"edge_cases":       fmt.Sprintf("%t", c.enableEdgeCases),
			"completeness":     fmt.Sprintf("%.1f", metrics.Completeness),
			"security_score":   fmt.Sprintf("%.1f", metrics.SecurityCoverage),
			"overall_quality":  fmt.Sprintf("%.1f", metrics.OverallScore),
		},
	}

	return result, nil
}

// GetModelName returns the enhanced mock model name
func (c *EnhancedMockClient) GetModelName() string {
	return c.modelName
}

// GetCapabilities returns enhanced capabilities
func (c *EnhancedMockClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            8000,
		Languages:            []string{"go"},
	}
}

// initializePatterns sets up test generation patterns
func (c *EnhancedMockClient) initializePatterns() {
	// CRUD Pattern
	c.patterns["crud_read"] = TestPattern{
		Name:        "CRUD Read",
		Description: "Tests for GET endpoints",
		Scenarios:   []string{"success", "not_found", "invalid_params", "auth"},
	}

	c.patterns["crud_create"] = TestPattern{
		Name:        "CRUD Create",
		Description: "Tests for POST endpoints",
		Scenarios:   []string{"success", "invalid_data", "duplicate", "auth"},
	}

	c.patterns["crud_update"] = TestPattern{
		Name:        "CRUD Update",
		Description: "Tests for PUT/PATCH endpoints",
		Scenarios:   []string{"success", "not_found", "invalid_data", "auth"},
	}

	c.patterns["crud_delete"] = TestPattern{
		Name:        "CRUD Delete",
		Description: "Tests for DELETE endpoints",
		Scenarios:   []string{"success", "not_found", "cascade", "auth"},
	}
}

// selectPattern chooses the appropriate test pattern
func (c *EnhancedMockClient) selectPattern(endpoint *parser.Endpoint) TestPattern {
	method := strings.ToUpper(endpoint.Method)

	switch method {
	case "GET":
		return c.patterns["crud_read"]
	case "POST":
		return c.patterns["crud_create"]
	case "PUT", "PATCH":
		return c.patterns["crud_update"]
	case "DELETE":
		return c.patterns["crud_delete"]
	default:
		return c.patterns["crud_read"]
	}
}

// identifyCategories identifies test categories for the endpoint
func (c *EnhancedMockClient) identifyCategories(endpoint *parser.Endpoint) []string {
	categories := []string{"integration", "api"}

	// Add method-specific categories
	method := strings.ToUpper(endpoint.Method)
	switch method {
	case "GET":
		categories = append(categories, "read", "query")
	case "POST":
		categories = append(categories, "create", "mutation")
	case "PUT", "PATCH":
		categories = append(categories, "update", "mutation")
	case "DELETE":
		categories = append(categories, "delete", "mutation")
	}

	// Add security category if enabled
	if c.enableSecurity {
		categories = append(categories, "security", "auth")
	}

	// Add edge cases category if enabled
	if c.enableEdgeCases {
		categories = append(categories, "edge-cases", "boundary")
	}

	return categories
}

// generateEnhancedTestCode creates comprehensive test code
func (c *EnhancedMockClient) generateEnhancedTestCode(endpoint *parser.Endpoint, pattern TestPattern) string {
	testName := fmt.Sprintf("Test%s%s", capitalize(endpoint.Method), sanitizePath(endpoint.Path))

	var testCases strings.Builder

	// Add header
	testCases.WriteString("package main\n\n")
	testCases.WriteString("import (\n")
	testCases.WriteString("\t\"net/http\"\n")
	testCases.WriteString("\t\"testing\"\n")
	testCases.WriteString("\t\"time\"\n\n")
	testCases.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	testCases.WriteString("\t\"github.com/stretchr/testify/require\"\n")
	testCases.WriteString(")\n\n")

	// Add main test function
	fmt.Fprintf(&testCases, "// %s tests the %s %s endpoint\n", testName, endpoint.Method, endpoint.Path)
	fmt.Fprintf(&testCases, "// Pattern: %s\n", pattern.Name)
	fmt.Fprintf(&testCases, "func %s(t *testing.T) {\n", testName)
	testCases.WriteString("\tbaseURL := \"http://localhost:8080\"\n")
	fmt.Fprintf(&testCases, "\tendpoint := \"%s\"\n\n", endpoint.Path)

	// Add test scenarios
	c.addSuccessTest(&testCases, endpoint)

	if c.enableEdgeCases {
		c.addEdgeCaseTests(&testCases, endpoint)
	}

	c.addErrorTests(&testCases, endpoint)

	if c.enableSecurity {
		c.addSecurityTests(&testCases, endpoint)
	}

	c.addPerformanceTest(&testCases, endpoint)

	testCases.WriteString("}\n")

	return testCases.String()
}

// addSuccessTest adds the happy path test
func (c *EnhancedMockClient) addSuccessTest(sb *strings.Builder, endpoint *parser.Endpoint) {
	sb.WriteString("\t// Test: Success scenario\n")
	sb.WriteString("\tt.Run(\"Success\", func(t *testing.T) {\n")
	fmt.Fprintf(sb, "\t\treq, err := http.NewRequest(\"%s\", baseURL+endpoint, nil)\n", strings.ToUpper(endpoint.Method))
	sb.WriteString("\t\trequire.NoError(t, err)\n\n")
	sb.WriteString("\t\tclient := &http.Client{Timeout: 10 * time.Second}\n")
	sb.WriteString("\t\tresp, err := client.Do(req)\n")
	sb.WriteString("\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\tdefer resp.Body.Close()\n\n")
	sb.WriteString("\t\t// Verify status code\n")

	expectedStatus := "http.StatusOK"
	if strings.ToUpper(endpoint.Method) == "POST" {
		expectedStatus = "http.StatusCreated"
	}
	fmt.Fprintf(sb, "\t\tassert.Equal(t, %s, resp.StatusCode)\n", expectedStatus)
	sb.WriteString("\t})\n\n")
}

// addEdgeCaseTests adds boundary and edge case tests
func (c *EnhancedMockClient) addEdgeCaseTests(sb *strings.Builder, endpoint *parser.Endpoint) {
	sb.WriteString("\t// Test: Edge cases\n")
	sb.WriteString("\tt.Run(\"EdgeCases\", func(t *testing.T) {\n")
	sb.WriteString("\t\tt.Run(\"EmptyResponse\", func(t *testing.T) {\n")
	fmt.Fprintf(sb, "\t\t\treq, err := http.NewRequest(\"%s\", baseURL+endpoint, nil)\n", strings.ToUpper(endpoint.Method))
	sb.WriteString("\t\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\t\treq.Header.Set(\"Accept\", \"application/json\")\n\n")
	sb.WriteString("\t\t\tclient := &http.Client{}\n")
	sb.WriteString("\t\t\tresp, err := client.Do(req)\n")
	sb.WriteString("\t\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\t\tdefer resp.Body.Close()\n")
	sb.WriteString("\t\t})\n")
	sb.WriteString("\t})\n\n")
}

// addErrorTests adds error handling tests
func (c *EnhancedMockClient) addErrorTests(sb *strings.Builder, endpoint *parser.Endpoint) {
	sb.WriteString("\t// Test: Error scenarios\n")
	sb.WriteString("\tt.Run(\"Errors\", func(t *testing.T) {\n")
	sb.WriteString("\t\tt.Run(\"NotFound\", func(t *testing.T) {\n")
	fmt.Fprintf(sb, "\t\t\treq, err := http.NewRequest(\"%s\", baseURL+\"/invalid/endpoint\", nil)\n", strings.ToUpper(endpoint.Method))
	sb.WriteString("\t\t\trequire.NoError(t, err)\n\n")
	sb.WriteString("\t\t\tclient := &http.Client{}\n")
	sb.WriteString("\t\t\tresp, err := client.Do(req)\n")
	sb.WriteString("\t\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\t\tdefer resp.Body.Close()\n\n")
	sb.WriteString("\t\t\tassert.Equal(t, http.StatusNotFound, resp.StatusCode)\n")
	sb.WriteString("\t\t})\n")
	sb.WriteString("\t})\n\n")
}

// addSecurityTests adds security-related tests
func (c *EnhancedMockClient) addSecurityTests(sb *strings.Builder, endpoint *parser.Endpoint) {
	sb.WriteString("\t// Test: Security scenarios\n")
	sb.WriteString("\tt.Run(\"Security\", func(t *testing.T) {\n")
	sb.WriteString("\t\tt.Run(\"Unauthorized\", func(t *testing.T) {\n")
	fmt.Fprintf(sb, "\t\t\treq, err := http.NewRequest(\"%s\", baseURL+endpoint, nil)\n", strings.ToUpper(endpoint.Method))
	sb.WriteString("\t\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\t\t// Don't set Authorization header\n\n")
	sb.WriteString("\t\t\tclient := &http.Client{}\n")
	sb.WriteString("\t\t\tresp, err := client.Do(req)\n")
	sb.WriteString("\t\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\t\tdefer resp.Body.Close()\n\n")
	sb.WriteString("\t\t\t// Should return 401 or 403\n")
	sb.WriteString("\t\t\tassert.Contains(t, []int{http.StatusUnauthorized, http.StatusForbidden}, resp.StatusCode)\n")
	sb.WriteString("\t\t})\n")
	sb.WriteString("\t})\n\n")
}

// addPerformanceTest adds performance validation test
func (c *EnhancedMockClient) addPerformanceTest(sb *strings.Builder, endpoint *parser.Endpoint) {
	sb.WriteString("\t// Test: Performance\n")
	sb.WriteString("\tt.Run(\"Performance\", func(t *testing.T) {\n")
	sb.WriteString("\t\tstart := time.Now()\n")
	fmt.Fprintf(sb, "\t\treq, err := http.NewRequest(\"%s\", baseURL+endpoint, nil)\n", strings.ToUpper(endpoint.Method))
	sb.WriteString("\t\trequire.NoError(t, err)\n\n")
	sb.WriteString("\t\tclient := &http.Client{}\n")
	sb.WriteString("\t\tresp, err := client.Do(req)\n")
	sb.WriteString("\t\trequire.NoError(t, err)\n")
	sb.WriteString("\t\tdefer resp.Body.Close()\n\n")
	sb.WriteString("\t\tduration := time.Since(start)\n")
	sb.WriteString("\t\t// Response should be under 2 seconds\n")
	sb.WriteString("\t\tassert.Less(t, duration, 2*time.Second, \"Response time should be under 2s\")\n")
	sb.WriteString("\t})\n\n")
}

// buildPrompt creates a comprehensive prompt
func (c *EnhancedMockClient) buildPrompt(endpoint *parser.Endpoint) string {
	return fmt.Sprintf("Generate comprehensive integration test for %s %s with security and edge cases",
		endpoint.Method, endpoint.Path)
}

// calculateQualityMetrics estimates test quality
func (c *EnhancedMockClient) calculateQualityMetrics(testCode string, _ *parser.Endpoint) TestQualityMetrics {
	metrics := TestQualityMetrics{}

	// Completeness: based on test scenarios covered
	scenarios := 0
	if strings.Contains(testCode, "Success") {
		scenarios++
	}
	if strings.Contains(testCode, "Errors") {
		scenarios++
	}
	if strings.Contains(testCode, "Security") {
		scenarios++
	}
	if strings.Contains(testCode, "EdgeCases") {
		scenarios++
	}
	if strings.Contains(testCode, "Performance") {
		scenarios++
	}
	metrics.Completeness = float64(scenarios) * 20.0 // 5 scenarios = 100%

	// Security coverage
	if strings.Contains(testCode, "Unauthorized") {
		metrics.SecurityCoverage = 70.0
	}
	if strings.Contains(testCode, "Authorization") {
		metrics.SecurityCoverage = 85.0
	}

	// Edge case coverage
	if strings.Contains(testCode, "EdgeCases") {
		metrics.EdgeCaseCoverage = 75.0
	}

	// Maintainability: based on structure
	if strings.Contains(testCode, "t.Run") {
		metrics.Maintainability = 80.0
	}

	// Overall score
	metrics.OverallScore = (metrics.Completeness + metrics.SecurityCoverage +
		metrics.EdgeCaseCoverage + metrics.Maintainability) / 4.0

	return metrics
}

// TestQualityMetrics represents test quality measurements
type TestQualityMetrics struct {
	Completeness     float64
	SecurityCoverage float64
	EdgeCaseCoverage float64
	Maintainability  float64
	OverallScore     float64
}
