package util

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type (
	HttpServer struct {
		handler http.Handler
		addr    string

		server       *http.Server
		listenerAddr string
	}
)

func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{handler: handler, addr: addr}
}

func (s *HttpServer) StartServer() error {
	if s.server != nil {
		return fmt.Errorf("server is already running")
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	server := &http.Server{Handler: s.handler}
	s.listenerAddr = listener.Addr().String()
	s.server = server
	go func() {
		slog.Debug("server starting", slog.String("addr", s.listenerAddr))
		err := s.server.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			slog.Debug("debug server is closed")
		} else {
			slog.Error("debug server error", slog.Any("error", err))
		}
	}()
	return nil
}

func (s *HttpServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.server = nil
	s.listenerAddr = ""
	return nil
}

func (s *HttpServer) GetListnerAddr() (string, error) {
	if s.server == nil || s.listenerAddr == "" {
		return "", fmt.Errorf("server is not running")
	}
	return s.listenerAddr, nil
}
