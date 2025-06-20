package http

import (
	"encoding/json"
	"in-memory-storage/internal/strings"
	"in-memory-storage/storage"
	"net/http"
	"strconv"
	"time"
)

type StringsController interface {
	Set(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

func NewStringsController(store storage.StringStore) StringsController {
	return &stringController{store: store}
}

type stringController struct {
	store storage.StringStore
}

func (sc *stringController) Set(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}
	value := r.URL.Query().Get("value")
	if value == "" {
		http.Error(w, ErrEmptyValue.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.store.Set(key, value, getTTLFromRequest(r)); err != nil {
		if err == storage.ErrAlreadyExists {
			http.Error(w, ErrKeyAlreadyExists.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content

	// TODO: Here we could also return the response to have the expires at field
}

func (sc *stringController) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	value, err := sc.store.Get(key)
	if err != nil {
		if err == storage.ErrNotFound || err == storage.ErrExpired {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&strings.GetResponse{
		Value:     value.Value,
		ExpiresAt: value.ExpiresAt.Format(time.RFC3339),
	})
	if err != nil {
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (sc *stringController) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}
	if err := sc.store.Remove(key); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (sc *stringController) Update(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}
	value := r.URL.Query().Get("value")
	if value == "" {
		http.Error(w, ErrEmptyValue.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.store.Update(key, value); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content

	// TODO: Here we could also return the response to have the expires at field
}

// getTTLFromRequest extracts an optional "ttl" query parameter (in seconds) from the request.
// Returns a time.Duration. If not present or invalid, returns 0.
func getTTLFromRequest(r *http.Request) time.Duration {
	ttlStr := r.URL.Query().Get("ttl")
	if ttlStr == "" {
		return 0
	}
	ttlSec, err := strconv.Atoi(ttlStr)
	if err != nil || ttlSec < 0 {
		return 0
	}
	return time.Duration(ttlSec) * time.Second
}
