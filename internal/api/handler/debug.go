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

func NewDebugHandler(di *injector.Injector) ApiHandler {
	return &DebugHandler{}
}

func (*DebugHandler) Register(mux *http.ServeMux) {
	mux.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "test")
	}))
}
