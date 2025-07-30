package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/pkg/util"
)

func main() {
	slog.SetDefault(util.NewLogger(util.LogOption{Level: util.LOG_LEVEL_DEBUG}))
	slog.Debug("test", slog.String("test", "t"))
}
