package middleware

import (
	"net/http"
	"strings"
)

func Authentication(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    	// Perform authentication logic here
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authToken := strings.TrimPrefix(authHeader, "Bearer ")
		if authToken != "my-secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
        next.ServeHTTP(w, r)
    }
}