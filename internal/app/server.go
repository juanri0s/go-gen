package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const version = "1.0.0"
const timeout = 15 * time.Second

// HealthHandler checks for the healthiness of the api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`HEALTHY ` + version))
}

// RepoHandler manages the creation of a GH repository.
func RepoHandler(w http.ResponseWriter, r *http.Request) {
	var g Generator
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&g)
	if err != nil {
		log.Error(err.Error())
	}

	token := g.Token
	if token == "" {
		log.Error("invalid GH token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m := g.Metadata
	if m.ProjectPath == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid path`))
		return
	}

	switch r.Method {
	case "POST":
		err := initGit(m.ProjectPath)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = setupService(m)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		repoCfg := &github.Repository{
			Name:          github.String(m.Name),
			Private:       github.Bool(m.IsPrivate),
			AutoInit:      github.Bool(true),
			DefaultBranch: github.String(m.MainBranch),
			// TODO: add in any other important configs
		}

		log.WithFields(log.Fields{
			"repo-name": m.Name,
		}).Info("creating Github Repository")

		t := time.Now()
		repo, _, err := client.Repositories.Create(ctx, "", repoCfg)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Warn("hit rate limit")
			w.WriteHeader(r.Response.StatusCode)
			return
		}
		if _, ok := err.(*github.AcceptedError); ok {
			w.WriteHeader(r.Response.StatusCode)
			w.Write([]byte(`Github 202 Accepted`))
			return
		}
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = setRepoURL(m.ProjectPath, repo.GetGitURL())
		if err != nil {
			log.Error("Error setting repo url ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.WithFields(log.Fields{
			"repo-name":      m.Name,
			"repo-link":      repo.GetHTMLURL(),
			"execution-time": time.Since(t),
		}).Info("successfully created Github Repository")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`Only POST is supported`))
	}
}

// initGit initializes a project as a git project
func initGit(p string) error {
	if p == "" {
		return fmt.Errorf("invalid path")
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to initalize project as git repo - %w", err)
	}
	return nil
}

// initMod initializes a Go project with go modules.
func initMod(p string) error {
	if p == "" {
		return fmt.Errorf("invalid path")
	}

	cmd := exec.Command("go", "mod", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to initalize project with go mod - %w", err)
	}
	return nil
}

// setRepoURL sets a new remote address for the git project.
func setRepoURL(p string, url string) error {
	if p == "" {
		return fmt.Errorf("invalid path")
	}

	if url == "" {
		return fmt.Errorf("invalid url")
	}

	cmd := exec.Command("git", "remote", "add", "origin", url)
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to set %s on repo  - %w", url, err)
	}
	return nil
}

// StartServer initializes the server and starts it with a graceful shutdown on signals.
func StartServer() {
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/repository", RepoHandler)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.WithFields(log.Fields{
			"signal": sig,
		}).Info("received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error(err.Error())
		}
	}()

	log.Info("starting server on port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err.Error())
		return
	}
}
