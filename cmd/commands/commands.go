package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "auth0-exercise-cli",
		Usage:                "interact with the service provisioner",
		Version:              "0.0.1",
		EnableBashCompletion: true,
	}

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a new service",
			Action: func(c *cli.Context) error {
				generateServiceDefault()
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "applys specs to the service being created",
					Action: func(c *cli.Context) error {
						generateServiceFromSpec()
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateServiceFromSpec() {
	fmt.Println("works 1")
}

func generateServiceDefault() {
	fmt.Println("works 2")
}
