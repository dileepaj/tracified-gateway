package middleware

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/utilities"
)

func HeaderReader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.Method {
		case http.MethodGet:
			next.ServeHTTP(w, r)
			break
		case http.MethodPost:
			if r.Header.Get("Content-Type") != "application/json" {
				utilities.HandleError(w, "Header Content-Type incorrect", http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
			break
		case http.MethodPut:
			if r.Header.Get("Content-Type") != "application/json" {
				utilities.HandleError(w, "Header Content-Type incorrect", http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
			break
		case http.MethodDelete:
			next.ServeHTTP(w, r)
			break
		default:
			utilities.HandleError(w, "Status method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
