package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"glens/internal/parser"
)

// AnthropicClient implements the Client interface for Anthropic Claude models
type AnthropicClient struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
}

// AnthropicRequest represents the request structure for Anthropic API
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
}

// AnthropicMessage represents a message in Anthropic format
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicResponse represents the response from Anthropic API
type AnthropicResponse struct {
	ID      string             `json:"id"`
	Type    string             `json:"type"`
	Role    string             `json:"role"`
	Content []AnthropicContent `json:"content"`
	Model   string             `json:"model"`
	Usage   AnthropicUsage     `json:"usage"`
}

// AnthropicContent represents content in the response
type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicUsage represents token usage
type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// NewAnthropicClient creates a new Anthropic client
func NewAnthropicClient() (*AnthropicClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "Anthropic"}
	}

	return &AnthropicClient{
		apiKey:    apiKey,
		baseURL:   "https://api.anthropic.com",
		model:     "claude-3-sonnet-20240229",
		maxTokens: 4000,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// GenerateTest generates integration test code using Anthropic Claude
func (c *AnthropicClient) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	startTime := time.Now()

	prompt := c.buildPrompt(endpoint)

	log.Debug().
		Str("model", c.model).
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Msg("Generating test with Anthropic Claude")

	request := AnthropicRequest{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	response, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: err.Error(),
		}
	}

	if len(response.Content) == 0 {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: "no content in response",
		}
	}

	testCode := response.Content[0].Text
	generationTime := time.Since(startTime)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         prompt,
		ModelUsed:      c.model,
		Framework:      "testify",
		TestCategories: []string{"happy-path", "error-handling", "boundary", "security"},
		GeneratedAt:    time.Now().Format(time.RFC3339),
		TokensUsed:     response.Usage.InputTokens + response.Usage.OutputTokens,
		GenerationTime: generationTime.String(),
		Metadata: map[string]string{
			"api_provider":  "anthropic",
			"input_tokens":  fmt.Sprintf("%d", response.Usage.InputTokens),
			"output_tokens": fmt.Sprintf("%d", response.Usage.OutputTokens),
		},
	}

	log.Info().
		Str("model", c.model).
		Dur("generation_time", generationTime).
		Int("tokens_used", result.TokensUsed).
		Msg("Test generation completed with Anthropic Claude")

	return result, nil
}

// GetModelName returns the model name
func (c *AnthropicClient) GetModelName() string {
	return "Anthropic Claude Sonnet"
}

// GetCapabilities returns the capabilities of Anthropic models
func (c *AnthropicClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            c.maxTokens,
		Languages:            []string{"go", "python", "javascript", "java", "rust"},
	}
}

// buildPrompt creates the detailed prompt for test generation
func (c *AnthropicClient) buildPrompt(endpoint *parser.Endpoint) string {
	var prompt bytes.Buffer

	prompt.WriteString("You are an expert software testing engineer specializing in API integration testing with Go.\n\n")
	prompt.WriteString("Generate comprehensive integration tests for the following OpenAPI endpoint using Go and the testify framework:\n\n")

	prompt.WriteString("**Endpoint Details:**\n")
	fmt.Fprintf(&prompt, "- Method: %s\n", endpoint.Method)
	fmt.Fprintf(&prompt, "- Path: %s\n", endpoint.Path)

	if endpoint.OperationID != "" {
		fmt.Fprintf(&prompt, "- Operation ID: %s\n", endpoint.OperationID)
	}

	if endpoint.Summary != "" {
		fmt.Fprintf(&prompt, "- Summary: %s\n", endpoint.Summary)
	}

	if endpoint.Description != "" {
		fmt.Fprintf(&prompt, "- Description: %s\n", endpoint.Description)
	}

	// Parameters
	if len(endpoint.Parameters) > 0 {
		prompt.WriteString("\n**Parameters:**\n")
		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			required := "optional"
			if param.Required {
				required = "required"
			}
			fmt.Fprintf(&prompt, "- %s (%s, %s): %s [Type: %s]\n",
				param.Name, param.In, required, param.Description, param.Schema.Type)
		}
	}

	// Request Body
	if endpoint.RequestBody != nil {
		prompt.WriteString("\n**Request Body:**\n")
		if endpoint.RequestBody.Description != "" {
			fmt.Fprintf(&prompt, "- Description: %s\n", endpoint.RequestBody.Description)
		}
		prompt.WriteString("- Content Types:\n")
		for contentType := range endpoint.RequestBody.Content {
			mediaType := endpoint.RequestBody.Content[contentType]
			fmt.Fprintf(&prompt, "  - %s: %s\n", contentType, mediaType.Schema.Type)
		}
	}

	// Responses
	if len(endpoint.Responses) > 0 {
		prompt.WriteString("\n**Expected Responses:**\n")
		for code, response := range endpoint.Responses {
			fmt.Fprintf(&prompt, "- %s: %s\n", code, response.Description)
		}
	}

	prompt.WriteString("\n**Requirements:**\n")
	prompt.WriteString("1. Use Go programming language with testify framework\n")
	prompt.WriteString("2. Include proper imports and package declaration\n")
	prompt.WriteString("3. Generate realistic test data and scenarios\n")
	prompt.WriteString("4. Cover all response status codes\n")
	prompt.WriteString("5. Test parameter validation (required vs optional)\n")
	prompt.WriteString("6. Include error handling scenarios\n")
	prompt.WriteString("7. Add boundary testing for limits and edge cases\n")
	prompt.WriteString("8. Consider security aspects (auth, validation)\n")
	prompt.WriteString("9. Add performance considerations where applicable\n")
	prompt.WriteString("10. Use descriptive test names and add comments\n")
	prompt.WriteString("11. Include setup and cleanup if necessary\n")
	prompt.WriteString("12. Make tests independent and idempotent\n\n")

	prompt.WriteString("**Test Categories to Include:**\n")
	prompt.WriteString("- Happy path tests with valid inputs\n")
	prompt.WriteString("- Error scenarios with invalid inputs\n")
	prompt.WriteString("- Boundary value testing\n")
	prompt.WriteString("- Security validation tests\n")
	prompt.WriteString("- Schema validation tests\n\n")

	prompt.WriteString("Generate complete, executable Go test code that follows best practices and can be run immediately.")

	return prompt.String()
}

// makeRequest makes an HTTP request to Anthropic API
func (c *AnthropicClient) makeRequest(ctx context.Context, request AnthropicRequest) (*AnthropicResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Msg("failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response AnthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// NewAnthropicClientWithModel creates a new Anthropic client with a specific model
func NewAnthropicClientWithModel(modelName string) (*AnthropicClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "Anthropic"}
	}

	return &AnthropicClient{
		apiKey:    apiKey,
		baseURL:   "https://api.anthropic.com",
		model:     modelName,
		maxTokens: 4000,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}
