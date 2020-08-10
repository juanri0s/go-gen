package main

import (
	"os"

	"github.com/juanri0s/auth0-exercise/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("auth0-cli started")
	err := app.StartCLI(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
