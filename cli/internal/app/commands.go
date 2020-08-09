package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Metadata is the metadata used to create a new service.
type Metadata struct {
	ProjectPath string
	Name        string
	Owner       string
	Version     string
	Copyright   bool
	License     bool
	Description string
	Entrypoint  string
	GitIgnore   bool
}

func generateServiceDefault() {
	fmt.Println("works 1")
}

func generateServiceFromFile(f string) {
	var metadata Metadata
	data, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	if strings.Contains(f, "yaml") || strings.Contains(f, "yml") {
		err = yaml.Unmarshal(data, &metadata)
		if err != nil {
			panic(err)
		}
	} else if strings.Contains(f, "json") {
		err = json.Unmarshal(data, &metadata)
		if err != nil {
			panic(err)
		}
	}

	// Instead of asking the user for a path, we would like them to run the command in the WD they want
	metadata.ProjectPath, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	generate(metadata)
}

func generate(m Metadata) {
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	const url = "http://localhost:8080/respository"
	resp, err := c.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var repo *github.Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(repo.GetURL())
}

// StartCLI TODO
func StartCLI() {
	app := &cli.App{
		Name:  "auth0-exercise-cli",
		Usage: "interact with the service provisioner",
		Authors: []*cli.Author{
			{
				Name:  "Creator",
				Email: "creator@auth0.com",
			},
		},
		Copyright: `
			Copyright 2020 Juan Rios. All rights reserved.
			Use of this source code is governed by an MIT License
			license that can be found in the LICENSE file.
		`,
		Version:              "0.0.1",
		EnableBashCompletion: true,
	}

	fileFlag := []cli.Flag{
		&cli.StringFlag{
			Name:  "file, f",
			Usage: "file used to manage custom metadata for the service",
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

				generateServiceFromFile(f)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
