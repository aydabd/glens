package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"glens/internal/parser"
)

// testEndpoint creates a simple endpoint for tests.
func testEndpoint(method, path string) *parser.Endpoint {
	return &parser.Endpoint{
		Method: method,
		Path:   path,
	}
}

// --- MockClient ---

func TestMockClient_GetModelName(t *testing.T) {
	c := NewMockClient("mock")
	assert.Equal(t, "mock", c.GetModelName())
}

func TestMockClient_DefaultModelName(t *testing.T) {
	c := NewMockClient("")
	assert.Equal(t, "mock", c.GetModelName())
}

func TestMockClient_GetCapabilities(t *testing.T) {
	c := NewMockClient("mock")
	caps := c.GetCapabilities()
	assert.True(t, caps.SupportsGoTests)
	assert.Contains(t, caps.SupportedFrameworks, "testify")
}

func TestMockClient_GenerateTest(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		wantFunc string
	}{
		{"GET users", "GET", "/users", "TestGETUsers"},
		{"POST posts", "POST", "/posts", "TestPOSTPosts"},
		{"GET user by id", "GET", "/users/{id}", "TestGETUsersId"},
	}

	c := NewMockClient("mock")
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := testEndpoint(tt.method, tt.path)
			result, err := c.GenerateTest(ctx, ep)
			require.NoError(t, err)

			assert.NotEmpty(t, result.TestCode)
			assert.Equal(t, "mock", result.ModelUsed)
			assert.Equal(t, "testify", result.Framework)
			assert.Contains(t, result.TestCode, tt.wantFunc)
		})
	}
}

// --- EnhancedMockClient ---

func TestEnhancedMockClient_GetModelName(t *testing.T) {
	c := NewEnhancedMockClient("enhanced-mock")
	assert.Equal(t, "enhanced-mock", c.GetModelName())
}

func TestEnhancedMockClient_DefaultModelName(t *testing.T) {
	c := NewEnhancedMockClient("")
	assert.Equal(t, "enhanced-mock", c.GetModelName())
}

func TestEnhancedMockClient_GetCapabilities(t *testing.T) {
	c := NewEnhancedMockClient("enhanced-mock")
	caps := c.GetCapabilities()
	assert.True(t, caps.SupportsGoTests)
	assert.True(t, caps.SupportsSecurityTest)
	assert.Contains(t, caps.SupportedFrameworks, "testify")
}

func TestEnhancedMockClient_GenerateTest_Scenarios(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		path          string
		wantTestFunc  string
		wantScenarios []string
	}{
		{
			name:         "GET users",
			method:       "GET",
			path:         "/users",
			wantTestFunc: "TestGETUsers",
			wantScenarios: []string{
				"Success",
				"EdgeCases",
				"Errors",
				"Security",
				"Performance",
			},
		},
		{
			name:         "POST posts",
			method:       "POST",
			path:         "/posts",
			wantTestFunc: "TestPOSTPosts",
			wantScenarios: []string{
				"Success",
				"Security",
				"Performance",
			},
		},
		{
			name:         "DELETE user",
			method:       "DELETE",
			path:         "/users/{id}",
			wantTestFunc: "TestDELETEUsersId",
			wantScenarios: []string{
				"Success",
				"Errors",
			},
		},
	}

	c := NewEnhancedMockClient("enhanced-mock")
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := testEndpoint(tt.method, tt.path)
			result, err := c.GenerateTest(ctx, ep)
			require.NoError(t, err)

			assert.NotEmpty(t, result.TestCode)
			assert.Equal(t, "enhanced-mock", result.ModelUsed)
			assert.Equal(t, "testify", result.Framework)
			assert.Contains(t, result.TestCode, tt.wantTestFunc)

			for _, scenario := range tt.wantScenarios {
				assert.Contains(t, result.TestCode, scenario,
					"expected scenario %q in test code", scenario)
			}

			// Metadata should include quality scores
			assert.NotEmpty(t, result.Metadata["completeness"])
			assert.NotEmpty(t, result.Metadata["overall_quality"])
		})
	}
}

