# Storage Library API Documentation

The `storage` package provides an in-memory data storage solution for Go applications. It offers thread-safe stores for strings and generic lists, with support for Time-To-Live (TTL) expiration.

## Table of Contents

-   [Errors](#errors)
-   [Value Struct](#value-struct)
-   [StringStore Interface](#stringstore-interface)
    -   [NewStringStore()](#newstringstore)
    -   [Set()](#set)
    -   [Get()](#get)
    -   [Update()](#update)
    -   [Remove()](#remove)
-   [ListStore Interface](#liststore-interface)
    -   [NewListStore()](#newliststore)
    -   [Set() (List)](#set-list)
    -   [Get() (List)](#get-list)
    -   [Update() (List)](#update-list)
    -   [Remove() (List)](#remove-list)
    -   [Push()](#push)
    -   [Pop()](#pop)

---

## Errors

The package defines the following error variables:

-   `ErrNotFound`: Returned when a requested item is not found in the store.
-   `ErrAlreadyExists`: Returned when trying to add an item that already exists in the store.
-   `ErrEmptyList`: Returned when trying to `Pop` an item from an empty list.
-   `ErrExpired`: Returned when trying to access an item whose TTL has expired.

---

## Value Struct

`Value[T any]` is a generic struct that holds a value and its expiration time.

```go
type Value[T any] struct {
    Value     T
    ExpiresAt time.Time
}
```

-   `Value`: The data of type `T` being stored.
-   `ExpiresAt`: The time at which the value expires. If `ExpiresAt` is the zero value, the item does not expire.

---

## StringStore Interface

An interface for storing and retrieving string values.

### `NewStringStore()`

Initializes a new `StringStore`.

-   **Signature:** `func NewStringStore() StringStore`
-   **Returns:** A new instance of `StringStore`.

### `Set()`

Stores a key-value pair with an optional TTL.

-   **Signature:** `func (ss *stringStore) Set(key, val string, ttl time.Duration) error`
-   **Parameters:**
    -   `key` (string): The key for the value.
    -   `val` (string): The string value to store.
    -   `ttl` (time.Duration): The time-to-live for the value. If `0`, the value never expires.
-   **Returns:** `ErrAlreadyExists` if the key is already in the store, otherwise `nil`.

### `Get()`

Retrieves a value by its key.

-   **Signature:** `func (ss *stringStore) Get(key string) (*Value[string], error)`
-   **Parameters:**
    -   `key` (string): The key of the value to retrieve.
-   **Returns:** A pointer to a `Value[string]` struct and `nil` error on success. Returns `ErrNotFound` if the key doesn't exist or `ErrExpired` if the key has expired (and removes it).

### `Update()`

Updates the value for an existing key.

-   **Signature:** `func (ss *stringStore) Update(key, val string) error`
-   **Parameters:**
    -   `key` (string): The key of the value to update.
    -   `val` (string): The new string value.
-   **Returns:** `ErrNotFound` if the key doesn't exist, otherwise `nil`.

### `Remove()`

Deletes a key-value pair from the store.

-   **Signature:** `func (ss *stringStore) Remove(key string) error`
-   **Parameters:**
    -   `key` (string): The key of the value to remove.
-   **Returns:** `ErrNotFound` if the key doesn't exist, otherwise `nil`.

---

## ListStore Interface

A generic interface for storing and retrieving lists.

### `NewListStore()`

Initializes a new generic `ListStore`.

-   **Signature:** `func NewListStore[T any]() ListStore[T]`
-   **Returns:** A new instance of `ListStore[T]` for the specified type `T`.

### `Set()` (List)

Stores a list with an optional TTL.

-   **Signature:** `func (ls *listStore[T]) Set(key string, list []T, ttl time.Duration) error`
-   **Parameters:**
    -   `key` (string): The key for the list.
    -   `list` ([]T): The list to store.
    -   `ttl` (time.Duration): The time-to-live for the list. If `0`, it never expires.
-   **Returns:** `ErrAlreadyExists` if the key is already in the store, otherwise `nil`.

### `Get()` (List)

Retrieves a list by its key.

-   **Signature:** `func (ls *listStore[T]) Get(key string) (*Value[[]T], error)`
-   **Parameters:**
    -   `key` (string): The key of the list to retrieve.
-   **Returns:** A pointer to a `Value[[]T]` struct and `nil` error on success. Returns `ErrNotFound` if the key doesn't exist or `ErrExpired` if the list has expired (and removes it).

### `Update()` (List)

Updates the list for an existing key.

-   **Signature:** `func (ls *listStore[T]) Update(key string, list []T) error`
-   **Parameters:**
    -   `key` (string): The key of the list to update.
    -   `list` ([]T): The new list.
-   **Returns:** `ErrNotFound` if the key doesn't exist, otherwise `nil`.

### `Remove()` (List)

Deletes a list from the store.

-   **Signature:** `func (ls *listStore[T]) Remove(key string) error`
-   **Parameters:**
    -   `key` (string): The key of the list to remove.
-   **Returns:** `ErrNotFound` if the key doesn't exist, otherwise `nil`.

### `Push()`

Adds a value to the end of a list.

-   **Signature:** `func (ls *listStore[T]) Push(key string, val T) error`
-   **Parameters:**
    -   `key` (string): The key of the list.
    -   `val` (T): The value to add.
-   **Returns:** `ErrNotFound` if the list doesn't exist, otherwise `nil`.

### `Pop()`

Removes and returns the first item from a list (FIFO).

-   **Signature:** `func (ls *listStore[T]) Pop(key string) (T, error)`
-   **Parameters:**
    -   `key` (string): The key of the list.
-   **Returns:** The first value from the list and `nil` error. Returns `ErrNotFound` if the list doesn't exist or `ErrEmptyList` if the list is empty. 
