// Package http provides a HTTP server implementation.
// It defines the available routes and handles incoming requests.
// It also provides a way to start and stop the server gracefully.
package http

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	*http.Server
}

// NewServer creates a new HTTP server with the providided port.
// It returns an error if the port is missing.
func NewServer(port string) (*Server, error) {
	if port == "" {
		return nil, errors.New("missing port")
	}

	s := &Server{}
	s.Server = &http.Server{
		Addr: ":" + port,
	}

	return s, nil
}

// Start starts the HTTP server.
// It returns an error if the server is not initialized or if it fails to start.
func (s *Server) Start() error {
	if s.Server == nil {
		return errors.New("server not initialized")
	}

	return s.Server.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
// It returns an error if the server is not initialized or if it fails to stop.
func (s *Server) Stop(ctx context.Context) error {
	if s.Server == nil {
		return errors.New("server not initialized")
	}

	return s.Server.Shutdown(ctx)
}
