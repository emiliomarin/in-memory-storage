package http

import (
	"encoding/json"
	"in-memory-storage/internal/strings"
	"in-memory-storage/storage"
	"log"
	"net/http"
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
	var req strings.SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}
	if req.Value == "" {
		http.Error(w, ErrEmptyValue.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.store.Set(req.Key, req.Value, time.Duration(req.TTL)*time.Second); err != nil {
		if err == storage.ErrAlreadyExists {
			http.Error(w, ErrKeyAlreadyExists.Error(), http.StatusConflict)
			return
		}
		log.Printf("ERROR: failed to set value for key %s: %v", req.Key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
		log.Printf("ERROR: failed to get value for key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&strings.GetResponse{
		Value:     value.Value,
		ExpiresAt: value.ExpiresAt.Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("ERROR: failed to marshal response for key %s: %v", key, err)
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		log.Printf("failed to write response: %v", err)
	}
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
		log.Printf("ERROR: failed to remove key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (sc *stringController) Update(w http.ResponseWriter, r *http.Request) {
	var req strings.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}
	if req.Value == "" {
		http.Error(w, ErrEmptyValue.Error(), http.StatusBadRequest)
		return
	}

	if err := sc.store.Update(req.Key, req.Value); err != nil {
		if err == storage.ErrNotFound || err == storage.ErrExpired {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to update key %s: %v", req.Key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
