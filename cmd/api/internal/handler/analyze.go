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
	SpecURL           string   `json:"spec_url"`
	Models            []string `json:"models"`
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
		writeProblem(w, r, http.StatusBadRequest, ProblemTypeValidation,
			"Validation Error", fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if req.SpecURL == "" {
		writeProblem(w, r, http.StatusBadRequest, ProblemTypeValidation,
			"Validation Error", "spec_url is required")
		return
	}

	runID, err := generateRunID()
	if err != nil {
		writeProblem(w, r, http.StatusInternalServerError, ProblemTypeInternal,
			"Internal Server Error", fmt.Sprintf("generate run id: %v", err))
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
