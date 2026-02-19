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
	case "mock":
		return NewMockClient("mock"), nil
	case "enhanced-mock", "mock-enhanced":
		return NewEnhancedMockClient("enhanced-mock"), nil

	// --- OpenAI ---
	case "gpt4", "openai", "gpt-4-turbo":
		return NewOpenAIClient()
	case "gpt-4o", "gpt4o":
		return NewOpenAIClientWithModel("gpt-4o")
	case "gpt-4o-mini", "gpt4o-mini":
		return NewOpenAIClientWithModel("gpt-4o-mini")
	// Latest OpenAI (2025)
	case "gpt-4.1", "gpt4.1":
		return NewOpenAIClientWithModel("gpt-4.1")
	case "gpt-4.1-mini":
		return NewOpenAIClientWithModel("gpt-4.1-mini")
	case "o3", "openai-o3":
		return NewOpenAIClientWithModel("o3")
	case "o4-mini", "openai-o4-mini":
		return NewOpenAIClientWithModel("o4-mini")

	// --- Anthropic ---
	case "sonnet4", "anthropic", "claude-3-sonnet":
		return NewAnthropicClient()
	case "claude-3.5-sonnet", "claude-3-5-sonnet":
		return NewAnthropicClientWithModel("claude-3-5-sonnet-20241022")
	// Latest Anthropic (2025)
	case "claude-3.7-sonnet", "claude-3-7-sonnet":
		return NewAnthropicClientWithModel("claude-3-7-sonnet-20250219")
	case "claude-opus-4", "claude-4-opus":
		return NewAnthropicClientWithModel("claude-opus-4-5")

	// --- Google ---
	case "flash-pro", "google", "gemini-1.5-flash":
		return NewGoogleClient()
	case "gemini-2.0-flash", "gemini-2-flash":
		return NewGoogleClientWithModel("gemini-2.0-flash")
	case "gemini-2.0-pro", "gemini-2-pro":
		return NewGoogleClientWithModel("gemini-2.0-pro")
	// Latest Google (2025)
	case "gemini-2.5-pro", "gemini-2-5-pro":
		return NewGoogleClientWithModel("gemini-2.5-pro-preview-03-25")
	case "gemini-2.5-flash", "gemini-2-5-flash":
		return NewGoogleClientWithModel("gemini-2.5-flash")

	// --- Ollama (local) ---
	case "ollama":
		return NewOllamaClient("")
	case "ollama_codellama":
		return NewOllamaClient("ollama")
	case "ollama_deepseekcoder", "deepseek-coder":
		return NewOllamaClient("ollama_deepseekcoder")
	case "ollama_qwen", "qwen-coder":
		return NewOllamaClient("ollama_qwen")
	// Latest open-source models via Ollama (2025)
	case "ollama_deepseek-r2", "deepseek-r2":
		return NewOllamaClient("ollama_deepseek-r2")
	case "ollama_qwen3", "qwen3":
		return NewOllamaClient("ollama_qwen3")
	case "ollama_llama4", "llama4":
		return NewOllamaClient("ollama_llama4")

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
