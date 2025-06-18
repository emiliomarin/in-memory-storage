package http_test

import (
	"context"
	"errors"
	"testing"

	"in-memory-storage/internal/http"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("it should return an error if the port is missing", func(t *testing.T) {
		_, err := http.NewServer("")
		expectedErr := errors.New("missing port")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("it should return no error if the server is initialized", func(t *testing.T) {
		srv, err := http.NewServer("8080")
		assert.NoError(t, err)
		assert.NotNil(t, srv)
	})
}

func TestServer_Start(t *testing.T) {
	t.Run("it should return an error if server is not initialized", func(t *testing.T) {
		srv := &http.Server{}
		srv.Server = nil

		err := srv.Start()
		expectedErr := errors.New("server not initialized")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("it should return error if ListenAndServe fails", func(t *testing.T) {
		srv, err := http.NewServer("invalid_port")
		assert.NoError(t, err)

		// Overwrite Addr to an invalid value to force error
		srv.Server.Addr = "invalid:port"
		err = srv.Start()
		assert.Error(t, err)
	})
}

func TestServer_Stop(t *testing.T) {
	t.Run("it should return an error if server is not initialized", func(t *testing.T) {
		srv := &http.Server{}
		srv.Server = nil

		err := srv.Stop(context.Background())
		expectedErr := errors.New("server not initialized")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("it should return no error if Shutdown succeeds", func(t *testing.T) {
		srv, err := http.NewServer("8080")
		assert.NoError(t, err)

		go func() {
			_ = srv.Start()
		}()

		ctx := context.Background()
		err = srv.Stop(ctx)
		assert.NoError(t, err)
	})
}
