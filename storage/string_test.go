package storage_test

import (
	"in-memory-storage/storage"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringStore_Set(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value", 0)
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		val         string
		expectedVal *storage.Value[string]
		expectedErr error
	}{
		"it should return an error if key already exists": {
			key:         "existing-key",
			val:         "new-value",
			expectedVal: &storage.Value[string]{Value: "existing-value"},
			expectedErr: storage.ErrAlreadyExists,
		},
		"it should set the value": {
			key:         "new-key",
			val:         "new-value",
			expectedVal: &storage.Value[string]{Value: "new-value"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Set(tc.key, tc.val, 0)
			assert.Equal(t, tc.expectedErr, err)

			val, err := store.Get(tc.key)
			assert.Nil(t, err)

			assert.Equal(t, tc.expectedVal, val)
		})
	}
}

func TestStringStore_Get(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value", 0)
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		expectedVal *storage.Value[string]
		expectedErr error
		setup       func()
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			expectedErr: storage.ErrNotFound,
		},
		"it should get the value": {
			key:         "existing-key",
			expectedVal: &storage.Value[string]{Value: "existing-value"},
		},
		"it should return an error if the key has expired": {
			key:         "expired-key",
			expectedErr: storage.ErrExpired,
			setup: func() {
				// Set an expired value
				err := store.Set("expired-key", "expired-value", time.Millisecond)
				assert.Nil(t, err)
				time.Sleep(2 * time.Millisecond) // Ensure the value is expired
			},
		},
		"it should return the expected value if the key hasn't expired": {
			key:         "ttl-valid-key",
			expectedVal: &storage.Value[string]{Value: "valid-value"},
			setup: func() {
				// Set a valid value with TTL
				err := store.Set("ttl-valid-key", "valid-value", time.Second)
				assert.Nil(t, err)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			val, err := store.Get(tc.key)
			assert.Equal(t, tc.expectedErr, err)
			if val != nil {
				assert.Equal(t, tc.expectedVal.Value, val.Value)
			}
		})
	}
}

func TestStringStore_Update(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value", 0)
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		val         string
		expectedVal storage.Value[string]
		expectedErr error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			val:         "new-value",
			expectedErr: storage.ErrNotFound,
		},
		"it should update the value": {
			key:         "existing-key",
			val:         "new-value",
			expectedVal: storage.Value[string]{Value: "new-value"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Update(tc.key, tc.val)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				val, err := store.Get(tc.key)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedVal, val)
			}
		})
	}
}

func TestStringStore_Remove(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value", 0)
	assert.Nil(t, err)

	testCases := map[string]struct {
		key            string
		expectedErr    error
		expectedGetErr error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			expectedErr: storage.ErrNotFound,
		},
		"it should remove the value": {
			key:            "existing-key",
			expectedGetErr: storage.ErrNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Remove(tc.key)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				_, err := store.Get(tc.key)
				assert.Equal(t, tc.expectedGetErr, err)
			}
		})
	}
}

func TestStringStore_ConcurrentSet(t *testing.T) {
	store := storage.NewStringStore()
	const n = 100
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key-" + strconv.Itoa(i)
			val := "val-" + strconv.Itoa(i)
			err := store.Set(key, val, 0)
			assert.Nil(t, err)
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		key := "key-" + strconv.Itoa(i)
		val, err := store.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, "val-"+strconv.Itoa(i), val.Value)
	}
}

func TestStringStore_ConcurrentUpdate(t *testing.T) {
	store := storage.NewStringStore()
	key := "existing-key"
	err := store.Set(key, "initial", 0)
	assert.Nil(t, err)

	const n = 100
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val := "val-" + strconv.Itoa(i)
			err := store.Update(key, val)
			assert.Nil(t, err)
		}(i)
	}
	wg.Wait()

	// Cannot guarantee the final value due to concurrent updates, but we can check that it contains the expected pattern
	val, err := store.Get(key)
	assert.Nil(t, err)
	assert.Contains(t, val.Value, "val-")
}
