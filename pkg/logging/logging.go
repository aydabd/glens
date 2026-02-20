// Package logging provides a generic logging setup backed by zerolog.
// It can be used in any Go project and never imports internal packages.
package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Level represents a logging level.
type Level string

// Supported log levels.
const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Format controls the output format of log messages.
type Format string

// Supported log formats.
const (
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
)

// Config holds logging configuration.
type Config struct {
	Level  Level
	Format Format
	Output io.Writer // defaults to os.Stderr when nil
}

// Setup configures the global zerolog logger with the given Config.
func Setup(cfg Config) {
	out := cfg.Output
	if out == nil {
		out = os.Stderr
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if cfg.Format == FormatConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	} else {
		log.Logger = zerolog.New(out).With().Timestamp().Logger()
	}

	switch cfg.Level {
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
