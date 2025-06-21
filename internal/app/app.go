// Package app provides a server application that handles incoming requests.
// It initializes the HTTP server and starts it to start serving requests.
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"in-memory-storage/internal/http"
	"in-memory-storage/storage"
)

const defaultTimeout = 5 * time.Second

type Application struct {
	httpServer *http.Server
	port       string

	// value used to determine the gap of time
	// required for shutdown the application
	timeout time.Duration
}

// New creates a new Application instance with the provided configuration.
func New(port string) (*Application, error) {
	stringStore := storage.NewStringStore()
	stringListStore := storage.NewListStore[string]()

	stringsCtrl := http.NewStringsController(stringStore)
	stringsListCtrl := http.NewStringListsController(stringListStore)

	// Get API key from environment variable
	apiKey := os.Getenv("API_KEY")

	httpServer, err := http.NewServer(port, stringsCtrl, stringsListCtrl, apiKey)
	if err != nil {
		return nil, err
	}

	return &Application{
		httpServer: httpServer,
		timeout:    defaultTimeout,
		port:       port,
	}, nil
}

// Start runs the HTTP server and waits for a termination signal.
func (app *Application) Start() {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(quitCh)

	go func() {
		if err := app.httpServer.Start(); err != nil {
			log.Fatal("error starting http server: ", err)
		}
	}()

	fmt.Println("Server is running in port", app.port, "... Press Ctrl+C to stop.")

	<-quitCh
	fmt.Println(nil, "Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
	defer cancel()

	if err := app.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("error shutting down http server: ", err)
	}
	fmt.Println("Server stopped gracefully.")
}
