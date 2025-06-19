package storage

import "time"

type stringStore struct {
	store map[string]Value[string]
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
	return set(ss.store, key, val, ttl)
}

// Get will return the value for the given key.
// It will return an error if not found.
func (ss *stringStore) Get(key string) (Value[string], error) {
	return get(ss.store, key)
}

// Update will update the value for the given key.
// It will return an error if not found.
func (ss *stringStore) Update(key, val string) error {
	return update(ss.store, key, val)
}

// Remove will delete the value linked to the given key.
// It will return an error if not found.
func (ss *stringStore) Remove(key string) error {
	return remove(ss.store, key)
}
