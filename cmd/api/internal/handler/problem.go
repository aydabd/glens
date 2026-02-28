package handler

import (
	"encoding/json"
	"net/http"
)

// ProblemDetail represents an RFC 9457 Problem Details response.
type ProblemDetail struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail"`
	Instance string            `json:"instance"`
	Errors   []ValidationError `json:"errors,omitempty"`
}

// ValidationError describes a single field validation failure.
type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Problem type URI constants.
const (
	ProblemTypeValidation = "https://glens.dev/errors/validation"
	ProblemTypeInternal   = "https://glens.dev/errors/internal"
)

// writeProblem writes an RFC 9457 Problem Details JSON response.
func writeProblem(w http.ResponseWriter, r *http.Request, status int, problemType, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	p := ProblemDetail{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "failed to encode problem response", http.StatusInternalServerError)
	}
}
