package handler

import (
	"encoding/json"
	"net/http"
)

// healthResponse is the JSON body returned by the health endpoint.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Health returns a handler that reports service health.
func Health(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, healthResponse{
			Status:  "ok",
			Version: version,
		})
	}
}

// writeJSON marshals v to JSON and writes it to w with the given status code.
// It encodes to a buffer first so that encoding failures are caught before
// headers are sent, avoiding a mixed/corrupted response body.
func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
	_, _ = w.Write([]byte("\n"))
}
