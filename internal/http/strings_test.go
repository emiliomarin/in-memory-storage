package http_test

import (
	"bytes"
	"encoding/json"
	gohttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"in-memory-storage/internal/http"
	"in-memory-storage/internal/strings"
	"in-memory-storage/storage"

	"github.com/stretchr/testify/assert"
)

func TestStringsController_Set(t *testing.T) {
	store := storage.NewStringStore()
	controller := http.NewStringsController(store)

	// Populate store with data
	err := store.Set("existing-key", "existing-value", 0)
	assert.NoError(t, err)

	testCases := map[string]struct {
		key            string
		value          string
		ttl            int
		expectedStatus int
		expectedError  error
	}{
		"it should return an error if the key is missing": {
			key:            "",
			value:          "foo",
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyKey,
		},
		"it should return an error if the value is missing": {
			key:            "foo",
			value:          "",
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyValue,
		},
		"it should return an error if the key already exists": {
			key:            "existing-key",
			value:          "existing-value",
			expectedStatus: gohttp.StatusConflict,
			expectedError:  http.ErrKeyAlreadyExists,
		},
		"success": {
			key:            "foo",
			value:          "bar",
			expectedStatus: gohttp.StatusNoContent,
		},
		"success with ttl": {
			key:            "bar",
			value:          "foo",
			ttl:            60,
			expectedStatus: gohttp.StatusNoContent,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(strings.SetRequest{
				Key:   tc.key,
				Value: tc.value,
				TTL:   int64(tc.ttl),
			})
			req := httptest.NewRequest(gohttp.MethodPost, "/strings", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			controller.Set(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedError != nil {
				assert.Contains(t, rr.Body.String(), tc.expectedError.Error())
			} else {
				// Check that the value was set in the store if no error
				storedValue, err := store.Get(tc.key)
				assert.NoError(t, err)
				assert.Equal(t, tc.value, storedValue.Value)
			}
		})
	}
}

func TestStringsController_Get(t *testing.T) {
	store := storage.NewStringStore()
	controller := http.NewStringsController(store)

	// Populate store with data
	err := store.Set("existing-key", "existing-value", 60*time.Second)
	assert.NoError(t, err)
	err = store.Set("expired-key", "expired-value", 1*time.Millisecond)
	assert.NoError(t, err)

	// Ensure the expired key is actually expired
	time.Sleep(2 * time.Millisecond)

	testCases := map[string]struct {
		key               string
		ttl               int
		expectedStatus    int
		expectedError     error
		expectedValue     string
		expectedExpiresAt string
	}{
		"it should return an error if the key is missing": {
			key:            "",
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyKey,
		},
		"it should return an error if the key is not found": {
			key:            "random-key",
			expectedStatus: gohttp.StatusNotFound,
			expectedError:  http.ErrKeyNotFound,
		},
		"it should return an error if the key has expired": {
			key:            "expired-key",
			expectedStatus: gohttp.StatusNotFound,
			expectedError:  http.ErrKeyNotFound,
		},
		"success": {
			key:            "existing-key",
			expectedStatus: gohttp.StatusOK,
			expectedValue:  "existing-value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(gohttp.MethodGet, "/strings?key="+tc.key, nil)
			rr := httptest.NewRecorder()

			controller.Get(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedError != nil {
				assert.Contains(t, rr.Body.String(), tc.expectedError.Error())
			} else {
				var response strings.GetResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, response.Value)
			}
		})
	}
}
