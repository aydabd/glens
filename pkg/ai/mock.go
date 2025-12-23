package ai

import (
	"context"
	"fmt"
	"time"

	"glens/pkg/parser"
)

// MockClient is a mock AI client for testing purposes
type MockClient struct {
	modelName string
}

// NewMockClient creates a new mock AI client
func NewMockClient(modelName string) *MockClient {
	if modelName == "" {
		modelName = "mock"
	}
	return &MockClient{
		modelName: modelName,
	}
}

// GenerateTest generates a mock test for demonstration purposes
func (c *MockClient) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	testCode := c.generateMockTestCode(endpoint)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         c.buildPrompt(endpoint),
		ModelUsed:      c.modelName,
		Framework:      "testify",
		TestCategories: []string{"integration", "api", "mock"},
		GeneratedAt:    time.Now().Format(time.RFC3339),
		GenerationTime: "50ms",
		Metadata: map[string]string{
			"mock": "true",
		},
	}

	return result, nil
}

// GetModelName returns the mock model name
func (c *MockClient) GetModelName() string {
	return c.modelName
}

// GetCapabilities returns mock capabilities
func (c *MockClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            4000,
		Languages:            []string{"go"},
	}
}

// generateMockTestCode creates a realistic-looking Go test
func (c *MockClient) generateMockTestCode(endpoint *parser.Endpoint) string {
	return fmt.Sprintf(`package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test%s%s tests the %s %s endpoint
func Test%s%s(t *testing.T) {
	// Setup
	baseURL := "http://localhost:8080"
	endpoint := "%s"
	
	// Test: Valid request
	t.Run("ValidRequest", func(t *testing.T) {
		req, err := http.NewRequest("%s", baseURL+endpoint, nil)
		require.NoError(t, err)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		// Verify status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK status")
	})
	
	// Test: Invalid endpoint (404)
	t.Run("InvalidEndpoint", func(t *testing.T) {
		req, err := http.NewRequest("%s", baseURL+"/invalid/endpoint", nil)
		require.NoError(t, err)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected 404 Not Found")
	})
}
`,
		capitalize(endpoint.Method), sanitizePath(endpoint.Path),
		endpoint.Method, endpoint.Path,
		capitalize(endpoint.Method), sanitizePath(endpoint.Path),
		endpoint.Path,
		endpoint.Method,
		endpoint.Method,
	)
}

// buildPrompt creates a simple prompt for the mock
func (c *MockClient) buildPrompt(endpoint *parser.Endpoint) string {
	return fmt.Sprintf("Generate test for %s %s", endpoint.Method, endpoint.Path)
}

// Helper functions
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func sanitizePath(path string) string {
	result := ""
	nextUpper := true
	
	for _, r := range path {
		if r == '/' || r == '{' || r == '}' || r == '-' {
			nextUpper = true
			continue
		}
		
		if nextUpper {
			result += string(r - 32)
			nextUpper = false
		} else {
			result += string(r)
		}
	}
	
	if result == "" {
		result = "Root"
	}
	
	return result
}
