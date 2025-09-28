package handler

import (
	"fmt"
	"net/http"

	"github.com/xiongjia/beacon/internal/core"
	"github.com/xiongjia/beacon/pkg/injector"
)

type (
	DebugHandler struct {
		dbgService core.DebugService
	}
)

func NewDebugHandler(di *injector.Injector) (ApiHandler, error) {
	dbgService, err := injector.Invoke[core.DebugService](di)
	if err != nil {
		return nil, err
	}
	return &DebugHandler{dbgService: dbgService}, nil
}

func (*DebugHandler) Register(mux *http.ServeMux) {
	mux.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "test")
	}))
}
