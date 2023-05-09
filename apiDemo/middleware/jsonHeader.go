package middleware

import (
	"net/http"
)

func HeaderReader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.Method {
		case http.MethodGet:
			next.ServeHTTP(w, r)
		case http.MethodPost:
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte("415 - Header Content-Type incorrect"))
				return
			}
			next.ServeHTTP(w, r)
		case http.MethodPut:
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte("415 - Header Content-Type incorrect"))
				return
			}
			next.ServeHTTP(w, r)
		case http.MethodDelete:
			next.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 - Status method not allowed"))
			return
		}
	})
}
