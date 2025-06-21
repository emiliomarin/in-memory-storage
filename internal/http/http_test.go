package http_test

import (
	"context"
	"errors"
	"testing"

	"in-memory-storage/internal/http"
	"in-memory-storage/storage"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	stringsCtrl := http.NewStringsController(storage.NewStringStore())
	listsCtrl := http.NewStringListsController(storage.NewListStore[string]())

	testCases := map[string]struct {
		port                 string
		stringsController    http.StringsController
		stringListController http.ListsController
		apiKey               string
		expectedErr          error
		wantNilSrv           bool
	}{
		"it should return an error if the port is missing": {
			port:                 "",
			stringsController:    stringsCtrl,
			stringListController: listsCtrl,
			apiKey:               "",
			expectedErr:          errors.New("missing port"),
			wantNilSrv:           true,
		},
		"it should return an error if strings controller is missing": {
			port:                 "8080",
			stringsController:    nil,
			stringListController: listsCtrl,
			apiKey:               "",
			expectedErr:          errors.New("missing strings controller"),
			wantNilSrv:           true,
		},
		"it should return an error if string list controller is missing": {
			port:                 "8080",
			stringsController:    stringsCtrl,
			stringListController: nil,
			apiKey:               "",
			expectedErr:          errors.New("missing string list controller"),
			wantNilSrv:           true,
		},
		"it should return no error if the server is initialized": {
			port:                 "8080",
			stringsController:    stringsCtrl,
			stringListController: listsCtrl,
			apiKey:               "",
			expectedErr:          nil,
			wantNilSrv:           false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			srv, err := http.NewServer(tc.port, tc.stringsController, tc.stringListController, tc.apiKey)
			assert.Equal(t, tc.expectedErr, err)
			if tc.wantNilSrv {
				assert.Nil(t, srv)
			} else {
				assert.NotNil(t, srv)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	stringsCtrl := http.NewStringsController(storage.NewStringStore())
	listsCtrl := http.NewStringListsController(storage.NewListStore[string]())

	t.Run("it should return an error if server is not initialized", func(t *testing.T) {
		srv := &http.Server{}
		srv.Server = nil

		err := srv.Start()
		expectedErr := errors.New("server not initialized")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("it should return error if ListenAndServe fails", func(t *testing.T) {
		srv, err := http.NewServer("invalid_port", stringsCtrl, listsCtrl, "")
		assert.NoError(t, err)

		// Overwrite Addr to an invalid value to force error
		srv.Addr = "invalid:port"
		err = srv.Start()
		assert.Error(t, err)
	})
}

func TestServer_Stop(t *testing.T) {
	stringsCtrl := http.NewStringsController(storage.NewStringStore())
	listsCtrl := http.NewStringListsController(storage.NewListStore[string]())

	t.Run("it should return an error if server is not initialized", func(t *testing.T) {
		srv := &http.Server{}
		srv.Server = nil

		err := srv.Stop(context.Background())
		expectedErr := errors.New("server not initialized")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("it should return no error if Shutdown succeeds", func(t *testing.T) {
		srv, err := http.NewServer("8080", stringsCtrl, listsCtrl, "")
		assert.NoError(t, err)

		go func() {
			_ = srv.Start()
		}()

		ctx := context.Background()
		err = srv.Stop(ctx)
		assert.NoError(t, err)
	})
}
