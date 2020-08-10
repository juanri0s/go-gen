package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
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

// new returns a new service with default metadata values.
func (m *Metadata) new() Metadata {
	return Metadata{
		ProjectPath: "",
		Name:        "default-repo",
		Owner:       "default-owner",
		Version:     "1.0.0",
		Copyright:   true,
		License:     true,
		Description: "A default service for Auth0",
		Entrypoint:  "/app",
		GitIgnore:   true,
	}
}

func generateServiceFromDefault() (string, error) {
	var m Metadata
	m = m.new()
	repo, err := generate(m)
	if err != nil {
		return "", err
	}

	return repo, nil
}

func generateServiceFromFile(f string) (string, error) {
	var m Metadata
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}

	if strings.Contains(f, "yaml") || strings.Contains(f, "yml") {
		err = yaml.Unmarshal(data, &m)
		if err != nil {
			return "", err
		}
	} else if strings.Contains(f, "json") {
		err = json.Unmarshal(data, &m)
		if err != nil {
			return "", err
		}
	}

	// Instead of asking the user for a path, we would like them to run the command in the WD they want
	m.ProjectPath, err = os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := generate(m)
	return repo, nil
}

func generate(m Metadata) (string, error) {
	log.WithFields(log.Fields{
		"repo-name": m.Name,
	}).Info("Generating Go service")
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	const url = "http://localhost:8080/repository"
	resp, err := c.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("HTTP request failed with HTTP Status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var repo *github.Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		return "", err
	}

	return repo.GetHTMLURL(), nil
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
				t := time.Now()
				f := c.String("file")
				var repo string
				var err error
				if f == "" {
					repo, err = generateServiceFromDefault()
					if err != nil {
						return fmt.Errorf("%w", err)
					}
				} else {
					repo, err = generateServiceFromFile(f)
					if err != nil {
						return fmt.Errorf("%w", err)
					}
				}

				log.WithFields(log.Fields{
					"repo-link":      repo,
					"execution-time": time.Since(t),
				}).Info("Service created successfully")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
