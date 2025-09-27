package engine

import (
	"log/slog"

	"github.com/xiongjia/beacon/internal/api"
	"github.com/xiongjia/beacon/internal/core"
	"github.com/xiongjia/beacon/pkg/injector"
)

type (
	Engine struct {
		di     *injector.Injector
		router *api.Router
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

func (*Engine) Start() error {
	return nil
}
