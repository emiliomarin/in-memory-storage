package http_test

import (
	"bytes"
	"encoding/json"
	gohttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"in-memory-storage/internal/http"
	"in-memory-storage/internal/lists"
	"in-memory-storage/storage"

	"github.com/stretchr/testify/assert"
)

func TestListsController_Set(t *testing.T) {
	store := storage.NewListStore[string]()
	controller := http.NewStringListsController(store)

	// Populate store with data
	err := store.Set("existing-key", []string{"a", "b"}, 0)
	assert.NoError(t, err)

	testCases := map[string]struct {
		key            string
		list           []string
		ttl            int64
		expectedStatus int
		expectedError  error
	}{
		"it should return an error if the key is missing": {
			key:            "",
			list:           []string{"foo"},
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyKey,
		},
		"it should return an error if the key already exists": {
			key:            "existing-key",
			list:           []string{"c", "d"},
			expectedStatus: gohttp.StatusConflict,
			expectedError:  http.ErrKeyAlreadyExists,
		},
		"success": {
			key:            "foo",
			list:           []string{"bar", "baz"},
			expectedStatus: gohttp.StatusNoContent,
		},
		"success with ttl": {
			key:            "bar",
			list:           []string{"foo"},
			ttl:            60,
			expectedStatus: gohttp.StatusNoContent,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(lists.SetRequest[string]{
				Key:  tc.key,
				List: tc.list,
				TTL:  tc.ttl,
			})
			req := httptest.NewRequest(gohttp.MethodPost, "/lists/strings", bytes.NewReader(payload))
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
				assert.ElementsMatch(t, tc.list, storedValue.Value)
			}
		})
	}
}

func TestListsController_Get(t *testing.T) {
	store := storage.NewListStore[string]()
	controller := http.NewStringListsController(store)

	// Populate store with data
	err := store.Set("existing-key", []string{"a", "b"}, 60*time.Second)
	assert.NoError(t, err)
	err = store.Set("expired-key", []string{"expired-value"}, 1*time.Millisecond)
	assert.NoError(t, err)

	// Ensure the expired key is actually expired
	time.Sleep(2 * time.Millisecond)

	testCases := map[string]struct {
		key            string
		expectedStatus int
		expectedError  error
		expectedList   []string
	}{
		"it should return an error if the key is missing": {
			key:            "",
			expectedStatus: gohttp.StatusBadRequest,
			expectedError:  http.ErrEmptyKey,
		},
		"it should return an error if the key does not exist": {
			key:            "non-existing-key",
			expectedStatus: gohttp.StatusNotFound,
			expectedError:  storage.ErrNotFound,
		},
		"it should return an error if the key has expired": {
			key:            "expired-key",
			expectedStatus: gohttp.StatusNotFound,
			expectedError:  http.ErrKeyNotFound,
		},
		"success": {
			key:            "existing-key",
			expectedStatus: gohttp.StatusOK,
			expectedList:   []string{"a", "b"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(gohttp.MethodGet, "/lists?key="+tc.key, nil)
			rr := httptest.NewRecorder()

			controller.Get(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedError != nil {
				assert.Contains(t, rr.Body.String(), tc.expectedError.Error())
			} else {
				var response lists.GetResponse[string]
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedList, response.List)
			}
		})
	}
}
