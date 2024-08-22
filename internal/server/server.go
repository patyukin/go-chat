package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	h          http.Handler
}

func New(h http.Handler) *Server {
	return &Server{
		h: h,
	}
}

func (s *Server) Run(addr string) error {
	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        s.h,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
