package storage_test

import (
	"errors"
	"in-memory-storage/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListStore_Set(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key         string
		list        []string
		expectedVal []string
		expectedErr error
	}{
		"it should return an error if key already exists": {
			key:         "existing-key",
			list:        []string{"new-value"},
			expectedVal: []string{"val1", "val2"},
			expectedErr: errors.New("key already exists"),
		},
		"it should set the value": {
			key:         "new-key",
			list:        []string{"new-val1", "new-val2"},
			expectedVal: []string{"new-val1", "new-val2"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Set(tc.key, tc.list)
			assert.Equal(t, tc.expectedErr, err)

			val, err := store.Get(tc.key)
			assert.Nil(t, err)

			assert.Equal(t, tc.expectedVal, val)
		})
	}
}

func TestListStore_Set_Int(t *testing.T) {
	store := storage.NewListStore[int]()

	// Populate existing values
	err := store.Set("existing-key", []int{1, 2})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key          string
		list         []int
		expectedList []int
		expectedErr  error
	}{
		"it should return an error if key already exists": {
			key:          "existing-key",
			list:         []int{3, 4},
			expectedList: []int{1, 2},
			expectedErr:  errors.New("key already exists"),
		},
		"it should set the value": {
			key:          "new-key",
			list:         []int{3, 4},
			expectedList: []int{3, 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Set(tc.key, tc.list)
			assert.Equal(t, tc.expectedErr, err)

			val, err := store.Get(tc.key)
			assert.Nil(t, err)

			assert.Equal(t, tc.expectedList, val)
		})
	}
}

func TestListStore_Get(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key          string
		expectedList []string
		expectedErr  error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			expectedErr: errors.New("not found"),
		},
		"it should get the value": {
			key:          "existing-key",
			expectedList: []string{"val1", "val2"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			val, err := store.Get(tc.key)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedList, val)
		})
	}
}

func TestListStore_Update(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key          string
		list         []string
		expectedList []string
		expectedErr  error
	}{
		"it should return an error if key not found": {
			key:          "new-key",
			expectedList: []string{"val1", "val2"},
			expectedErr:  errors.New("not found"),
		},
		"it should update the value": {
			key:          "existing-key",
			list:         []string{"val3", "val4"},
			expectedList: []string{"val3", "val4"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Update(tc.key, tc.list)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				val, err := store.Get(tc.key)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedList, val)
			}
		})
	}
}

func TestListStore_Remove(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key            string
		expectedErr    error
		expectedGetErr error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			expectedErr: errors.New("not found"),
		},
		"it should remove the value": {
			key:            "existing-key",
			expectedGetErr: errors.New("not found"),
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

func TestListStore_Push(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key          string
		val          string
		expectedList []string
		expectedErr  error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			val:         "val3",
			expectedErr: errors.New("not found"),
		},
		"it should get the value": {
			key:          "existing-key",
			val:          "val3",
			expectedList: []string{"val1", "val2", "val3"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Push(tc.key, tc.val)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				val, err := store.Get(tc.key)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedList, val)
			}
		})
	}
}

func TestListStore_Pop(t *testing.T) {
	store := storage.NewListStore[string]()

	// Populate existing values
	err := store.Set("existing-key", []string{"val1", "val2"})
	assert.Nil(t, err)

	testCases := map[string]struct {
		key          string
		val          string
		expectedList []string
		expectedErr  error
	}{
		"it should return an error if key not found": {
			key:         "new-key",
			val:         "val3",
			expectedErr: errors.New("not found"),
		},
		"it should get the value": {
			key:          "existing-key",
			val:          "val3",
			expectedList: []string{"val1", "val2", "val3"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := store.Push(tc.key, tc.val)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				val, err := store.Get(tc.key)
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedList, val)
			}
		})
	}
}
