package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/tolikproh/banners-rotation/internal/config"
	"github.com/tolikproh/banners-rotation/internal/logger"
)

type Server struct {
	cfg *config.Config
	log *logger.Logger
	app Application
	srv *http.Server
}

func New(cfg *config.Config, log *logger.Logger, app Application) *Server {
	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return &Server{cfg: cfg, log: log, app: app, srv: srv}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("http server started", "address", s.srv.Addr)

	s.srv.Handler = s.initRoute()
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Debug("http server is shutting down...")

	return s.srv.Shutdown(ctx)
}
