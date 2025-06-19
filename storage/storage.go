// Package storage defines an in-memory library to persist different data types.
// It defines and implements StringStore and ListStore interfaces with functionality to
// persist and retrieve string data and list of any type.
package storage

import (
	"time"
)

// Value is a generic struct that holds a value of type T and its expiration time.
// It can be used to store any type of data with an expiration time.
type Value[T any] struct {
	Value     T
	ExpiresAt *time.Time
}

// StringStore defines an interface for storing and retrieving string values.
type StringStore interface {
	Get(key string) (Value[string], error)
	Set(key string, val string) error
	Update(key string, val string) error
	Remove(key string) error
}

// ListStore defines an interface for storing and retrieving lists of any type.
type ListStore[T any] interface {
	Get(key string) (Value[[]T], error)
	Set(key string, list []T) error
	Update(key string, list []T) error
	Remove(key string) error
	Push(key string, val T) error
	Pop(key string) (T, error)
}

func set[T any](store map[string]Value[T], key string, val T) error {
	if _, ok := store[key]; ok {
		return ErrAlreadyExists
	}
	store[key] = Value[T]{Value: val}
	return nil
}

func get[T any](store map[string]Value[T], key string) (Value[T], error) {
	if _, ok := store[key]; !ok {
		var zero Value[T]
		return zero, ErrNotFound
	}
	return store[key], nil
}

func update[T any](store map[string]Value[T], key string, val T) error {
	if _, ok := store[key]; !ok {
		return ErrNotFound
	}

	store[key] = Value[T]{Value: val, ExpiresAt: store[key].ExpiresAt}
	return nil
}

func remove[T any](store map[string]Value[T], key string) error {
	if _, ok := store[key]; !ok {
		return ErrNotFound
	}
	delete(store, key)

	return nil
}
