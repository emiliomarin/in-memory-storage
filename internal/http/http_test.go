package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"in-memory-storage/internal/http"
	"in-memory-storage/internal/strings"
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

func TestStringsSetEndpoint_E2E(t *testing.T) {
	stringStore := storage.NewStringStore()
	listsStore := storage.NewListStore[string]()

	stringsCtrl := http.NewStringsController(stringStore)
	listsCtrl := http.NewStringListsController(listsStore)

	validAPIKey := "valid-api-key"

	testCases := map[string]struct {
		apiKey         string
		requestBody    any
		expectedStatus int
		expectedError  error
		setupStore     func()
		verifyStore    func(t *testing.T)
	}{
		"should return 401 when wrong auth is provided": {
			apiKey:         "invalid-api-key",
			requestBody:    strings.SetRequest{Key: "test-key", Value: "test-value"},
			expectedStatus: gohttp.StatusUnauthorized,
			expectedError:  http.ErrUnauthorized,
		},
		"should return 400 when request body is invalid JSON": {
			apiKey:         validAPIKey,
			requestBody:    "invalid json",
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrInvalidBody,
		},
		"should return 400 when key is missing": {
			apiKey:         validAPIKey,
			requestBody:    strings.SetRequest{Key: "", Value: "test-value"},
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyKey,
		},
		"should return 400 when value is missing": {
			apiKey:         validAPIKey,
			requestBody:    strings.SetRequest{Key: "test-key", Value: ""},
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyValue,
		},
		"should return 409 when key already exists": {
			apiKey:         validAPIKey,
			requestBody:    strings.SetRequest{Key: "existing-key", Value: "new-value"},
			expectedStatus: gohttp.StatusConflict,
			expectedError:  http.ErrKeyAlreadyExists,
			setupStore: func() {
				_ = stringStore.Set("existing-key", "existing-value", 0)
			},
		},
		"should return 204 and store value when successful": {
			apiKey:         validAPIKey,
			requestBody:    strings.SetRequest{Key: "new-key", Value: "new-value", TTL: 60},
			expectedStatus: gohttp.StatusNoContent,
			verifyStore: func(t *testing.T) {
				value, err := stringStore.Get("new-key")
				assert.NoError(t, err)
				assert.Equal(t, "new-value", value.Value)
				assert.False(t, value.ExpiresAt.IsZero())
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupStore != nil {
				tc.setupStore()
			}

			srv, err := http.NewServer("8080", stringsCtrl, listsCtrl, validAPIKey)
			assert.NoError(t, err)

			var body []byte
			if str, ok := tc.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(gohttp.MethodPost, "/strings", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Add authorization header if API key is set
			if tc.apiKey != "" {
				req.Header.Set("Authorization", "Bearer "+tc.apiKey)
			}

			rr := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedError != nil {
				assert.Contains(t, rr.Body.String(), tc.expectedError.Error())
			}

			if tc.verifyStore != nil {
				tc.verifyStore(t)
			}
		})
	}
}
