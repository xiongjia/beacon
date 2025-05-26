package util

import (
	"log/slog"
	"os"
	"time"
)

type (
	LogOption struct {
		Level     slog.Level
		AddSource bool
	}
)

func NewLog(opts LogOption) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: opts.Level, AddSource: opts.AddSource})
	return slog.New(handler)
}

func InitDefaultLog(opts LogOption) {
	slog.SetDefault(NewLog(opts))
	time.Now()
}
