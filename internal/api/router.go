package api

import (
	"net/http"

	"github.com/xiongjia/beacon/internal/api/handler"
	"github.com/xiongjia/beacon/pkg/injector"
	"github.com/xiongjia/beacon/pkg/util"
)

type (
	Router struct {
		muxGroup *util.GroupMux
	}
)

func newHandle(handler handler.ApiHandler) *http.ServeMux {
	mux := http.NewServeMux()
	handler.Register(mux)
	return mux
}

func NewRouter(di *injector.Injector) (*Router, error) {
	// create handlers
	dbgHandler, _ := handler.NewDebugHandler(di)

	muxGroup := util.NewMainMuxGroup()
	muxGroup.Group("/api/v1/debug", newHandle(dbgHandler))
	return &Router{muxGroup: muxGroup}, nil
}

func (r *Router) ServerHandler() *http.ServeMux {
	return r.muxGroup.MainMux()
}
