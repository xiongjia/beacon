package engine

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/xiongjia/beacon/internal/api"
	"github.com/xiongjia/beacon/internal/core"
	"github.com/xiongjia/beacon/pkg/injector"
	"github.com/xiongjia/beacon/pkg/util"
)

type (
	Engine struct {
		mutx   sync.Mutex
		di     *injector.Injector
		router *api.Router
		server *util.HttpServer
	}
)

func NewEngine() (*Engine, error) {
	di := injector.NewInjector()
	if err := injector.Provide(di, core.NewDebugService); err != nil {
		slog.Error("create service error")
		return nil, err
	}
	// TODO create other services

	// create API router
	router, err := api.NewRouter(di)
	if err != nil {
		return nil, err
	}
	return &Engine{di: di, router: router}, nil
}

func (e *Engine) Start() error {
	e.mutx.Lock()
	defer e.mutx.Unlock()
	if e.server != nil {
		return fmt.Errorf("debug server is already running")
	}
	serv := util.NewHttpServer(":8080", e.router.ServerHandler())
	if err := serv.StartServer(); err != nil {
		slog.Error("start server error", slog.Any("error", err))
		return err
	}
	e.server = serv
	return nil
}

func (e *Engine) Stop(timeout time.Duration) error {
	e.mutx.Lock()
	defer e.mutx.Unlock()
	if e.server == nil {
		return fmt.Errorf("server is not running")
	}
	if err := e.server.Shutdown(timeout); err != nil {
		slog.Error("stop server error")
		return err
	}
	e.server = nil
	return nil
}
