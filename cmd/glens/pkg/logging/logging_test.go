package logging_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"glens/tools/glens/pkg/logging"
)

func TestSetup_defaults(t *testing.T) {
	prevLogger := log.Logger
	prevLevel := zerolog.GlobalLevel()
	t.Cleanup(func() {
		log.Logger = prevLogger
		zerolog.SetGlobalLevel(prevLevel)
	})

	var buf bytes.Buffer
	logging.Setup(logging.Config{
		Level:  logging.LevelInfo,
		Format: logging.FormatJSON,
		Output: &buf,
	})
	// No panic is the primary assertion; global logger is reconfigured.
}

func TestSetup_console(t *testing.T) {
	prevLogger := log.Logger
	prevLevel := zerolog.GlobalLevel()
	t.Cleanup(func() {
		log.Logger = prevLogger
		zerolog.SetGlobalLevel(prevLevel)
	})

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
