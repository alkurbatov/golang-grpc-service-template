// Package httpserver implements handy wrap around HTTP server
// to group common settings and tasks inside single entity.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

// NB (alkurbatov): Set reasonable timeouts, see:
// https://habr.com/ru/company/ispring/blog/560032/
const (
	defaultReadTimeout       = 5 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultIdleTimeout       = 120 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
)

// Server wraps HTTP server entity and handy means
// to simplify work with the entity.
type Server struct {
	ctx    context.Context
	Server *http.Server
}

// NewProdServer creates and initializes new instance of HTTP server
// preconfigured for production.
func NewProdServer(router http.Handler, address string) *Server {
	httpServer := &http.Server{
		Handler:           router,
		Addr:              address,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	return &Server{Server: httpServer}
}

// New creates and initializes new instance of HTTP server
// with default settings.
// Not recommended for production.
func New(router http.Handler, address string) *Server {
	httpServer := &http.Server{
		Handler:           router,
		Addr:              address,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	return &Server{Server: httpServer}
}

// Start launches the HTTP server. This is non-blocking call.
func (s *Server) Start() context.Context {
	if s.Running() {
		return s.ctx
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	s.ctx = ctx

	go func() {
		err := s.Server.ListenAndServe()

		cancel(err)
	}()

	return s.ctx
}

// Running returns true if HTTP server is running.
func (s *Server) Running() bool {
	return s.ctx != nil && s.ctx.Err() == nil
}

// Stop immediately stops the HTTP server and forcibly closes all connections.
func (s *Server) Stop() error {
	if !s.Running() {
		return nil
	}

	return s.Server.Close()
}
