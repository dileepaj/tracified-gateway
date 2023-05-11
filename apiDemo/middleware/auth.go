package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/golang-jwt/jwt/v5"
)

// permissions
// 1= read
// 2= submit XDR
// "permissions": [1,2,4],
// "exp":1654330908
func Authentication(requiredPermission []int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utilities.HandleError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Check that the Authorization header starts with "Bearer".
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utilities.HandleError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key used to sign the token.
			return []byte("my-secret-key"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utilities.HandleError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			utilities.HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check that the token is valid and has not expired.
		if !token.Valid {
			utilities.HandleError(w, "token is expired", http.StatusUnauthorized)
			return
		}

		// Check the permissions in the token.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utilities.HandleError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		permissions, ok := claims["permissions"].([]interface{})
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if !checkPermissions(permissions, requiredPermission) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Add the token claims to the request context.
		ctx := context.WithValue(r.Context(), "tokenClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func checkPermissions(permissions []interface{}, requiredPermission []int) bool {
	// Create a map of required elements for faster lookups
	requiredMap := make(map[int]bool)
	for _, req := range requiredPermission {
		requiredMap[req] = true
	}

	// Check if all required elements are in the array
	for _, val := range permissions {
		if requiredMap[val.(int)] {
			delete(requiredMap, val.(int))
			if len(requiredMap) == 0 {
				return true
			}
		}
	}
	return false
}
