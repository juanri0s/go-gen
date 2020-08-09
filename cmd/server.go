package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// HealthHandler TODO:
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`HEALTHY ` + version))
}

// RepoHandler TODO
func RepoHandler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`Invalid Github Auth Token`))
		return
	}

	// TODO: read in a path
	path := "test"
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Invalid Path`))
		return
	}

	switch r.Method {
	case "POST":
		createProject(path)
		initRepo(path)

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		repoCfg := &github.Repository{
			Name:    github.String("test-repo1"),
			Private: github.Bool(true),
			// TODO: add in any other important configs
		}
		repo, _, err := client.Repositories.Create(ctx, "", repoCfg)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		setRepoURL(path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Successfully created new repo:` + repo.GetURL()))
	default:
		fmt.Fprintf(w, "Only POST is supported")
	}
}

func createProject(p string) {
	// Check if the directory exists first
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.Mkdir(p, 0755)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func initRepo(p string) {
	// initalize the repo as a git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func setRepoURL(p string) {
	// TODO: take the github repo url and set it as the origin
	cmd := exec.Command("git", "remote", "add", "origin")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
}
