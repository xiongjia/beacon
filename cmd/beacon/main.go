package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/internal/engine"
	"github.com/xiongjia/beacon/pkg/logger"
)

func main() {
	// TODO Load config

	// init log setting
	logger.SetDefaultSlog(
		logger.LoggerWithLevel("debug"),
		logger.LoggerWithSource(true),
	)
	slog.Debug("test", slog.String("test", "t"))

	// TODO loading engine config
	eng, err := engine.NewEngine()
	if err != nil {
		slog.Error("create engine error", slog.Any("error", err))
		return
	}

	// TODO engine stop and other config
	if err := eng.Start(); err != nil {
		slog.Error("engine error", slog.Any("error", err))
	}
}
