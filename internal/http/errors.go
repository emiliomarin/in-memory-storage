package http

import "errors"

var (
	// ErrEmptyKey is returned when the request contains an empty key.
	ErrEmptyKey = errors.New("key cannot be empty")
	// ErrEmptyValue is returned when the request contains an empty value.
	ErrEmptyValue = errors.New("value cannot be empty")
	// ErrKeyAlreadyExists is returned when trying to add an item that already exists in the store
	ErrKeyAlreadyExists = errors.New("key already exists")
	// ErrKeyNotFound is returned when the requested key is not found in the store.
	ErrKeyNotFound = errors.New("key not found")
	// ErrUnauthorized is returned when the request does not have a valid API key.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrInvalidBody is returned when the request body is invalid.
	ErrInvalidBody = errors.New("invalid request body")
)
