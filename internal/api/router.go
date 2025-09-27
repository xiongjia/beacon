package api

import (
	"github.com/xiongjia/beacon/pkg/injector"
)

type (
	Router struct {
	}
)

func NewRouter(di *injector.Injector) (*Router, error) {
	// mux := http.NewServeMux()
	// mux.Handle("/api/v1/debug", )

	return &Router{}, nil
}
