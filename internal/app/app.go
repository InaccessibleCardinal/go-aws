package app

import (
	"context"
	"go-aws/internal/controllers"
	"go-aws/internal/env"
	usersRepo "go-aws/internal/repos/users"
	"go-aws/internal/routers/types"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	usersController *controllers.UsersController
	router          types.Router
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) routes() {
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middleware.URLFormat)

	a.router.Get("/users/{id}", a.usersController.GetUser)
	a.router.Put("/users", a.usersController.PutUser)
}

func New(router types.Router, usersController *controllers.UsersController) *App {
	return &App{router: router, usersController: usersController}
}

func Run() error {
	env.Load(".env")
	ctx := context.Background()
	usersController := controllers.NewUsersController(usersRepo.InitUsersDB(ctx))
	app := New(chi.NewRouter(), usersController)

	app.routes()
	if err := http.ListenAndServe(":8888", app); err != nil {
		return err
	}
	return nil
}