func TestEnhancedMockClient_GenerateTest_ValidGoSyntax(t *testing.T) {
	c := NewEnhancedMockClient("enhanced-mock")
	ctx := context.Background()

	ep := testEndpoint("GET", "/items")
	result, err := c.GenerateTest(ctx, ep)
	require.NoError(t, err)

	// Basic syntax checks
	assert.Contains(t, result.TestCode, "package main")
	assert.Contains(t, result.TestCode, "import (")
	assert.Contains(t, result.TestCode, "testing")
	assert.Contains(t, result.TestCode, "testify/assert")
	assert.Contains(t, result.TestCode, "testify/require")
	assert.True(t, strings.Contains(result.TestCode, "func Test"),
		"should contain a test function")
}

func TestEnhancedMockClient_Categories(t *testing.T) {
	c := NewEnhancedMockClient("enhanced-mock")
	ctx := context.Background()

	ep := testEndpoint("POST", "/users")
	result, err := c.GenerateTest(ctx, ep)
	require.NoError(t, err)

	assert.Contains(t, result.TestCategories, "integration")
	assert.Contains(t, result.TestCategories, "security")
}

// --- Manager ---

func TestManager_MockModel(t *testing.T) {
	m, err := NewManager([]string{"mock"})
	require.NoError(t, err)

	models := m.GetAvailableModels()
	assert.Contains(t, models, "mock")
}

func TestManager_EnhancedMockModel(t *testing.T) {
	m, err := NewManager([]string{"enhanced-mock"})
	require.NoError(t, err)

	ctx := context.Background()
	ep := testEndpoint("GET", "/ping")

	code, _, err := m.GenerateTest(ctx, "enhanced-mock", ep)
	require.NoError(t, err)
	assert.NotEmpty(t, code)
}

func TestManager_UnknownModel(t *testing.T) {
	_, err := NewManager([]string{"unknown-model-xyz"})
	assert.Error(t, err)
}

func TestManager_ModelNotFound(t *testing.T) {
	m, err := NewManager([]string{"mock"})
	require.NoError(t, err)

	ctx := context.Background()
	ep := testEndpoint("GET", "/ping")

	_, _, err = m.GenerateTest(ctx, "nonexistent", ep)
	assert.Error(t, err)
}

// TestCreateClient_RequiresAPIKey verifies that cloud models return an error
// when the required environment variable is not set.
// Env vars are process-global, so subtests must not run in parallel.
func TestCreateClient_RequiresAPIKey(t *testing.T) {
	// Clear all API keys for the duration of the test so the result is
	// deterministic regardless of the developer's local environment.
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("MISTRAL_API_KEY", "")

	cloudModels := []string{
		// OpenAI
		"gpt-4o", "gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano",
		"o3", "o3-mini", "o4-mini", "codex",
		// Anthropic
		"claude-3.5-sonnet", "claude-3.7-sonnet",
		"claude-sonnet-4", "claude-opus-4", "claude-haiku-4",
		// Google
		"gemini-2.5-pro", "gemini-2.5-flash",
		"gemini-2.0-flash", "gemini-2.0-pro",
		// Mistral
		"mistral", "mistral-medium", "mistral-small",
		"codestral", "mistral-nemo",
	}

	for _, name := range cloudModels {
		name := name
		t.Run(name, func(t *testing.T) {
			_, err := createClient(name)
			// Without an API key set the client must return an error.
			assert.Error(t, err, "model %q should require an API key", name)
		})
	}
}

// TestCreateClient_LocalModels verifies that local Ollama models don't require
// an API key to construct.
func TestCreateClient_LocalModels(t *testing.T) {
	localModels := []string{
		"mock", "enhanced-mock",
	}
	for _, name := range localModels {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c, err := createClient(name)
			assert.NoError(t, err, "model %q should not need an API key", name)
			assert.NotNil(t, c)
		})
	}
}

// --- Helper functions ---

func TestCapitalize(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"get", "Get"},
		{"GET", "GET"},
		{"post", "Post"},
		{"", ""},
		{"123", "123"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, capitalize(tt.in), "capitalize(%q)", tt.in)
	}
}

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"/users", "Users"},
		{"/users/{id}", "UsersId"},
		{"/", "Root"},
		{"", "Root"},
		{"/some-path", "SomePath"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, sanitizePath(tt.in), "sanitizePath(%q)", tt.in)
	}
}
