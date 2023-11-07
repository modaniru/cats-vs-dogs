package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) InitRouter() *chi.Mux{
	router := chi.NewRouter()
	sub := chi.NewRouter()
	{
		sub.Get("/", s.GetAll)
	}
	r := sub.With(s.isAvailable)
	{
		r.Get("/{key}", s.GetValue)
		r.Put("/{key}", s.IncreaseValue)
	}
	router.Mount("/v1", sub)
	return router
}

type GetValueResponse struct{
	Count int `json:"count"`
}

//validate middleware
func (s *Server) GetValue(w http.ResponseWriter, r *http.Request){
	key := chi.URLParam(r, "key")
	value, err := s.MyStorage.Get(r.Context(), key)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(GetValueResponse{Count: value})
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}


//validate middleware
func (s *Server) IncreaseValue(w http.ResponseWriter, r *http.Request){
	key := chi.URLParam(r, "key")
	err := s.MyStorage.Increase(r.Context(), key)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetAllResponse struct{
	CatCount int `json:"cat_count"`
	DogCount int `json:"dog_count"`
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request){
	cat, err := s.MyStorage.Get(r.Context(), "cat")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dog, err := s.MyStorage.Get(r.Context(), "dog")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(GetAllResponse{DogCount: dog, CatCount: cat})
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}