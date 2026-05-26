package logging

import "log/slog"

func NewLoggerMock(m *HandlerMock) *slog.Logger {
	return slog.New(&ContextHandler{m})
}
