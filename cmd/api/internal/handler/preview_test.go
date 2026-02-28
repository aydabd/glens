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

func TestAnalyzePreview_ValidRequest_Returns200(t *testing.T) {
	body := `{"spec_url":"https://example.com/api.json"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze/preview", strings.NewReader(body))
	rec := httptest.NewRecorder()

	AnalyzePreview(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp previewResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/api.json", resp.SpecURL)
	assert.NotEmpty(t, resp.Endpoints, "endpoints should not be empty")
}

func TestAnalyzePreview_ResponseContainsEndpointCategories(t *testing.T) {
	body := `{"spec_url":"https://example.com/api.json"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze/preview", strings.NewReader(body))
	rec := httptest.NewRecorder()

	AnalyzePreview(rec, req)

	var resp previewResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	for _, ep := range resp.Endpoints {
		assert.NotEmpty(t, ep.Path, "endpoint path must not be empty")
		assert.NotEmpty(t, ep.Method, "endpoint method must not be empty")
		assert.NotEmpty(t, ep.RiskLevel, "endpoint risk_level must not be empty")
	}
}

func TestAnalyzePreview_MissingSpecURL_Returns400(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze/preview", strings.NewReader(body))
	rec := httptest.NewRecorder()

	AnalyzePreview(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "spec_url is required")
}

func TestAnalyzePreview_InvalidJSON_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{bad json`},
		{"empty body", ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze/preview", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			AnalyzePreview(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp map[string]string
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Contains(t, resp["error"], "invalid request body")
		})
	}
}

func TestAnalyzePreview_WrongMethod_Returns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/analyze/preview", AnalyzePreview)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/analyze/preview", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}
