package middleware

import (
	"net/http"
)

func JSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.Method {
		case "GET":
			next.ServeHTTP(w, r)
		case "POST":
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte("415 - Header Content-Type incorrect"))
				return
			}
			next.ServeHTTP(w, r)
		case "PUT":
			next.ServeHTTP(w, r)
		case "DELETE":
			next.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 - Status method not allowed"))
			return
		}
	})
}
