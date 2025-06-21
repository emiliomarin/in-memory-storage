package http

import (
	"encoding/json"
	"in-memory-storage/internal/lists"
	"in-memory-storage/storage"
	"log"
	"net/http"
	"time"
)

type ListsController interface {
	Get(w http.ResponseWriter, r *http.Request)
	Set(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Push(w http.ResponseWriter, r *http.Request)
	Pop(w http.ResponseWriter, r *http.Request)
}

func NewStringListsController(store storage.ListStore[string]) ListsController {
	return &stringListsController{store: store}
}

type stringListsController struct {
	store storage.ListStore[string]
}

func (slc *stringListsController) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	value, err := slc.store.Get(key)
	if err != nil {
		if err == storage.ErrNotFound || err == storage.ErrExpired {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&lists.GetResponse[string]{
		List:      value.Value,
		ExpiresAt: value.ExpiresAt.Format(time.RFC3339),
	})
	if err != nil {
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (slc *stringListsController) Set(w http.ResponseWriter, r *http.Request) {
	var req lists.SetRequest[string]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	if err := slc.store.Set(req.Key, req.List, time.Duration(req.TTL)*time.Second); err != nil {
		if err == storage.ErrAlreadyExists {
			http.Error(w, ErrKeyAlreadyExists.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (slc *stringListsController) Update(w http.ResponseWriter, r *http.Request) {
	var req lists.UpdateRequest[string]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	if err := slc.store.Update(req.Key, req.List); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (slc *stringListsController) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	if err := slc.store.Remove(key); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (slc *stringListsController) Push(w http.ResponseWriter, r *http.Request) {
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

	if err := slc.store.Push(key, value); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (slc *stringListsController) Pop(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	value, err := slc.store.Pop(key)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&lists.PopResponse[string]{
		Value: value,
	})
	if err != nil {
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(res); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
