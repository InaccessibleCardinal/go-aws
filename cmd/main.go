package main

import (
	"context"
	"go-aws/internal/controllers"
	"go-aws/internal/env"
	usersRepo "go-aws/internal/repos/users"
	"go-aws/internal/server"
)

func main() {
	env.Load(".env")
	ctx := context.Background()
	usersController := controllers.NewUsersController(usersRepo.InitUsersDB(ctx))

	app := server.New(usersController)
	app.Run()
}
