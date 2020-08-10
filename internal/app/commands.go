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

// Generator is the service generator using the GH token for repo creation and the metadata for the configured service.
type Generator struct {
	Token    string
	Metadata Metadata
}

// Metadata is the metadata used to create a new service.
type Metadata struct {
	ProjectPath  string
	Name         string
	Owner        string
	Version      string
	HasCopyright bool
	HasLicense   bool
	Imports      string
	Description  string
	Entrypoint   string
	HasGitIgnore bool
	MainBranch   string
	IsPrivate    bool
}

// new returns a new service with default metadata values.
func (m *Metadata) new() Metadata {
	return Metadata{
		ProjectPath:  "C:\\Users\\Juan\\go\\src\\github.com\\juanri0s\\test\\",
		Name:         "default-repo",
		Owner:        "default-owner",
		Version:      "1.0.0",
		Imports:      DefaultImports,
		Description:  "A default service for auth0",
		Entrypoint:   "default-service",
		MainBranch:   "main",
		IsPrivate:    true,
		HasCopyright: true,
		HasLicense:   true,
		HasGitIgnore: true,
	}
}

// generateServiceFromDefault creates a new service based on the default configurations.
func generateServiceFromDefault(token string) (string, error) {
	var g Generator
	var m Metadata
	m = m.new()
	g.Metadata = m
	g.Token = token
	repo, err := g.generate()
	if err != nil {
		return "", err
	}

	return repo, nil
}

// getWorkingDirectory returns the current working directory.
func getWorkingDirectory() string {
	p, err := os.Getwd()
	if err != nil {
		log.Error("Error getting working directory")
		return ""
	}
	return p
}

// generateServiceFromFile creates a new service based on the given file configurations.
func generateServiceFromFile(f string, token string) (string, error) {
	var g Generator
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
	m.ProjectPath = getWorkingDirectory()

	g.Metadata = m
	g.Token = token

	repo, err := g.generate()
	return repo, nil
}

// generate takes a generator and generates the service through the api.
func (g *Generator) generate() (string, error) {
	log.WithFields(log.Fields{
		"repo-name": g.Metadata.Name,
	}).Info("generating Auth0 service")
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := json.Marshal(g)
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

// StartCLI initializes the CLI start and exit.
func StartCLI(args []string) error {
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

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "file, f",
			Usage: "file used to manage custom metadata for the service (optional)",
		},
		&cli.StringFlag{
			Name:     "token",
			Usage:    "GH token to interact with GH API (required)",
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a new Auth0 service",
			Flags:   flags,
			Action: func(c *cli.Context) error {
				t := time.Now()
				f := c.String("file")
				token := c.String("token")
				var repo string
				var err error
				if f == "" {
					repo, err = generateServiceFromDefault(token)
					if err != nil {
						return fmt.Errorf("%w", err)
					}
				} else {
					repo, err = generateServiceFromFile(f, token)
					if err != nil {
						return fmt.Errorf("%w", err)
					}
				}

				log.WithFields(log.Fields{
					"repo-link":      repo,
					"execution-time": time.Since(t),
				}).Info("successfully created Auth0 service")
				return nil
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
