package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func (s *Server) InitRouter() *chi.Mux {
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
	router.HandleFunc("/w", s.Websocket)
	router.Mount("/v1", sub)
	return router
}

type GetValueResponse struct {
	Count int `json:"count"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//посмотреть
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReadWebsocket struct {
	Candidate string `json:"candidate"`
}

var candidates = map[string]bool{"dog": true, "cat": true}
var connections = map[*websocket.Conn]bool{}

func (s *Server) Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	connections[conn] = true

	defer func() {
		conn.Close()
		delete(connections, conn)
	}()

	for {
		candidate := &ReadWebsocket{}
		err := conn.ReadJSON(&candidate)
		if err != nil {
			log.Println(err.Error())
			break
		}

		if !candidates[candidate.Candidate] {
			log.Println("candidate not found")
			break
		}

		err = s.MyStorage.Increase(r.Context(), candidate.Candidate)
		if err != nil {
			log.Println(err.Error())
			break
		}

		res, err := s.GetAllValues(context.Background())

		WriteToAllConnections(res)
	}
}

func (s *Server) GetAllValues(ctx context.Context) (*GetAllResponse, error) {
	cat, err := s.MyStorage.Get(ctx, "cat")
	if err != nil {
		return nil, err
	}

	dog, err := s.MyStorage.Get(ctx, "dog")
	if err != nil {
		return nil, err
	}

	return &GetAllResponse{CatCount: cat, DogCount: dog}, nil
}

func WriteToAllConnections(object interface{}) {
	for key := range connections {
		client := key
		go func() {
			err := client.WriteJSON(object)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	}
}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value, err := s.MyStorage.Get(r.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(GetValueResponse{Count: value})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (s *Server) IncreaseValue(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	err := s.MyStorage.Increase(r.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetAllResponse struct {
	CatCount int `json:"cat_count"`
	DogCount int `json:"dog_count"`
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	cat, err := s.MyStorage.Get(r.Context(), "cat")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dog, err := s.MyStorage.Get(r.Context(), "dog")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(GetAllResponse{DogCount: dog, CatCount: cat})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
