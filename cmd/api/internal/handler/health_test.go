package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth_GET_ReturnsStatusOK(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"returns status ok with version", "1.0.0"},
		{"empty version string", ""},
		{"dev version", "dev-abc123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			rec := httptest.NewRecorder()

			Health(tt.version)(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var resp healthResponse
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, "ok", resp.Status)
			assert.Equal(t, tt.version, resp.Version)
		})
	}
}

func TestHealth_WrongMethod_Returns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", Health("1.0.0"))

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/healthz", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}
