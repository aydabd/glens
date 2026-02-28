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

func TestProblemDetail_ContentType(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		method  string
		path    string
		body    string
	}{
		{
			name:    "analyze invalid JSON",
			handler: Analyze,
			method:  http.MethodPost,
			path:    "/api/v1/analyze",
			body:    `{bad`,
		},
		{
			name:    "analyze missing spec_url",
			handler: Analyze,
			method:  http.MethodPost,
			path:    "/api/v1/analyze",
			body:    `{"models":["gpt-4o"]}`,
		},
		{
			name:    "preview invalid JSON",
			handler: AnalyzePreview,
			method:  http.MethodPost,
			path:    "/api/v1/analyze/preview",
			body:    `{bad`,
		},
		{
			name:    "preview missing spec_url",
			handler: AnalyzePreview,
			method:  http.MethodPost,
			path:    "/api/v1/analyze/preview",
			body:    `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"),
				"error responses must use application/problem+json")
		})
	}
}

func TestProblemDetail_RequiredFields(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	Analyze(rec, req)

	var p ProblemDetail
	err := json.NewDecoder(rec.Body).Decode(&p)
	require.NoError(t, err)

	assert.NotEmpty(t, p.Type, "type must be set")
	assert.NotEmpty(t, p.Title, "title must be set")
	assert.NotZero(t, p.Status, "status must be set")
	assert.NotEmpty(t, p.Detail, "detail must be set")
	assert.NotEmpty(t, p.Instance, "instance must be set")
}

func TestProblemDetail_InstanceMatchesRequestPath(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		path    string
		body    string
	}{
		{
			name:    "analyze path",
			handler: Analyze,
			path:    "/api/v1/analyze",
			body:    `{}`,
		},
		{
			name:    "preview path",
			handler: AnalyzePreview,
			path:    "/api/v1/analyze/preview",
			body:    `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			var p ProblemDetail
			err := json.NewDecoder(rec.Body).Decode(&p)
			require.NoError(t, err)
			assert.Equal(t, tt.path, p.Instance,
				"instance must match the request path")
		})
	}
}

func TestProblemDetail_StatusMatchesHTTPCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(`{bad`))
	rec := httptest.NewRecorder()

	Analyze(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var p ProblemDetail
	err := json.NewDecoder(rec.Body).Decode(&p)
	require.NoError(t, err)
	assert.Equal(t, rec.Code, p.Status,
		"ProblemDetail.status must match the HTTP status code")
}

func TestProblemDetail_TypeIsValidURI(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	Analyze(rec, req)

	var p ProblemDetail
	err := json.NewDecoder(rec.Body).Decode(&p)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(p.Type, "https://"),
		"type URI must start with https://")
}
