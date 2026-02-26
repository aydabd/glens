package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- newOllamaLocal / local model shortcuts ---

func TestNewOllamaLocal_CreatesClient(t *testing.T) {
	tests := []struct {
		modelName string
	}{
		{"mistral"},
		{"mistral-nemo"},
		{"llama3"},
		{"phi4"},
		{"gemma2"},
	}

	for _, tt := range tests {
		t.Run(tt.modelName, func(t *testing.T) {
			t.Parallel()
			c, err := newOllamaLocal(tt.modelName)
			require.NoError(t, err)
			assert.NotNil(t, c)
			assert.Equal(t, "ollama:"+tt.modelName, c.GetModelName())
		})
	}
}

func TestCreateClient_LocalModelShortcuts(t *testing.T) {
	tests := []struct {
		shortcut      string
		wantModelName string
	}{
		// Mistral local
		{"mistral-local", "ollama:mistral"},
		{"mistral7b", "ollama:mistral"},
		{"mistral-nemo-local", "ollama:mistral-nemo"},
		{"mistral-small-local", "ollama:mistral-small"},
		// Meta Llama local
		{"llama3-local", "ollama:llama3"},
		{"llama3", "ollama:llama3"},
		{"llama3.1-local", "ollama:llama3.1"},
		{"llama3.1", "ollama:llama3.1"},
		{"llama3.2-local", "ollama:llama3.2"},
		{"llama3.2", "ollama:llama3.2"},
		// Microsoft Phi local
		{"phi3-local", "ollama:phi3"},
		{"phi3", "ollama:phi3"},
		{"phi4-local", "ollama:phi4"},
		{"phi4", "ollama:phi4"},
		// Google Gemma local
		{"gemma2-local", "ollama:gemma2"},
		{"gemma2", "ollama:gemma2"},
		{"gemma3-local", "ollama:gemma3"},
		{"gemma3", "ollama:gemma3"},
	}

	for _, tt := range tests {
		t.Run(tt.shortcut, func(t *testing.T) {
			t.Parallel()
			c, err := createClient(tt.shortcut)
			require.NoError(t, err, "shortcut %q should not require an API key", tt.shortcut)
			assert.NotNil(t, c)
			assert.Equal(t, tt.wantModelName, c.GetModelName(),
				"shortcut %q: unexpected model name", tt.shortcut)
		})
	}
}

func TestCreateClient_LocalModelShortcuts_Capabilities(t *testing.T) {
	c, err := createClient("mistral-local")
	require.NoError(t, err)

	caps := c.GetCapabilities()
	assert.True(t, caps.SupportsGoTests)
	assert.Contains(t, caps.SupportedFrameworks, "testify")
}

// --- OllamaClient.PullModel ---

func TestOllamaClient_PullModel_Success(t *testing.T) {
	// Simulate a streaming Ollama pull response
	lines := []OllamaPullResponse{
		{Status: "pulling manifest"},
		{Status: "pulling layer", Total: 1000, Completed: 500},
		{Status: "pulling layer", Total: 1000, Completed: 1000},
		{Status: "pull complete"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/pull", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		for _, line := range lines {
			_ = enc.Encode(line)
		}
	}))
	defer srv.Close()

	client := newTestOllamaClient(t, srv.URL)

	var out strings.Builder
	err := client.PullModel(context.Background(), "mistral", &out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "pulling manifest")
}

func TestOllamaClient_PullModel_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"model not found"}`))
	}))
	defer srv.Close()

	client := newTestOllamaClient(t, srv.URL)
	err := client.PullModel(context.Background(), "nonexistent", io.Discard)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestOllamaClient_PullModel_ConnectionRefused(t *testing.T) {
	client := newTestOllamaClient(t, "http://127.0.0.1:1") // nothing listening
	err := client.PullModel(context.Background(), "mistral", io.Discard)
	assert.Error(t, err)
}

// newTestOllamaClient builds an OllamaClient pointed at the given base URL.
func newTestOllamaClient(t *testing.T, baseURL string) *OllamaClient {
	t.Helper()
	c, err := NewOllamaClient("")
	require.NoError(t, err)
	c.baseURL = baseURL
	return c
}
