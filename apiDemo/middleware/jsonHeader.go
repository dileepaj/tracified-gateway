package middleware

import "net/http"

func FilterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("405 - Header Content-Type incorrect"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}
