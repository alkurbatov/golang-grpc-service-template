package logging_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

func newDefaultLogger(t *testing.T, level string, useJSON bool) *slog.Logger {
	t.Helper()

	err := logging.Setup(level, useJSON)
	require.NoError(t, err)

	return slog.Default()
}

func TestSetupLogLevel(t *testing.T) {
	tt := []struct {
		name     string
		aliases  []string
		expected slog.Level
	}{
		{
			name:     "Set debug level",
			aliases:  []string{"debug", "DEBUG", "Debug"},
			expected: slog.LevelDebug,
		},
		{
			name:     "Enable info level",
			aliases:  []string{"info", "INFO", "Info"},
			expected: slog.LevelInfo,
		},
		{
			name: "Enable warning level",
			aliases: []string{
				"warn", "warning",
				"WARN", "WARNING",
				"Warn", "Warning",
			},
			expected: slog.LevelWarn,
		},
		{
			name: "Enable error level",
			aliases: []string{
				"error", "critical", "fatal",
				"ERROR", "CRITICAL", "FATAL",
				"Error", "Critical", "Fatal",
			},
			expected: slog.LevelError,
		},
	}

	for _, tc := range tt {
		for _, alias := range tc.aliases {
			t.Run(tc.name+" with "+alias, func(t *testing.T) {
				sut := newDefaultLogger(t, alias, true)
				require.True(t, sut.Enabled(t.Context(), tc.expected))
			})
		}
	}
}

func TestSetupOnInvalidLogLevel(t *testing.T) {
	err := logging.Setup("trace", true)
	require.ErrorIs(t, err, logging.ErrInvalidLevel)
}

func TestSetupJSONLogging(t *testing.T) {
	sut := newDefaultLogger(t, "info", true)

	h, ok := sut.Handler().(*logging.ContextHandler)
	require.True(t, ok)

	require.IsType(t, &slog.JSONHandler{}, h.Handler)
}

func TestSetupTextLogging(t *testing.T) {
	sut := newDefaultLogger(t, "info", false)

	h, ok := sut.Handler().(*logging.ContextHandler)
	require.True(t, ok)

	require.IsType(t, &slog.TextHandler{}, h.Handler)
}
