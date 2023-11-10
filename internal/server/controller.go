package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

func (s *Server) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	// разобраться с cors
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	sub := chi.NewRouter()
	t := sub.With(s.Test)
	{
		t.Get("/", s.GetAll)
	}
	r := t.With(s.isAvailable)
	r = r.With(s.DataChanged)
	{
		r.Get("/{key}", s.GetValue)
		r.Put("/{key}", s.IncreaseValue)
	}
	router.Mount("/v1", sub)
	router.HandleFunc("/w", s.Websocket)
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

	res, err := s.GetAllValues(context.Background())
	if err != nil {
		slog.Error(fmt.Errorf("disconect, error: %w", err).Error())
		return
	}

	err = conn.WriteJSON(res)
	if err != nil {
		slog.Error(fmt.Errorf("disconect, error: %w", err).Error())
		return
	}

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				slog.Debug("close connection...")
				break
			}
			slog.Error(fmt.Errorf("disconect, error: %w", err).Error())
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		if err != nil {
			slog.Error(fmt.Errorf("disconect, error: %w", err).Error())
			break
		}
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
	slog.Debug(fmt.Sprintf("notify all clients. count: %d", len(connections)))
	for key := range connections {
		client := key
		go func() {
			err := client.WriteJSON(object)
			if err != nil {
				slog.Error(err.Error())
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
