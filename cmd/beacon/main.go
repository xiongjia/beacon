package main

import (
	"log/slog"

	"github.com/xiongjia/beacon/pkg/util"
)

func main() {
	util.InitDefaultLog(util.LogOption{Level: "debug"})
	slog.Debug("test", slog.Int("n1", 1))
}
