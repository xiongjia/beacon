package dbg

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"path"
	"strconv"
	"sync"
	"time"
)

type (
	DebugOption func(*dbgOption)

	dbgOption struct {
		baseUrl string
		addr    string
	}

	DebugManager struct {
		addr string
		mux  *http.ServeMux
		mutx sync.Mutex

		server *http.Server
	}
)

const (
	DEFAULT_BASEURL = "/debug"
)

var (
	DEFAULT_ADDR = net.JoinHostPort("127.0.0.1", "6060")
)

func NewDebugManager(opts ...DebugOption) *DebugManager {
	dbgOpts := dbgOption{
		baseUrl: DEFAULT_BASEURL,
		addr:    DEFAULT_ADDR,
	}
	for _, opt := range opts {
		opt(&dbgOpts)
	}

	mux := http.NewServeMux()
	// Debugger handler for pprof
	mux.HandleFunc("GET "+path.Join(dbgOpts.baseUrl, "/pprof")+"/", pprof.Index)
	mux.HandleFunc("GET "+path.Join(dbgOpts.baseUrl, "/pprof/cmdline"), pprof.Cmdline)
	mux.HandleFunc("GET "+path.Join(dbgOpts.baseUrl, "/pprof/profile"), pprof.Profile)
	mux.HandleFunc("GET "+path.Join(dbgOpts.baseUrl, "/pprof/symbol"), pprof.Symbol)
	mux.HandleFunc("GET "+path.Join(dbgOpts.baseUrl, "/pprof/trace"), pprof.Trace)
	// Debugger Handler for expvar
	mux.Handle("GET "+path.Join(dbgOpts.baseUrl, "/vars"), expvar.Handler())
	return &DebugManager{mux: mux, addr: dbgOpts.addr}
}

func DebugManagerWithBaseUrl(baseUrl string) DebugOption {
	return func(opt *dbgOption) {
		if baseUrl != "" {
			opt.baseUrl = baseUrl
		}
	}
}

func DebugManagerWithAddr(host string, port int16) DebugOption {
	return func(opt *dbgOption) {
		opt.addr = net.JoinHostPort(host, strconv.Itoa(int(port)))
	}
}

func (d *DebugManager) Start() error {
	d.mutx.Lock()
	defer d.mutx.Unlock()

	if d.server != nil {
		slog.Error("debug server is already running", slog.String("addr", d.addr))
		return fmt.Errorf("debug server already running")
	}
	d.server = &http.Server{Addr: d.addr, Handler: d.mux}
	go func() {
		slog.Debug("debug server starting", slog.String("addr", d.addr))
		err := d.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Debug("debug server is closed")
		} else {
			slog.Error("debug server error", slog.Any("error", err))
		}
	}()
	return nil
}

func (d *DebugManager) Stop(timeout time.Duration) error {
	d.mutx.Lock()
	defer d.mutx.Unlock()
	if d.server == nil {
		return fmt.Errorf("debug server not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := d.server.Shutdown(ctx)
	if err == nil {
		d.server = nil
		slog.Debug("debug server is stopped")
	} else {
		slog.Error("stop debug server error", slog.Any("error", err))
	}
	return err
}
