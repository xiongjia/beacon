package dbg

import (
	"expvar"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"sync"
	"time"

	"github.com/xiongjia/beacon/pkg/util"
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

		server *util.HttpServer
	}
)

const (
	// the default base url. The default pprof url is /debug/pprof/
	// This value can be changed via DebugManagerWithBaseUrl("/your-prefix/debug")
	DEFAULT_BASEURL = "/debug"
)

var (
	// the default listner address for debug server
	// This value can be set via DebugManagerWithAddr()
	DEFAULT_ADDR = net.JoinHostPort("127.0.0.1", "6060")
)

func newDebugMux() *http.ServeMux {
	mux := http.NewServeMux()
	// Debugger handler for pprof
	mux.HandleFunc("GET  /pprof/", pprof.Index)
	mux.HandleFunc("GET /pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /pprof/trace", pprof.Trace)
	// Debugger Handler for expvar
	mux.Handle("GET /vars", expvar.Handler())
	return mux
}

// NewDebugManager allocats and return a new [DebugManager].
//
// Debug Server options: DebugManagerWithBaseUrl(), DebugManagerWithAddr()
// Debug Server interfaces:
// - To start the debug server: [DebugManager.Start()]
// - To stop the debug server: [DebugManager.Stop(timemout)]
func NewDebugManager(opts ...DebugOption) *DebugManager {
	dbgOpts := dbgOption{
		baseUrl: DEFAULT_BASEURL,
		addr:    DEFAULT_ADDR,
	}
	for _, opt := range opts {
		opt(&dbgOpts)
	}
	return &DebugManager{
		addr: dbgOpts.addr,
		mux:  util.NewMainMuxGroup().Group(dbgOpts.baseUrl, newDebugMux()),
	}
}

// The default base url is "/debug"
// Example:  DebugManagerWithBaseUrl("/internal/debug") => The pprof url is "/internal/debug/pprof"
func DebugManagerWithBaseUrl(baseUrl string) DebugOption {
	return func(opt *dbgOption) {
		if baseUrl != "" {
			opt.baseUrl = baseUrl
		}
	}
}

// The default listner address is "127.0.0.1:6060"
func DebugManagerWithAddr(host string, port int16) DebugOption {
	return func(opt *dbgOption) {
		opt.addr = net.JoinHostPort(host, strconv.Itoa(int(port)))
	}
}

func (d *DebugManager) DebugServerAddr() (string, error) {
	d.mutx.Lock()
	defer d.mutx.Unlock()
	if d.server == nil {
		return "", fmt.Errorf("debug server not running")
	}
	return d.server.GetListnerAddr()
}

// DebugManager.Start() starts the debug http server.
// It will return error when the server already running.
func (d *DebugManager) Start() error {
	d.mutx.Lock()
	defer d.mutx.Unlock()

	if d.server != nil {
		slog.Error("debug server is already running", slog.String("addr", d.addr))
		return fmt.Errorf("debug server already running")
	}
	dbgServer := util.NewHttpServer(d.addr, d.mux)
	err := dbgServer.StartServer()
	if err != nil {
		slog.Error("debug server start", slog.Any("error", err))
		return err
	}
	d.server = dbgServer
	return nil
}

// DebugManager.Stop(timeout) stops the debug http server.
// It will return error when the server is not running.
func (d *DebugManager) Stop(timeout time.Duration) error {
	d.mutx.Lock()
	defer d.mutx.Unlock()
	if d.server == nil {
		return fmt.Errorf("debug server not running")
	}
	if err := d.server.Shutdown(timeout); err != nil {
		slog.Error("stop debug server error", slog.Any("error", err))
		return err
	}
	d.server = nil
	slog.Debug("debug server is stopped")
	return nil
}
