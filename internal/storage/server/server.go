package server

import "github.com/modaniru/bio-vue/cat-vs-dogs/internal/storage"

type Server struct{
	MyStorage storage.Storage
}

func NewServer(storage storage.Storage) *Server{
	return &Server{MyStorage: storage}
}