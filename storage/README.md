## Storage Library API

This library provides in-memory storage for strings and lists of any type, with simple CRUD operations.

### StringStore

The `StringStore` interface provides methods to store, retrieve, update, and remove string values by key.

#### Methods

- **Set(key, val string) error**
Stores a new key/value pair. Returns an error if the key already exists.

- **Get(key string) (string, error)**
Retrieves the value for the given key. Returns an error if the key is not found.

- **Update(key, val string) error**
Updates the value for the given key. Returns an error if the key is not found.

- **Remove(key string) error**
Deletes the value for the given key. Returns an error if the key is not found.

#### Usage

```go
import "github.com/emiliomarin/in-memory-storage/storage"

store := storage.NewStringStore()
err := store.Set("foo", "bar")
val, err := store.Get("foo")
err = store.Update("foo", "baz")
err = store.Remove("foo")
```

---

### ListStore

The `ListStore[T]` interface provides methods to store, retrieve, update, and remove lists of any type, as well as push and pop elements.

#### Methods

- **Set(key string, list []T) error**
Stores a new key/list pair. Returns an error if the key already exists.

- **Get(key string) ([]T, error)**
Retrieves the list for the given key. Returns an error if the key is not found.

- **Update(key string, list []T) error**
Updates the list for the given key. Returns an error if the key is not found.

- **Remove(key string) error**
Deletes the list for the given key. Returns an error if the key is not found.

- **Push(key string, val T) error**
Appends a value to the list for the given key. Returns an error if the key is not found.

- **Pop(key string) (T, error)**
Removes and returns the first element (FIFO) from the list for the given key. Returns an error if the key is not found or the list is empty.

#### Usage

```go
import "github.com/emiliomarin/in-memory-storage/storage"

listStore := storage.NewListStore[int]()
err := listStore.Set("numbers", []int{1, 2, 3})
err = listStore.Push("numbers", 4)
val, err := listStore.Pop("numbers")
list, err := listStore.Get("numbers")
err = listStore.Update("numbers", []int{10, 20})
err = listStore.Remove("numbers")
```
