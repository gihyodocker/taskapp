package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

type HTTPServer struct {
	port        int
	mux         *chi.Mux
	server      *http.Server
	gracePeriod time.Duration
}

type Option func(*HTTPServer)

func WithGracePeriod(d time.Duration) Option {
	return func(s *HTTPServer) {
		s.gracePeriod = d
	}
}

func NewHTTPServer(port int, opts ...Option) *HTTPServer {
	r := chi.NewRouter()
	s := &HTTPServer{
		port: port,
		mux:  r,
		server: &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", port),
		},
		gracePeriod: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *HTTPServer) Get(pattern string, handlerFn http.HandlerFunc) {
	s.mux.Get(pattern, handlerFn)
}

func (s *HTTPServer) Post(pattern string, handlerFn http.HandlerFunc) {
	s.mux.Post(pattern, handlerFn)
}

func (s *HTTPServer) Put(pattern string, handlerFn http.HandlerFunc) {
	s.mux.Put(pattern, handlerFn)
}

func (s *HTTPServer) Delete(pattern string, handlerFn http.HandlerFunc) {
	s.mux.Delete(pattern, handlerFn)
}

func (s *HTTPServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *HTTPServer) Serve(ctx context.Context) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		doneCh <- s.run()
	}()

	<-ctx.Done()
	s.stop()
	return <-doneCh
}

func (s *HTTPServer) run() error {
	slog.Info(fmt.Sprintf("http api is running on %d", s.port))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to listen and serve http api", err)
		return err
	}
	return nil
}

func (s *HTTPServer) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.gracePeriod)
	defer cancel()
	slog.Info("stopping http api")
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown http api", err)
		return err
	}
	return nil
}
