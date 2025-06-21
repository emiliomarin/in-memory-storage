package http

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	apiKey string
}

func NewAuthMiddleware(apiKey string) *AuthMiddleware {
	return &AuthMiddleware{
		apiKey: apiKey,
	}
}

func (am *AuthMiddleware) WithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !am.authenticate(r) {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func (am *AuthMiddleware) authenticate(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return false
	}

	// Check for Bearer token
	if parts[0] == "Bearer" && parts[1] == am.apiKey {
		return true
	}

	return false
}
