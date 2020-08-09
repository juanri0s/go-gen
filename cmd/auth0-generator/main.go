package main

import (
	"log"

	"github.com/juanri0s/auth0-exercise/internal/app"
)

func main() {
	log.Println("auth0-generator has started")
	app.StartServer()
	app.StartCLI()
}
