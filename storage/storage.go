// Package storage defines an in-memory library to persist different data types.
// It defines and implements StringStore and ListStore interfaces with functionality to
// persist and retrieve string data and list of any type.
package storage

import "errors"

type StringStore interface {
	Get(key string) (string, error)
	Set(key, val string) error // TODO: Set TTL
	Update(key, val string) error
	Remove(key string) error
}

type ListStore[T any] interface {
	Get(key string) ([]T, error)
	Set(key string, list []T) error // TODO: Set TTL
	Update(key string, list []T) error
	Remove(key string) error
	Push(key string, val T) error
	Pop(key string) (T, error)
}

func set[T any](store map[string]T, key string, val T) error {
	if _, ok := store[key]; ok {
		return errors.New("key already exists")
	}
	store[key] = val
	return nil
}

func get[T any](store map[string]T, key string) (T, error) {
	if _, ok := store[key]; !ok {
		var zero T
		return zero, ErrNotFound
	}
	return store[key], nil
}

func update[T any](store map[string]T, key string, val T) error {
	if _, ok := store[key]; !ok {
		return ErrNotFound
	}

	store[key] = val
	return nil
}

func remove[T any](store map[string]T, key string) error {
	if _, ok := store[key]; !ok {
		return ErrNotFound
	}
	delete(store, key)

	return nil
}
