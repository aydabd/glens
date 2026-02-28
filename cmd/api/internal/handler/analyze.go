package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// analyzeRequest is the JSON body for the analyze endpoint.
type analyzeRequest struct {
	SpecURL          string   `json:"spec_url"`
	Models           []string `json:"models"`
	ApprovedEndpoints []string `json:"approved_endpoints"`
	SkippedEndpoints  []string `json:"skipped_endpoints"`
}

// analyzeResponse is returned when an analysis run is accepted.
type analyzeResponse struct {
	RunID  string `json:"run_id"`
	Status string `json:"status"`
}

// Analyze handles POST /api/v1/analyze requests.
func Analyze(w http.ResponseWriter, r *http.Request) {
	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if req.SpecURL == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("spec_url is required"))
		return
	}

	runID, err := generateRunID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("generate run id: %w", err))
		return
	}

	writeJSON(w, http.StatusAccepted, analyzeResponse{
		RunID:  runID,
		Status: "accepted",
	})
}

func generateRunID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
