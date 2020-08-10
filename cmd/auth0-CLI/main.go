package main

import (
	"os"

	"github.com/juanri0s/auth0-exercise/internal/app"
)

func main() {
	err := app.StartCLI(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
