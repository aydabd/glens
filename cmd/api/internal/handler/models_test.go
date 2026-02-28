package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModels_GET_ReturnsModelList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/models", nil)
	rec := httptest.NewRecorder()

	Models(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp struct {
		Models []model `json:"models"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Models, "models list should not be empty")
}

func TestModels_ContainsExpectedModels(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/models", nil)
	rec := httptest.NewRecorder()

	Models(rec, req)

	var resp struct {
		Models []model `json:"models"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	ids := make([]string, len(resp.Models))
	for i, m := range resp.Models {
		ids[i] = m.ID
	}

	expectedIDs := []string{"gpt-4o", "gpt-4o-mini", "claude-sonnet-4-20250514", "claude-3-5-haiku-20241022"}
	for _, id := range expectedIDs {
		assert.Contains(t, ids, id)
	}
}

func TestModels_ModelFieldsPopulated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/models", nil)
	rec := httptest.NewRecorder()

	Models(rec, req)

	var resp struct {
		Models []model `json:"models"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	for _, m := range resp.Models {
		assert.NotEmpty(t, m.ID, "model ID must not be empty")
		assert.NotEmpty(t, m.Name, "model Name must not be empty")
		assert.NotEmpty(t, m.Provider, "model Provider must not be empty")
	}
}

func TestModels_WrongMethod_Returns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/models", Models)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/models", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		})
	}
}
