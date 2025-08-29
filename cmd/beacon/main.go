package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/pkg/logger"
)

func main() {
	logger.SetDefaultSlog(
		logger.LoggerWithLevel("debug"),
		logger.LoggerWithSource(true),
	)
	slog.Debug("test", slog.String("test", "t"))
}
