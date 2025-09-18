package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"glens/pkg/parser"
)

// OllamaClient implements Client interface for Ollama local LLM
type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
	config     OllamaConfig
}

// OllamaConfig holds configuration for Ollama client
type OllamaConfig struct {
	BaseURL       string  `mapstructure:"base_url"`
	Model         string  `mapstructure:"model"`
	Timeout       string  `mapstructure:"timeout"`
	Temperature   float64 `mapstructure:"temperature"`
	MaxTokens     int     `mapstructure:"max_tokens"`
	ContextLength int     `mapstructure:"context_length"`
	NumPredict    int     `mapstructure:"num_predict"`
	TopK          int     `mapstructure:"top_k"`
	TopP          float64 `mapstructure:"top_p"`
	RepeatPenalty float64 `mapstructure:"repeat_penalty"`
	Seed          int     `mapstructure:"seed"`
}

// OllamaGenerateRequest represents the request structure for Ollama API
type OllamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Format  string                 `json:"format,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// OllamaGenerateResponse represents the response structure from Ollama API
type OllamaGenerateResponse struct {
	Model          string `json:"model"`
	Response       string `json:"response"`
	Done           bool   `json:"done"`
	Context        []int  `json:"context,omitempty"`
	TotalTime      int64  `json:"total_duration,omitempty"`
	LoadTime       int64  `json:"load_duration,omitempty"`
	PromptEvalTime int64  `json:"prompt_eval_duration,omitempty"`
	EvalTime       int64  `json:"eval_duration,omitempty"`
}

// OllamaModel represents a model in Ollama
type OllamaModel struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
}

// OllamaModelsResponse represents the response from /api/tags endpoint
type OllamaModelsResponse struct {
	Models []OllamaModel `json:"models"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(configKey string) (*OllamaClient, error) {
	var config OllamaConfig

	// Use provided config key or default to "ollama"
	if configKey == "" {
		configKey = "ollama"
	}

	if err := viper.UnmarshalKey(fmt.Sprintf("ai_models.%s", configKey), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Ollama config: %w", err)
	}

	// Set defaults
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}
	if config.Model == "" {
		config.Model = "codellama:7b-instruct"
	}
	if config.Temperature == 0 {
		config.Temperature = 0.1
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 4000
	}

	// Parse timeout
	timeout := 300 * time.Second
	if config.Timeout != "" {
		if parsedTimeout, err := time.ParseDuration(config.Timeout); err == nil {
			timeout = parsedTimeout
		}
	}

	client := &OllamaClient{
		baseURL: config.BaseURL,
		model:   config.Model,
		config:  config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}

	return client, nil
}

// GenerateTest generates integration test code using Ollama
func (c *OllamaClient) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	startTime := time.Now()

	prompt := c.buildPrompt(endpoint)

	log.Info().
		Str("model", c.model).
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Msg("Generating test with Ollama")

	// Create request
	req := OllamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false, // Use non-streaming for simplicity
		Options: map[string]interface{}{
			"temperature":    c.config.Temperature,
			"num_predict":    c.config.NumPredict,
			"top_k":          c.config.TopK,
			"top_p":          c.config.TopP,
			"repeat_penalty": c.config.RepeatPenalty,
		},
	}

	if c.config.Seed >= 0 {
		req.Options["seed"] = c.config.Seed
	}

	// Make API call
	response, err := c.generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate with Ollama: %w", err)
	}

	generationTime := time.Since(startTime)

	// Extract test code from response
	testCode := c.extractTestCode(response.Response)

	result := &TestGenerationResult{
		TestCode:       testCode,
		Prompt:         prompt,
		ModelUsed:      c.model,
		Framework:      "testify", // Default framework
		TestCategories: []string{"integration", "api"},
		GeneratedAt:    time.Now().Format(time.RFC3339),
		GenerationTime: generationTime.String(),
		Metadata: map[string]string{
			"ollama_version":          "latest",
			"total_duration_ms":       fmt.Sprintf("%d", response.TotalTime/1000000),
			"eval_duration_ms":        fmt.Sprintf("%d", response.EvalTime/1000000),
			"prompt_eval_duration_ms": fmt.Sprintf("%d", response.PromptEvalTime/1000000),
		},
	}

	log.Info().
		Str("model", c.model).
		Dur("generation_time", generationTime).
		Int("response_length", len(response.Response)).
		Msg("Test generation completed")

	return result, nil
}

// GetModelName returns the Ollama model name
func (c *OllamaClient) GetModelName() string {
	return fmt.Sprintf("ollama:%s", c.model)
}

