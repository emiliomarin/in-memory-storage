package storage

import (
	"sync"
	"time"
)

type stringStore struct {
	store map[string]Value[string]
	// Mutex to handle concurrent access to memory
	mu sync.RWMutex
}

// NewStringStore initializes a new string store
func NewStringStore() StringStore {
	return &stringStore{
		store: map[string]Value[string]{}, // TODO: Here we could init with existing data
	}
}

// Set will store the given key/value pair.
// It will check if the key already exists and return an error.
func (ss *stringStore) Set(key, val string, ttl time.Duration) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return set(ss.store, key, val, ttl)
}

// Get will return the value for the given key.
// It will return an error if not found.
func (ss *stringStore) Get(key string) (Value[string], error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return get(ss.store, key)
}

// Update will update the value for the given key.
// It will return an error if not found.
func (ss *stringStore) Update(key, val string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return update(ss.store, key, val)
}

// Remove will delete the value linked to the given key.
// It will return an error if not found.
func (ss *stringStore) Remove(key string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return remove(ss.store, key)
}
