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

// OpenAIClient implements the Client interface for OpenAI GPT models
type OpenAIClient struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "OpenAI"}
	}

	return &OpenAIClient{
		apiKey:    apiKey,
		baseURL:   "https://api.openai.com/v1",
		model:     "gpt-4-turbo",
		maxTokens: 4000,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// GenerateTest generates integration test code using OpenAI GPT
func (c *OpenAIClient) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	startTime := time.Now()

	prompt := c.buildPrompt(endpoint)

	log.Debug().
		Str("model", c.model).
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Msg("Generating test with OpenAI")

	request := OpenAIRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: c.getSystemPrompt(),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   c.maxTokens,
		Temperature: 0.7,
	}

	response, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: err.Error(),
		}
	}

	if len(response.Choices) == 0 {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: "no response choices returned",
		}
	}

	testCode := response.Choices[0].Message.Content
	generationTime := time.Since(startTime)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         prompt,
		ModelUsed:      c.model,
		Framework:      "testify",
		TestCategories: []string{"happy-path", "error-handling", "boundary", "security"},
		GeneratedAt:    time.Now().Format(time.RFC3339),
		TokensUsed:     response.Usage.TotalTokens,
		GenerationTime: generationTime.String(),
		Metadata: map[string]string{
			"api_provider":      "openai",
			"finish_reason":     response.Choices[0].FinishReason,
			"prompt_tokens":     fmt.Sprintf("%d", response.Usage.PromptTokens),
			"completion_tokens": fmt.Sprintf("%d", response.Usage.CompletionTokens),
		},
	}

	log.Info().
		Str("model", c.model).
		Dur("generation_time", generationTime).
		Int("tokens_used", response.Usage.TotalTokens).
		Msg("Test generation completed with OpenAI")

	return result, nil
}

// GetModelName returns the model name
func (c *OpenAIClient) GetModelName() string {
	return "OpenAI GPT-4"
}

// GetCapabilities returns the capabilities of OpenAI models
func (c *OpenAIClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            c.maxTokens,
		Languages:            []string{"go", "python", "javascript", "java"},
	}
}

// getSystemPrompt returns the system prompt for test generation
func (c *OpenAIClient) getSystemPrompt() string {
	return `You are an expert software testing engineer specializing in API integration testing with Go.

Your task is to generate comprehensive integration tests for OpenAPI endpoints using the Go programming language and the testify framework.

Generate tests that cover:
1. Happy path scenarios with valid inputs
2. Error handling with invalid inputs
3. Boundary testing with edge cases
4. Security testing for authentication/authorization
5. Performance considerations

Requirements:
- Use Go with testify framework
- Include proper imports and package declaration
- Add comprehensive assertions
- Generate realistic test data
- Include setup and teardown if needed
- Add clear test names and descriptions
- Handle different HTTP methods appropriately
- Validate response schemas and status codes
- Test both positive and negative scenarios

Provide clean, production-ready Go test code that can be executed immediately.`
}

// buildPrompt creates the detailed prompt for test generation
func (c *OpenAIClient) buildPrompt(endpoint *parser.Endpoint) string {
	var prompt bytes.Buffer

	prompt.WriteString("Generate comprehensive integration tests for this OpenAPI endpoint:\n\n")
	fmt.Fprintf(&prompt, "**Method:** %s\n", endpoint.Method)
	fmt.Fprintf(&prompt, "**Path:** %s\n", endpoint.Path)

	if endpoint.OperationID != "" {
		fmt.Fprintf(&prompt, "**Operation ID:** %s\n", endpoint.OperationID)
	}

	if endpoint.Summary != "" {
		fmt.Fprintf(&prompt, "**Summary:** %s\n", endpoint.Summary)
	}

	if endpoint.Description != "" {
		fmt.Fprintf(&prompt, "**Description:** %s\n\n", endpoint.Description)
	}

	// Parameters
	if len(endpoint.Parameters) > 0 {
		prompt.WriteString("**Parameters:**\n")
		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			required := "optional"
			if param.Required {
				required = "required"
			}
			fmt.Fprintf(&prompt, "- %s (%s, %s): %s - Type: %s\n",
				param.Name, param.In, required, param.Description, param.Schema.Type)
		}
		prompt.WriteString("\n")
	}

	// Request Body
	if endpoint.RequestBody != nil {
		prompt.WriteString("**Request Body:**\n")
		if endpoint.RequestBody.Description != "" {
			fmt.Fprintf(&prompt, "Description: %s\n", endpoint.RequestBody.Description)
		}
		prompt.WriteString("Content Types:\n")
		for contentType := range endpoint.RequestBody.Content {
			mediaType := endpoint.RequestBody.Content[contentType]
			fmt.Fprintf(&prompt, "- %s: %s\n", contentType, mediaType.Schema.Type)
		}
		prompt.WriteString("\n")
	}

	// Responses
	if len(endpoint.Responses) > 0 {
		prompt.WriteString("**Expected Responses:**\n")
		for code, response := range endpoint.Responses {
			fmt.Fprintf(&prompt, "- %s: %s\n", code, response.Description)
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("Generate Go integration tests using testify that:\n")
	prompt.WriteString("1. Test all documented response codes\n")
	prompt.WriteString("2. Validate request/response schemas\n")
	prompt.WriteString("3. Include error scenarios\n")
	prompt.WriteString("4. Test parameter validation\n")
	prompt.WriteString("5. Include performance assertions\n")
	prompt.WriteString("6. Add security considerations\n")
	prompt.WriteString("\nProvide complete, executable Go test code.")

	return prompt.String()
}

// makeRequest makes an HTTP request to OpenAI API
func (c *OpenAIClient) makeRequest(ctx context.Context, request OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// NewOpenAIClientWithModel creates a new OpenAI client with a specific model
func NewOpenAIClientWithModel(modelName string) (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "OpenAI"}
	}

	return &OpenAIClient{
		apiKey:    apiKey,
		baseURL:   "https://api.openai.com/v1",
		model:     modelName,
		maxTokens: 4000,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// NewMistralClient creates a client for Mistral AI (OpenAI-compatible API).
// Requires MISTRAL_API_KEY env var.
func NewMistralClient(modelName string) (*OpenAIClient, error) {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "Mistral"}
	}

	return &OpenAIClient{
		apiKey:    apiKey,
		baseURL:   "https://api.mistral.ai/v1",
		model:     modelName,
		maxTokens: 4000,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}
