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

	"glens/pkg/parser"
)

// GoogleClient implements the Client interface for Google Gemini models
type GoogleClient struct {
	apiKey    string
	baseURL   string
	model     string
	maxTokens int
	client    *http.Client
	projectID string
}

// GoogleRequest represents the request structure for Google Gemini API
type GoogleRequest struct {
	Contents         []GoogleContent        `json:"contents"`
	GenerationConfig GoogleGenerationConfig `json:"generationConfig"`
}

// GoogleContent represents content in Google format
type GoogleContent struct {
	Parts []GooglePart `json:"parts"`
}

// GooglePart represents a part of content
type GooglePart struct {
	Text string `json:"text"`
}

// GoogleGenerationConfig represents generation configuration
type GoogleGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"topP"`
	TopK            int     `json:"topK"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

// GoogleResponse represents the response from Google Gemini API
type GoogleResponse struct {
	Candidates    []GoogleCandidate   `json:"candidates"`
	UsageMetadata GoogleUsageMetadata `json:"usageMetadata"`
}

// GoogleCandidate represents a response candidate
type GoogleCandidate struct {
	Content      GoogleContent `json:"content"`
	FinishReason string        `json:"finishReason"`
	Index        int           `json:"index"`
}

// GoogleUsageMetadata represents usage metadata
type GoogleUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// NewGoogleClient creates a new Google Gemini client
func NewGoogleClient() (*GoogleClient, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyMissing{Model: "Google"}
	}

	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		projectID = "default-project" // Use a default if not specified
	}

	return &GoogleClient{
		apiKey:    apiKey,
		baseURL:   "https://generativelanguage.googleapis.com/v1beta",
		model:     "gemini-1.5-flash",
		maxTokens: 4000,
		projectID: projectID,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// GenerateTest generates integration test code using Google Gemini
func (c *GoogleClient) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	startTime := time.Now()

	prompt := c.buildPrompt(endpoint)

	log.Debug().
		Str("model", c.model).
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Msg("Generating test with Google Gemini")

	request := GoogleRequest{
		Contents: []GoogleContent{
			{
				Parts: []GooglePart{
					{
						Text: prompt,
					},
				},
			},
		},
		GenerationConfig: GoogleGenerationConfig{
			Temperature:     0.7,
			TopP:            0.8,
			TopK:            40,
			MaxOutputTokens: c.maxTokens,
		},
	}

	response, err := c.makeRequest(ctx, request)
	if err != nil {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: err.Error(),
		}
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, ErrGenerationFailed{
			Model:  c.GetModelName(),
			Reason: "no content in response candidates",
		}
	}

	testCode := response.Candidates[0].Content.Parts[0].Text
	generationTime := time.Since(startTime)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         prompt,
		ModelUsed:      c.model,
		Framework:      "testify",
		TestCategories: []string{"happy-path", "error-handling", "boundary", "security"},
		GeneratedAt:    time.Now().Format(time.RFC3339),
		TokensUsed:     response.UsageMetadata.TotalTokenCount,
		GenerationTime: generationTime.String(),
		Metadata: map[string]string{
			"api_provider":          "google",
			"finish_reason":         response.Candidates[0].FinishReason,
			"prompt_token_count":    fmt.Sprintf("%d", response.UsageMetadata.PromptTokenCount),
			"candidate_token_count": fmt.Sprintf("%d", response.UsageMetadata.CandidatesTokenCount),
		},
	}

	log.Info().
		Str("model", c.model).
		Dur("generation_time", generationTime).
		Int("tokens_used", response.UsageMetadata.TotalTokenCount).
		Msg("Test generation completed with Google Gemini")

	return result, nil
}

// GetModelName returns the model name
func (c *GoogleClient) GetModelName() string {
	return "Google Gemini Flash Pro"
}

// GetCapabilities returns the capabilities of Google models
func (c *GoogleClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            c.maxTokens,
		Languages:            []string{"go", "python", "javascript", "java", "cpp", "rust"},
	}
}

