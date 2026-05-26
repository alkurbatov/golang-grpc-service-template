package logging_test

import (
	"log/slog"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

func TestLogWithoutAttrs(t *testing.T) {
	m := &logging.HandlerMock{}
	l := logging.NewLoggerMock(m)

	l.InfoContext(t.Context(), "Test message")

	require.Len(t, m.Records, 1)
	snaps.MatchSnapshot(t, m.DumpRecordedMsg(0))
}

func TestLogWithContextAttrs(t *testing.T) {
	m := &logging.HandlerMock{}
	l := logging.NewLoggerMock(m).With(
		slog.String("child", "logger"),
	)

	ctx := logging.Context(t.Context(),
		slog.String("ctx_key1", "value1"),
		slog.String("ctx_key2", "value2"),
	)
	ctx = logging.Context(ctx,
		slog.String("ctx_key3", "value3"),
	)

	l.InfoContext(ctx, "Test message",
		slog.String("duration", "150ms"),
	)

	require.Len(t, m.Records, 1)
	snaps.MatchSnapshot(t, m.DumpRecordedMsg(0))
}
