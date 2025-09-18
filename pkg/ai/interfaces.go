package ai

import (
	"context"

	"glens/pkg/parser"
)

// Client defines the interface for AI model clients
type Client interface {
	// GenerateTest generates integration test code for an endpoint
	GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (*TestGenerationResult, error)

	// GetModelName returns the name/identifier of the AI model
	GetModelName() string

	// GetCapabilities returns the capabilities of this AI model
	GetCapabilities() ModelCapabilities
}

// TestGenerationResult contains the result of test generation
type TestGenerationResult struct {
	TestCode       string            `json:"test_code"`
	Prompt         string            `json:"prompt"`
	ModelUsed      string            `json:"model_used"`
	Framework      string            `json:"framework"`
	TestCategories []string          `json:"test_categories"`
	GeneratedAt    string            `json:"generated_at"`
	TokensUsed     int               `json:"tokens_used,omitempty"`
	GenerationTime string            `json:"generation_time"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ModelCapabilities describes what the AI model can do
type ModelCapabilities struct {
	SupportsGoTests      bool     `json:"supports_go_tests"`
	SupportsSecurityTest bool     `json:"supports_security_test"`
	SupportedFrameworks  []string `json:"supported_frameworks"`
	MaxTokens            int      `json:"max_tokens"`
	Languages            []string `json:"languages"`
}

// Manager manages multiple AI model clients
type Manager struct {
	clients map[string]Client
}

// NewManager creates a new AI manager with specified models
func NewManager(modelNames []string) (*Manager, error) {
	manager := &Manager{
		clients: make(map[string]Client),
	}

	for _, modelName := range modelNames {
		client, err := createClient(modelName)
		if err != nil {
			return nil, err
		}
		manager.clients[modelName] = client
	}

	return manager, nil
}

// GenerateTest generates a test using the specified AI model
func (m *Manager) GenerateTest(ctx context.Context, modelName string, endpoint *parser.Endpoint) (testCode, modelUsed string, err error) {
	client, exists := m.clients[modelName]
	if !exists {
		return "", "", ErrModelNotFound{Model: modelName}
	}

	result, err := client.GenerateTest(ctx, endpoint)
	if err != nil {
		return "", "", err
	}

	return result.TestCode, result.Prompt, nil
}

// GetAvailableModels returns the names of all available AI models
func (m *Manager) GetAvailableModels() []string {
	var models []string
	for name := range m.clients {
		models = append(models, name)
	}
	return models
}

// GetModelCapabilities returns capabilities for a specific model
func (m *Manager) GetModelCapabilities(modelName string) (ModelCapabilities, error) {
	client, exists := m.clients[modelName]
	if !exists {
		return ModelCapabilities{}, ErrModelNotFound{Model: modelName}
	}

	return client.GetCapabilities(), nil
}

// createClient creates an AI client based on model name
func createClient(modelName string) (Client, error) {
	switch modelName {
	case "gpt4", "openai":
		return NewOpenAIClient()
	case "sonnet4", "anthropic":
		return NewAnthropicClient()
	case "flash-pro", "google":
		return NewGoogleClient()
	case "ollama":
		return NewOllamaClient("")
	case "ollama_codellama":
		return NewOllamaClient("ollama")
	case "ollama_deepseekcoder":
		return NewOllamaClient("ollama_deepseekcoder")
	case "ollama_qwen":
		return NewOllamaClient("ollama_qwen")
	default:
		// Check if it's a custom Ollama model (format: ollama:model-name)
		if len(modelName) > 7 && modelName[:7] == "ollama:" {
			// For custom models, use default ollama config but override model name
			client, err := NewOllamaClient("")
			if err != nil {
				return nil, err
			}
			// Override the model name - need to modify the client struct
			return &OllamaClientWithModel{
				client: client,
				model:  modelName[7:], // Remove "ollama:" prefix
			}, nil
		}
		return nil, ErrUnsupportedModel{Model: modelName}
	}
}
