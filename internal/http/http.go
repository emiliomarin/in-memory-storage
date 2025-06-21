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

	stringsController    StringsController
	stringListController ListsController
	authMiddleware       *AuthMiddleware
}

// NewServer creates a new HTTP server with the providided port.
// It returns an error if the port is missing.
func NewServer(
	port string,
	stringsController StringsController,
	stringListController ListsController,
	apiKey string,
) (*Server, error) {
	if port == "" {
		return nil, errors.New("missing port")
	}
	if stringsController == nil {
		return nil, errors.New("missing strings controller")
	}
	if stringListController == nil {
		return nil, errors.New("missing string list controller")
	}

	s := &Server{
		stringsController:    stringsController,
		stringListController: stringListController,
		authMiddleware:       NewAuthMiddleware(apiKey),
	}
	s.Server = &http.Server{
		Addr: ":" + port,
	}

	// Set the handler to the server's routes
	s.Handler = s.routes()

	return s, nil
}

// Start starts the HTTP server.
// It returns an error if the server is not initialized or if it fails to start.
func (s *Server) Start() error {
	if s.Server == nil {
		return errors.New("server not initialized")
	}

	return s.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
// It returns an error if the server is not initialized or if it fails to stop.
func (s *Server) Stop(ctx context.Context) error {
	if s.Server == nil {
		return errors.New("server not initialized")
	}

	return s.Shutdown(ctx)
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// String routes
	mux.HandleFunc("/strings", s.authMiddleware.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.stringsController.Set(w, r)
		case http.MethodGet:
			s.stringsController.Get(w, r)
		case http.MethodDelete:
			s.stringsController.Delete(w, r)
		case http.MethodPut:
			s.stringsController.Update(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// String list routes
	mux.HandleFunc("/lists/strings", s.authMiddleware.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.stringListController.Set(w, r)
		case http.MethodGet:
			s.stringListController.Get(w, r)
		case http.MethodDelete:
			s.stringListController.Delete(w, r)
		case http.MethodPut:
			s.stringListController.Update(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/lists/strings/push", s.authMiddleware.WithAuth(s.stringListController.Push))
	mux.HandleFunc("/lists/strings/pop", s.authMiddleware.WithAuth(s.stringListController.Pop))

	return mux
}
