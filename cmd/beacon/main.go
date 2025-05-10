package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/pkg/util"
)

func main() {
	util.InitDefaultLog(util.LogOption{Level: slog.LevelDebug, AddSource: true})
	slog.Debug("test")
}
