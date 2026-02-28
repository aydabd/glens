package events

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvent_CreatesValidEvent(t *testing.T) {
	before := time.Now()
	e := NewEvent(TypeAnalyzeCompleted, "ws-1", AnalyzeCompletedPayload{
		RunID:  "run-1",
		Passed: 5,
		Failed: 2,
	})
	after := time.Now()

	assert.Equal(t, TypeAnalyzeCompleted, e.EventType)
	assert.Equal(t, "ws-1", e.WorkspaceID)

	_, err := uuid.Parse(e.EventID)
	require.NoError(t, err, "EventID must be a valid UUID")

	assert.False(t, e.Timestamp.Before(before), "timestamp should not be before call")
	assert.False(t, e.Timestamp.After(after), "timestamp should not be after call")
	assert.NoError(t, ValidateEvent(e))
}

func TestValidateEvent_MissingFields(t *testing.T) {
	tests := []struct {
		name   string
		event  Event
		errMsg string
	}{
		{
			name:   "missing event_type",
			event:  Event{EventID: "id", Timestamp: time.Now(), WorkspaceID: "ws"},
			errMsg: "event_type is required",
		},
		{
			name:   "missing event_id",
			event:  Event{EventType: "t", Timestamp: time.Now(), WorkspaceID: "ws"},
			errMsg: "event_id is required",
		},
		{
			name:   "missing timestamp",
			event:  Event{EventType: "t", EventID: "id", WorkspaceID: "ws"},
			errMsg: "timestamp is required",
		},
		{
			name:   "missing workspace_id",
			event:  Event{EventType: "t", EventID: "id", Timestamp: time.Now()},
			errMsg: "workspace_id is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvent(tt.event)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestValidateEvent_AllFieldsPresent(t *testing.T) {
	e := Event{
		EventType:   TypeTestFailed,
		EventID:     "abc-123",
		Timestamp:   time.Now(),
		WorkspaceID: "ws-1",
	}
	assert.NoError(t, ValidateEvent(e))
}

func TestPayloadTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		payload   any
	}{
		{
			name:      "AnalyzeCompletedPayload",
			eventType: TypeAnalyzeCompleted,
			payload:   AnalyzeCompletedPayload{RunID: "r1", Passed: 10, Failed: 1},
		},
		{
			name:      "TestFailedPayload",
			eventType: TypeTestFailed,
			payload: TestFailedPayload{
				RunID:          "r1",
				EndpointPath:   "/users",
				EndpointMethod: "GET",
				Model:          "gpt-4",
				ErrorMessage:   "status 500",
			},
		},
		{
			name:      "ReportGeneratedPayload",
			eventType: TypeReportGenerated,
			payload:   ReportGeneratedPayload{RunID: "r1", ReportURL: "https://example.com/report"},
		},
		{
			name:      "SecretStoredPayload",
			eventType: TypeSecretStored,
			payload:   SecretStoredPayload{SecretRef: "ref-1", Name: "api-key"},
		},
		{
			name:      "ExportScheduledPayload",
			eventType: TypeExportScheduled,
			payload:   ExportScheduledPayload{DatasetID: "ds-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvent(tt.eventType, "ws-test", tt.payload)
			assert.Equal(t, tt.eventType, e.EventType)
			assert.Equal(t, tt.payload, e.Payload)
			assert.NoError(t, ValidateEvent(e))
		})
	}
}
