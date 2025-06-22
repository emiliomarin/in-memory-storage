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
		log.Printf("ERROR: failed to get value for key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&lists.GetResponse[string]{
		List:      value.Value,
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

func (slc *stringListsController) Set(w http.ResponseWriter, r *http.Request) {
	var req lists.SetRequest[string]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
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
		log.Printf("ERROR: failed to set list for key %s: %v", req.Key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (slc *stringListsController) Update(w http.ResponseWriter, r *http.Request) {
	var req lists.UpdateRequest[string]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	if err := slc.store.Update(req.Key, req.List); err != nil {
		if err == storage.ErrNotFound || err == storage.ErrExpired {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to update list for key %s: %v", req.Key, err)
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
		log.Printf("ERROR: failed to remove list for key %s: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (slc *stringListsController) Push(w http.ResponseWriter, r *http.Request) {
	var req lists.PushRequest[string]
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

	if err := slc.store.Push(req.Key, req.Value); err != nil {
		if err == storage.ErrNotFound || err == storage.ErrExpired {
			http.Error(w, ErrKeyNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to push to list for key %s: %v", req.Key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (slc *stringListsController) Pop(w http.ResponseWriter, r *http.Request) {
	var req lists.PopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, ErrInvalidBody.Error(), http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		http.Error(w, ErrEmptyKey.Error(), http.StatusBadRequest)
		return
	}

	value, err := slc.store.Pop(req.Key)
	if err != nil {
		if err == storage.ErrNotFound || err == storage.ErrEmptyList || err == storage.ErrExpired {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to pop from list for key %s: %v", req.Key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&lists.PopResponse[string]{
		Value: value,
	})
	if err != nil {
		log.Printf("ERROR: failed to marshal response for key %s: %v", req.Key, err)
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(res); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
