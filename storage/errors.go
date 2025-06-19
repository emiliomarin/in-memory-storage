package storage

import "errors"

var (
	// ErrNotFound is returned when the requested item is not found in the store.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when trying to add an item that already exists in the store
	ErrAlreadyExists = errors.New("already exists")
	// ErrEmptyList is returned when trying to pop an item from an empty list
	ErrEmptyList = errors.New("list is empty")
	// ErrExpired is returned when trying to access an expired item
	ErrExpired = errors.New("expired")
)
