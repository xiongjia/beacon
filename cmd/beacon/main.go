package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/pkg/logger"

	"github.com/xiongjia/beacon/pkg/dbg"
)

func main() {
	logger.SetDefaultSlog(
		logger.LoggerWithLevel("debug"),
		logger.LoggerWithSource(true),
	)
	slog.Debug("test", slog.String("test", "t"))

	dbgMgr := dbg.NewDebugManager(dbg.DebugManagerWithAddr("", 8081))
	err := dbgMgr.Start()
	if err != nil {
		slog.Error("dbg manager start error", slog.Any("error", err))
	}
	select {}
}
