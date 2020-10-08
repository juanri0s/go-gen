package main

import (
	"os"

	"github.com/juanri0s/go-gen/internal/app"
)

func main() {
	err := app.StartCLI(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