// GetCapabilities returns the capabilities of the Ollama model
func (c *OllamaClient) GetCapabilities() ModelCapabilities {
	return ModelCapabilities{
		SupportsGoTests:      true,
		SupportsSecurityTest: true,
		SupportedFrameworks:  []string{"testify", "ginkgo", "standard"},
		MaxTokens:            c.config.MaxTokens,
		Languages:            []string{"go", "json", "yaml"},
	}
}

// HealthCheck verifies if Ollama is running and the model is available
func (c *OllamaClient) HealthCheck(ctx context.Context) error {
	// Check if Ollama is running
	resp, err := c.httpClient.Get(c.baseURL + "/api/version")
	if err != nil {
		return fmt.Errorf("ollama server not accessible: %w", err)
	}
	if closeErr := resp.Body.Close(); closeErr != nil {
		log.Debug().Err(closeErr).Msg("failed to close response body")
	}

	// Check if model is available
	models, err := c.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	for _, model := range models {
		if model.Name == c.model {
			return nil // Model found
		}
	}

	return fmt.Errorf("model %s not found in Ollama. Available models: %v", c.model, c.getModelNames(models))
}

// ListModels returns available models in Ollama
func (c *OllamaClient) ListModels(ctx context.Context) ([]OllamaModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var modelsResp OllamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	return modelsResp.Models, nil
}

// generate makes a generation request to Ollama
func (c *OllamaClient) generate(ctx context.Context, req OllamaGenerateRequest) (*OllamaGenerateResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Ollama: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// buildPrompt creates a prompt optimized for local LLMs to generate Go integration tests
func (c *OllamaClient) buildPrompt(endpoint *parser.Endpoint) string {
	prompt := fmt.Sprintf(`You are a Go developer writing integration tests. Generate a complete Go test function for this OpenAPI endpoint:

Endpoint: %s %s
Summary: %s
Description: %s

Requirements:
1. Use the testify framework
2. Create a test function that covers:
   - Valid request with expected response
   - Invalid request handling
   - Status code validation
   - Response structure validation
3. Include proper imports
4. Use realistic test data
5. Handle authentication if required
6. Test error cases

Generate ONLY the Go test code, no explanations:

`, endpoint.Method, endpoint.Path, endpoint.Summary, endpoint.Description)

	// Add parameters information if available
	if len(endpoint.Parameters) > 0 {
		prompt += "Parameters:\n"
		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			prompt += fmt.Sprintf("- %s (%s): %s\n", param.Name, param.In, param.Description)
		}
		prompt += "\n"
	}

	// Add response information if available
	if len(endpoint.Responses) > 0 {
		prompt += "Expected Responses:\n"
		for code, response := range endpoint.Responses {
			prompt += fmt.Sprintf("- %s: %s\n", code, response.Description)
		}
		prompt += "\n"
	}

	prompt += "```go\n"

	return prompt
}

// extractTestCode extracts Go test code from the Ollama response
func (c *OllamaClient) extractTestCode(response string) string {
	// Ollama responses often include the code block markers
	// Extract code between ```go and ``` markers
	startMarker := "```go"
	endMarker := "```"

	startIdx := bytes.Index([]byte(response), []byte(startMarker))
	if startIdx != -1 {
		startIdx += len(startMarker)
		endIdx := bytes.Index([]byte(response[startIdx:]), []byte(endMarker))
		if endIdx != -1 {
			return response[startIdx : startIdx+endIdx]
		}
	}

	// If no code blocks found, return the response as-is
	// (some models might not use markdown formatting)
	return response
}

// getModelNames extracts model names from OllamaModel slice
func (c *OllamaClient) getModelNames(models []OllamaModel) []string {
	names := make([]string, len(models))
	for i, model := range models {
		names[i] = model.Name
	}
	return names
}

// OllamaClientWithModel wraps OllamaClient to override the model name
type OllamaClientWithModel struct {
	client *OllamaClient
	model  string
}

// GenerateTest delegates to the wrapped client but uses custom model name
func (c *OllamaClientWithModel) GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error) {
	// Temporarily override the model name
	originalModel := c.client.model
	c.client.model = c.model
	defer func() {
		c.client.model = originalModel
	}()

	return c.client.GenerateTest(ctx, endpoint)
}

// GetModelName returns the custom model name
func (c *OllamaClientWithModel) GetModelName() string {
	return fmt.Sprintf("ollama:%s", c.model)
}

// GetCapabilities delegates to the wrapped client
func (c *OllamaClientWithModel) GetCapabilities() ModelCapabilities {
	return c.client.GetCapabilities()
}
