package server

import (
	"fmt"
	"go-aws/internal/controllers"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	addr = ":8080"
)

type Server struct {
	server          *http.Server
	router          *chi.Mux
	usersController *controllers.UsersController
}

func (s *Server) route() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)

	s.router.Get("/", s.handleHome)
	s.router.Get("/users/{id}", s.usersController.GetUser)
	s.router.Put("/users", s.usersController.PutUser)
}

func New(usersController *controllers.UsersController) *Server {
	router := chi.NewMux()
	return &Server{
		router:          router,
		server:          &http.Server{Addr: addr, Handler: router},
		usersController: usersController,
	}
}

func (a *Server) Run() {
	a.route()
	log.Printf("starting server on %s", addr)
	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("error starting server: %s", err)
	}
}

func (a *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `welcome to the app`)
}
