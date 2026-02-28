package events

import (
	"context"
	"fmt"
)

// Publisher publishes domain events to a messaging backend.
type Publisher interface {
	Publish(ctx context.Context, topic string, event Event) error
}

// TopicMap maps event types to their topic names.
var TopicMap = map[string]string{
	TypeAnalyzeCompleted: "glens-analyze",
	TypeTestFailed:       "glens-test-results",
	TypeReportGenerated:  "glens-reports",
	TypeSecretStored:     "glens-secrets",
	TypeExportScheduled:  "glens-export",
}

// TopicForEvent returns the topic name for the given event type.
func TopicForEvent(eventType string) (string, error) {
	topic, ok := TopicMap[eventType]
	if !ok {
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}
	return topic, nil
}
