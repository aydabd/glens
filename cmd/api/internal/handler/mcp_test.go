package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCP_ToolsList_ReturnsTools(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/mcp", strings.NewReader(body))
	rec := httptest.NewRecorder()

	MCP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp jsonRPCResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, float64(1), resp.ID)
	assert.Nil(t, resp.Error)

	tools, ok := resp.Result.([]any)
	require.True(t, ok, "result should be a list of tools")
	assert.Len(t, tools, 2)

	// Verify tool names
	names := make([]string, len(tools))
	for i, tool := range tools {
		m := tool.(map[string]any)
		names[i] = m["name"].(string)
	}
	assert.Contains(t, names, "analyze")
	assert.Contains(t, names, "models")
}

func TestMCP_ToolsCall_ReturnsStubResult(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"analyze"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/mcp", strings.NewReader(body))
	rec := httptest.NewRecorder()

	MCP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp jsonRPCResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, float64(2), resp.ID)
	assert.Nil(t, resp.Error)

	result, ok := resp.Result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "stub", result["status"])
}

func TestMCP_UnknownMethod_ReturnsError(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":3,"method":"unknown/method"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/mcp", strings.NewReader(body))
	rec := httptest.NewRecorder()

	MCP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp jsonRPCResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, float64(3), resp.ID)
	require.NotNil(t, resp.Error)
	assert.Equal(t, -32601, resp.Error.Code)
	assert.Equal(t, "method not found", resp.Error.Message)
}

func TestMCP_InvalidJSON_ReturnsParseError(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{not valid}`},
		{"empty body", ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/mcp", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			MCP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp jsonRPCResponse
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, "2.0", resp.JSONRPC)
			require.NotNil(t, resp.Error)
			assert.Equal(t, -32700, resp.Error.Code)
			assert.Contains(t, resp.Error.Message, "parse error")
		})
	}
}

func TestMCP_WrongMethod_Returns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/mcp", MCP)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/mcp", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}
