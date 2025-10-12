package main

import (
	"context"
	"log/slog"

	todoapp "grantjames.github.io/todo-app"
)

type traceIdContextHandler struct{ h slog.Handler }

func traceFrom(ctx context.Context) string {
	if v := ctx.Value(todoapp.TraceIdKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (a *traceIdContextHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return a.h.Enabled(ctx, lvl)
}

func (a *traceIdContextHandler) Handle(ctx context.Context, r slog.Record) error {
	rr := r.Clone()
	rr.AddAttrs(slog.String("trace_id", traceFrom(ctx)))
	return a.h.Handle(ctx, rr)
}

func (a *traceIdContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceIdContextHandler{h: a.h.WithAttrs(attrs)}
}

func (a *traceIdContextHandler) WithGroup(name string) slog.Handler {
	return &traceIdContextHandler{h: a.h.WithGroup(name)}
}
