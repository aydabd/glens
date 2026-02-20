package ai

import "fmt"

// ErrModelNotFound is returned when a requested AI model is not available
type ErrModelNotFound struct {
	Model string
}

func (e ErrModelNotFound) Error() string {
	return fmt.Sprintf("AI model '%s' not found", e.Model)
}

// ErrUnsupportedModel is returned when a model name is not supported
type ErrUnsupportedModel struct {
	Model string
}

func (e ErrUnsupportedModel) Error() string {
	return fmt.Sprintf("AI model '%s' is not supported", e.Model)
}

// ErrAPIKeyMissing is returned when an API key is missing for a model
type ErrAPIKeyMissing struct {
	Model string
}

func (e ErrAPIKeyMissing) Error() string {
	return fmt.Sprintf("API key missing for AI model '%s'", e.Model)
}

// ErrGenerationFailed is returned when test generation fails
type ErrGenerationFailed struct {
	Model  string
	Reason string
}

func (e ErrGenerationFailed) Error() string {
	return fmt.Sprintf("test generation failed for model '%s': %s", e.Model, e.Reason)
}

// ErrRateLimited is returned when API rate limits are exceeded
type ErrRateLimited struct {
	Model      string
	RetryAfter string
}

func (e ErrRateLimited) Error() string {
	return fmt.Sprintf("rate limited for model '%s', retry after: %s", e.Model, e.RetryAfter)
}
