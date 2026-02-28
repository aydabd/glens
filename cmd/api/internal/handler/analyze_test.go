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

func TestAnalyze_ValidRequest_Returns202(t *testing.T) {
	body := `{"spec_url":"https://example.com/api.json","models":["gpt-4o"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Analyze(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp analyzeResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "accepted", resp.Status)
	assert.NotEmpty(t, resp.RunID, "run_id must not be empty")
	assert.Len(t, resp.RunID, 32, "run_id should be 32 hex characters")
}

func TestAnalyze_InvalidJSON_Returns400(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{invalid`},
		{"empty body", ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			Analyze(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))

			var resp ProblemDetail
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, ProblemTypeValidation, resp.Type)
			assert.Equal(t, "Validation Error", resp.Title)
			assert.Equal(t, http.StatusBadRequest, resp.Status)
			assert.Contains(t, resp.Detail, "invalid request body")
			assert.Equal(t, "/api/v1/analyze", resp.Instance)
		})
	}
}

func TestAnalyze_MissingSpecURL_Returns400(t *testing.T) {
	body := `{"models":["gpt-4o"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(body))
	rec := httptest.NewRecorder()

	Analyze(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))

	var resp ProblemDetail
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, ProblemTypeValidation, resp.Type)
	assert.Equal(t, "Validation Error", resp.Title)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Detail, "spec_url is required")
	assert.Equal(t, "/api/v1/analyze", resp.Instance)
}

func TestAnalyze_WrongMethod_Returns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/analyze", Analyze)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/analyze", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}

func TestAnalyze_UniqueRunIDs(t *testing.T) {
	body := `{"spec_url":"https://example.com/api.json"}`
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(body))
		rec := httptest.NewRecorder()

		Analyze(rec, req)

		var resp analyzeResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		assert.False(t, ids[resp.RunID], "run_id should be unique")
		ids[resp.RunID] = true
	}
}
