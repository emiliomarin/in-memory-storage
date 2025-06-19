package storage

import "time"

type listStore[T any] struct {
	store map[string]Value[[]T]
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
	return set(ls.store, key, list, ttl)
}

// Get will return the value for the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Get(key string) (Value[[]T], error) {
	return get(ls.store, key)
}

// Update will update the value for the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Update(key string, list []T) error {
	return update(ls.store, key, list)
}

// Remove will delete the value linked to the given key.
// It will return an error if the list is not found.
func (ls *listStore[T]) Remove(key string) error {
	return remove(ls.store, key)
}

// Push will add the given value to the existing list.
// It will return an error if the list is not found
func (ls *listStore[T]) Push(key string, val T) error {
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
