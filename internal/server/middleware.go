package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) isAvailable(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key != "cat" && key != "dog"{
			http.Error(w, "bad param", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}