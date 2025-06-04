package main

import (
	"context"
	"go-aws/internal/env"
	"go-aws/internal/server"
)

func main() {
	env.Load(".env")
	ctx := context.Background()

	_ = ctx

	app := server.Server{
		// db: usersRepo.InitUsersDB(ctx),

	}
	app.Run()
}
