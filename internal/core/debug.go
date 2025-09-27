package core

import (
	"log/slog"
	"time"

	"github.com/xiongjia/beacon/pkg/injector"
)

type (
	DebugService interface {
		StartDebugServer() error
		StopDebugServer(timeout time.Duration) error
	}

	debugServiceImpl struct{}
)

func NewDebugService(inj *injector.Injector) (DebugService, error) {
	// TODO
	return &debugServiceImpl{}, nil
}

func (*debugServiceImpl) StartDebugServer() error {
	// XXX TODO
	slog.Debug("start debug server")
	return nil
}

func (*debugServiceImpl) StopDebugServer(timeout time.Duration) error {
	// XXX TODO
	slog.Debug("stop debug server")
	return nil
}
