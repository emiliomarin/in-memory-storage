package storage

import (
	"errors"
	"sync"
	"time"
)

type listStore[T any] struct {
	store map[string]Value[[]T]
	// Mutex to handle concurrent access to memory
	mu sync.RWMutex
}

// NewListStore initializes a list store for the given data type
func NewListStore[T any]() ListStore[T] {
	return &listStore[T]{
		store: map[string]Value[[]T]{}, // TODO: Here we could init with existing data
	}
}

// Set will store the given key/value pair.
// It will check if the key already exists and return an error.
func (ls *listStore[T]) Set(key string, list []T, ttl time.Duration) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return set(ls.store, key, list, ttl)
}

// Get will return the value for the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Get(key string) (*Value[[]T], error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	value, err := get(ls.store, key)
	if err != nil {
		return nil, err
	}

	// If the value has an expiration time and it is in the past, remove it
	// and return an error indicating it has expired.
	if !value.ExpiresAt.IsZero() && value.ExpiresAt.Before(time.Now()) {
		if err := remove(ls.store, key); err != nil {
			return nil, errors.New("failed to remove expired key: " + err.Error())
		}
		return nil, ErrExpired
	}

	return &value, nil
}

// Update will update the value for the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Update(key string, list []T) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return update(ls.store, key, list)
}

// Remove will delete the value linked to the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Remove(key string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return remove(ls.store, key)
}

// Push will add the given value to the existing list.
// It will return an error if the list is not found
func (ls *listStore[T]) Push(key string, val T) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if _, ok := ls.store[key]; !ok {
		return ErrNotFound
	}

	v := ls.store[key]
	v.Value = append(v.Value, val)
	ls.store[key] = v

	return nil
}

// Pop will retrieve and remove the first item from the list. Applying FIFO.
// It will check that the list exists and that it's not empty.
func (ls *listStore[T]) Pop(key string) (T, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var zero T

	if _, ok := ls.store[key]; !ok {
		return zero, ErrNotFound
	}

	if len(ls.store[key].Value) == 0 {
		return zero, ErrEmptyList
	}

	val := ls.store[key].Value[0]

	newList := ls.store[key]
	newList.Value = ls.store[key].Value[1:]
	ls.store[key] = newList

	return val, nil
}
