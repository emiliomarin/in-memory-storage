package storage_test

import (
	"in-memory-storage/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringStore_Set(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value")
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		val         string
		expectedVal string
		expectedErr error
	}{
		"it should return an error if key already exists": {
			key:         "existing-key",
			val:         "new-value",
			expectedVal: "existing-value",
			expectedErr: storage.ErrAlreadyExists,
		},
		"it should set the value": {
			key:         "new-key",
			val:         "new-value",
			expectedVal: "new-value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Set(tc.key, tc.val)
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
	err := store.Set("existing-key", "existing-value")
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		expectedVal string
		expectedErr error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			expectedErr: storage.ErrNotFound,
		},
		"it should get the value": {
			key:         "existing-key",
			expectedVal: "existing-value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			val, err := store.Get(tc.key)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedVal, val)
		})
	}
}

func TestStringStore_Update(t *testing.T) {
	store := storage.NewStringStore()

	// Populate existing values
	err := store.Set("existing-key", "existing-value")
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		val         string
		expectedVal string
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
			expectedVal: "new-value",
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
	err := store.Set("existing-key", "existing-value")
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
