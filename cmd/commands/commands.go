package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Metadata is the metadata used to create a new service.
type Metadata struct {
	Name        string
	Owner       string
	Version     string
	Copyright   bool
	License     bool
	Description string
	Entrypoint  string
	GitIgnore   bool
}

func main() {
	app := &cli.App{
		Name:  "auth0-exercise-cli",
		Usage: "interact with the service provisioner",
		Authors: []*cli.Author{
			{
				Name:  "Creator",
				Email: "creator@auth0.com",
			},
		},
		Copyright:            "2020 COPYRIGHT",
		Version:              "0.0.1",
		EnableBashCompletion: true,
	}

	fileFlag := []cli.Flag{
		&cli.StringFlag{
			Name:  "file, f",
			Usage: "spec file used to manage options for the service",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a new service",
			Flags:   fileFlag,
			Action: func(c *cli.Context) error {
				f := c.String("file")
				if f == "" {
					generateServiceDefault()
				}

				generateServiceFromSpec(f)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateServiceDefault() {
	fmt.Println("works 1")
}

func generateServiceFromSpec(f string) {
	var metadata Metadata
	data, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &metadata)
	if err != nil {
		panic(err)
	}
}
