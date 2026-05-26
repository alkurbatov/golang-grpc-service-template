package logging

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

var _ slog.Handler = (*HandlerMock)(nil)

type HandlerMock struct {
	attrs   []slog.Attr
	Records []slog.Record
}

func (m *HandlerMock) Enabled(ctx context.Context, level slog.Level) bool {
	// NB (alkurbatov): We always want the handler to be enabled during tests.
	return true
}

func (m *HandlerMock) Handle(
	ctx context.Context,
	record slog.Record, //nolint: gocritic //arg copy is forced by slog.Handler interface
) error {
	record.AddAttrs(m.attrs...)
	m.Records = append(m.Records, record)

	return nil
}

func (m *HandlerMock) WithAttrs(attrs []slog.Attr) slog.Handler {
	m.attrs = slices.Concat(m.attrs, attrs)
	return m
}

func (m *HandlerMock) WithGroup(name string) slog.Handler {
	return m
}

// DumpRecordedMsg returns content of particular recorded message.
func (m *HandlerMock) DumpRecordedMsg(index int) string {
	msg := []string{m.Records[index].Message}

	m.Records[index].Attrs(func(attr slog.Attr) bool {
		if attr.Key == "request_id" {
			tokens := strings.Split(attr.Value.String(), "-")
			msg = append(msg, attr.Key+"="+tokens[0])

			return true
		}

		if attr.Key == "duration" {
			msg = append(msg, attr.Key+"=xxx")
			return true
		}

		msg = append(msg, attr.String())

		return true
	})

	return strings.Join(msg, " ")
}
