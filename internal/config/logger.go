package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandlerLog() *slog.Logger {
	// Setup JSON file output
	file, _ := os.Create("resources/logs.json")
	jsonHandler := slog.NewJSONHandler(file, nil)

	// Setup Terminal (Tint) output
	terminalHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "2006-01-02 3:04PM",
	})

	// Combine loggers
	return slog.New(&MultiHandler{
		handlers: []slog.Handler{jsonHandler, terminalHandler},
	})
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return true
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, r.Clone())
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	// Implementation for groups...
	return m
}
