package events

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Event is the base envelope for all domain events.
type Event struct {
	EventType   string    `json:"event_type"`
	EventID     string    `json:"event_id"`
	Timestamp   time.Time `json:"timestamp"`
	WorkspaceID string    `json:"workspace_id"`
	Payload     any       `json:"payload"`
}

// Event type constants.
const (
	TypeAnalyzeCompleted = "analyze.completed"
	TypeTestFailed       = "test.failed"
	TypeReportGenerated  = "report.generated"
	TypeSecretStored     = "secret.stored"
	TypeExportScheduled  = "export.scheduled"
)

// AnalyzeCompletedPayload is sent when an analysis run finishes.
type AnalyzeCompletedPayload struct {
	RunID  string `json:"run_id"`
	Passed int    `json:"passed"`
	Failed int    `json:"failed"`
}

// TestFailedPayload is sent when an individual test fails.
type TestFailedPayload struct {
	RunID          string `json:"run_id"`
	EndpointPath   string `json:"endpoint_path"`
	EndpointMethod string `json:"endpoint_method"`
	Model          string `json:"model"`
	ErrorMessage   string `json:"error_message"`
}

// ReportGeneratedPayload is sent when a report is ready.
type ReportGeneratedPayload struct {
	RunID     string `json:"run_id"`
	ReportURL string `json:"report_url"`
}

// SecretStoredPayload is sent when a secret is stored.
type SecretStoredPayload struct {
	SecretRef string `json:"secret_ref"`
	Name      string `json:"name"`
}

// ExportScheduledPayload is sent when an export job is queued.
type ExportScheduledPayload struct {
	DatasetID string `json:"dataset_id"`
}

// NewEvent creates an Event with a generated UUID and current timestamp.
func NewEvent(eventType, workspaceID string, payload any) Event {
	return Event{
		EventType:   eventType,
		EventID:     uuid.New().String(),
		Timestamp:   time.Now(),
		WorkspaceID: workspaceID,
		Payload:     payload,
	}
}

// ValidateEvent checks that required fields are present.
func ValidateEvent(e Event) error {
	if e.EventType == "" {
		return errors.New("event_type is required")
	}
	if e.EventID == "" {
		return errors.New("event_id is required")
	}
	if e.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if e.WorkspaceID == "" {
		return errors.New("workspace_id is required")
	}
	if e.Payload == nil {
		return errors.New("payload is required")
	}
	return nil
}
