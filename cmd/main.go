package main

import (
	"go-aws/internal/app"
	"go-aws/internal/env"
	"log"
)

func main() {
	env.Load(".env")
	if err := app.Run(); err != nil {
		log.Fatal("error starting app: ", err)
	}

}
