package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// previewRequest is the JSON body for the analyze preview endpoint.
type previewRequest struct {
	SpecURL string `json:"spec_url"`
}

// endpointCategory represents an endpoint with its risk level.
type endpointCategory struct {
	Path      string `json:"path"`
	Method    string `json:"method"`
	RiskLevel string `json:"risk_level"`
}

// previewResponse is returned by the analyze preview endpoint.
type previewResponse struct {
	SpecURL   string             `json:"spec_url"`
	Endpoints []endpointCategory `json:"endpoints"`
}

// AnalyzePreview handles POST /api/v1/analyze/preview requests.
func AnalyzePreview(w http.ResponseWriter, r *http.Request) {
	var req previewRequest
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

	// Stub: return placeholder endpoint categories.
	writeJSON(w, http.StatusOK, previewResponse{
		SpecURL: req.SpecURL,
		Endpoints: []endpointCategory{
			{Path: "/pets", Method: "GET", RiskLevel: "safe"},
			{Path: "/pets", Method: "POST", RiskLevel: "medium"},
			{Path: "/pets/{id}", Method: "DELETE", RiskLevel: "high"},
		},
	})
}
