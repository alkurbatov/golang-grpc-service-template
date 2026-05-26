package logging

import (
	"context"
	"log/slog"
	"slices"
)

type ctxKey struct{}

var _ slog.Handler = (*ContextHandler)(nil)

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying handler.
func (h *ContextHandler) Handle(
	ctx context.Context,
	record slog.Record, //nolint: gocritic //arg copy is forced by slog.Handler interface
) error {
	if attrs, ok := ctx.Value(ctxKey{}).([]slog.Attr); ok {
		record.AddAttrs(attrs...)
	}

	return h.Handler.Handle(ctx, record)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{h.Handler.WithAttrs(attrs)}
}

// Context adds slog attributes to the provided context so that it will be
// included in any slog.Record created with such context.
func Context(parent context.Context, attrs ...slog.Attr) context.Context {
	var newAttrs []slog.Attr
	if prevAttrs, ok := parent.Value(ctxKey{}).([]slog.Attr); ok {
		newAttrs = prevAttrs
	}

	newAttrs = slices.Concat(newAttrs, attrs)

	return context.WithValue(parent, ctxKey{}, newAttrs)
}