// buildPrompt creates the detailed prompt for test generation
func (c *GoogleClient) buildPrompt(endpoint *parser.Endpoint) string {
	var prompt bytes.Buffer

	prompt.WriteString("As an expert software testing engineer, generate comprehensive integration tests for this OpenAPI endpoint using Go and testify.\n\n")

	prompt.WriteString("**ENDPOINT SPECIFICATION:**\n")
	prompt.WriteString(fmt.Sprintf("Method: %s\n", endpoint.Method))
	prompt.WriteString(fmt.Sprintf("Path: %s\n", endpoint.Path))

	if endpoint.OperationID != "" {
		prompt.WriteString(fmt.Sprintf("Operation ID: %s\n", endpoint.OperationID))
	}

	if endpoint.Summary != "" {
		prompt.WriteString(fmt.Sprintf("Summary: %s\n", endpoint.Summary))
	}

	if endpoint.Description != "" {
		prompt.WriteString(fmt.Sprintf("Description: %s\n", endpoint.Description))
	}

	// Parameters
	if len(endpoint.Parameters) > 0 {
		prompt.WriteString("\n**PARAMETERS:**\n")
		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			required := "Optional"
			if param.Required {
				required = "Required"
			}
			prompt.WriteString(fmt.Sprintf("• %s (%s, %s): %s [Type: %s]\n",
				param.Name, param.In, required, param.Description, param.Schema.Type))
		}
	}

	// Request Body
	if endpoint.RequestBody != nil {
		prompt.WriteString("\n**REQUEST BODY:**\n")
		if endpoint.RequestBody.Description != "" {
			prompt.WriteString(fmt.Sprintf("Description: %s\n", endpoint.RequestBody.Description))
		}
		prompt.WriteString("Supported Content Types:\n")
		for contentType := range endpoint.RequestBody.Content {
			mediaType := endpoint.RequestBody.Content[contentType]
			prompt.WriteString(fmt.Sprintf("• %s: %s\n", contentType, mediaType.Schema.Type))
		}
	}

	// Responses
	if len(endpoint.Responses) > 0 {
		prompt.WriteString("\n**EXPECTED RESPONSES:**\n")
		for code, response := range endpoint.Responses {
			prompt.WriteString(fmt.Sprintf("• HTTP %s: %s\n", code, response.Description))
		}
	}

	prompt.WriteString("\n**REQUIREMENTS:**\n")
	prompt.WriteString("Generate Go integration tests that include:\n\n")
	prompt.WriteString("1. **Package and Imports**: Proper Go package declaration with necessary imports\n")
	prompt.WriteString("2. **Test Structure**: Use testify/assert and testify/suite if needed\n")
	prompt.WriteString("3. **Happy Path Tests**: Valid requests with expected successful responses\n")
	prompt.WriteString("4. **Error Handling**: Invalid inputs, missing required fields, wrong types\n")
	prompt.WriteString("5. **Boundary Testing**: Edge cases, limits, empty values, max values\n")
	prompt.WriteString("6. **Security Tests**: Authentication validation, authorization checks\n")
	prompt.WriteString("7. **Schema Validation**: Response structure and data type validation\n")
	prompt.WriteString("8. **HTTP Method Specific**: Appropriate tests for the HTTP method\n")
	prompt.WriteString("9. **Parameter Testing**: All parameter types (path, query, header)\n")
	prompt.WriteString("10. **Performance Checks**: Response time assertions where relevant\n\n")

	prompt.WriteString("**CODE STANDARDS:**\n")
	prompt.WriteString("• Use descriptive test names (TestEndpoint_Scenario_ExpectedResult)\n")
	prompt.WriteString("• Include setup and teardown functions if needed\n")
	prompt.WriteString("• Add comments explaining complex test scenarios\n")
	prompt.WriteString("• Use table-driven tests for multiple scenarios\n")
	prompt.WriteString("• Generate realistic test data\n")
	prompt.WriteString("• Include proper error checking and assertions\n")
	prompt.WriteString("• Make tests independent and idempotent\n\n")

	prompt.WriteString("Generate complete, executable Go test code that can be run immediately without modifications.")

	return prompt.String()
}

// makeRequest makes an HTTP request to Google Gemini API
func (c *GoogleClient) makeRequest(ctx context.Context, request GoogleRequest) (*GoogleResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, c.model, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	var response GoogleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
