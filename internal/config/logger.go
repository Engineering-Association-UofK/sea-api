package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type MultiHandler struct {
	handlers []slog.Handler
	level    slog.Level
}

func NewMultiHandlerLog(level slog.Level) *slog.Logger {
	// Setup JSON file output
	file, _ := os.OpenFile(App.ResourcesDir+"/logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}

	return slog.New(&MultiHandler{
		level: level,
		handlers: []slog.Handler{
			slog.NewJSONHandler(file, opts),
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: "2006-01-02 3:04PM",
			}),
		},
	})
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return l >= m.level
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
