package server

import (
	"log"
	"net/http"
	"scraper/config"
	"scraper/handlers"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	handler *handlers.Handler
}

func NewServer() *Server {
	server := &Server{
		router: chi.NewRouter(),
	}

	return server
}

func (s *Server) ServerHTTP() {
	s.Routes()

	conf := config.GetConfig()

	log.Printf("✨ HTTP Server running on http://localhost:%s/", conf.ServerPort)

	log.Fatal(http.ListenAndServe("localhost:"+conf.ServerPort, s.router))
}

func (s *Server) SetHandler(handler *handlers.Handler) {
	s.handler = handler
}
