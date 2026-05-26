// Package logging provides different shortcuts and utilities to simplify work with slog.
package logging

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

var ErrInvalidLevel = errors.New("invalid log level value")

func parseLogLevel(level string) (slog.Level, error) {
	// NB (alkurbatov): Log level names are kept in sync with Python services
	// to make configuration compatible.
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil

	case "info":
		return slog.LevelInfo, nil

	case "warn", "warning":
		return slog.LevelWarn, nil

	case "error", "critical", "fatal":
		return slog.LevelError, nil
	}

	return 0, ErrInvalidLevel
}

// Setup initializes default slog logger.
func Setup(level string, useJSON bool) error {
	logLevel, err := parseLogLevel(level)
	if err != nil {
		return err
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	var handler slog.Handler
	if useJSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(&ContextHandler{handler})
	slog.SetDefault(logger)

	return nil
}
