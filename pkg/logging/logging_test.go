package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"glens/pkg/logging"
)

func TestSetup_defaults(_ *testing.T) {
	var buf bytes.Buffer
	logging.Setup(logging.Config{
		Level:  logging.LevelInfo,
		Format: logging.FormatJSON,
		Output: &buf,
	})
	// No panic is the primary assertion; global logger is reconfigured.
}

func TestSetup_console(_ *testing.T) {
	var buf bytes.Buffer
	logging.Setup(logging.Config{
		Level:  logging.LevelDebug,
		Format: logging.FormatConsole,
		Output: &buf,
	})
}

func TestLevelConstants(t *testing.T) {
	levels := []logging.Level{
		logging.LevelDebug,
		logging.LevelInfo,
		logging.LevelWarn,
		logging.LevelError,
	}
	for _, l := range levels {
		if strings.TrimSpace(string(l)) == "" {
			t.Errorf("unexpected empty level")
		}
	}
}
