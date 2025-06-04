package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	addr = ":8080"
)

type genID func(string) string
type Server struct {
	db          DynamoIface
	idGenerator genID
	server      *http.Server
	router      *chi.Mux
	userTable   string
}

func (s *Server) route() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)

	s.router.Get("/", s.handleHome)
	s.router.Get("/users/{id}", s.GetUser)
	s.router.Post("/users", s.PutUser)
}

func (a *Server) Run() {
	a.route()
	log.Printf("starting server on %s", addr)
	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("error starting server: %s", err)
	}
}

func (a *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	log.Println("invoked home")
	fmt.Fprintln(w, `welcome to the app`)
}
