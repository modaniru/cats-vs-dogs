package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) isAvailable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key != "cat" && key != "dog" {
			http.Error(w, "bad param", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) DataChanged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		res, err := s.GetAllValues(context.Background())
		if err != nil {
			slog.Error(err.Error())
			return
		}

		WriteToAllConnections(res)
	})
}

func (s *Server) Test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.RawQuery)
		next.ServeHTTP(w, r)
	})
}
