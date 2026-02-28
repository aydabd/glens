package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopicForEvent_ValidTypes(t *testing.T) {
	tests := []struct {
		eventType string
		wantTopic string
	}{
		{TypeAnalyzeCompleted, "glens-analyze"},
		{TypeTestFailed, "glens-test-results"},
		{TypeReportGenerated, "glens-reports"},
		{TypeSecretStored, "glens-secrets"},
		{TypeExportScheduled, "glens-export"},
	}
	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			topic, err := TopicForEvent(tt.eventType)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTopic, topic)
		})
	}
}

func TestTopicForEvent_UnknownType(t *testing.T) {
	_, err := TopicForEvent("unknown.event")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")
}

func TestTopicMap_CoversAllEventTypes(t *testing.T) {
	allTypes := []string{
		TypeAnalyzeCompleted,
		TypeTestFailed,
		TypeReportGenerated,
		TypeSecretStored,
		TypeExportScheduled,
	}
	for _, et := range allTypes {
		_, ok := TopicMap[et]
		assert.True(t, ok, "TopicMap missing entry for %s", et)
	}
}
