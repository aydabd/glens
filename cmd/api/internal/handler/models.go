package handler

import "net/http"

// model represents a supported AI model.
type model struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

var supportedModels = []model{
	{ID: "gpt-4o", Name: "GPT-4o", Provider: "openai"},
	{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Provider: "openai"},
	{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", Provider: "anthropic"},
	{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", Provider: "anthropic"},
}

// Models handles GET /api/v1/models requests.
func Models(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"models": supportedModels,
	})
}
