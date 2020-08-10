package app

import (
	"context"
	"encoding/json"
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
const timeout = 120 * time.Second

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
		w.Write([]byte(`invalid Github Token`))
		return
	}

	// TODO: read in a path
	path := "test"
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`invalid path`))
		return
	}

	switch r.Method {
	case "POST":
		// var m
		// decoder := json.NewDecoder(r.Body)
		// err := decoder.Decode(&m)
		// if err != nil {
		// 	panic(err)
		// }

		err := initRepo(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		repoCfg := &github.Repository{
			Name:    github.String("test-repo"),
			Private: github.Bool(true),
			// TODO: add in any other important configs
		}

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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = setRepoURL(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(repo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`Only POST is supported`))
	}
}

func initRepo(p string) error {
	// initalize the repo as a git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func setRepoURL(p string) error {
	// TODO: take the github repo url and set it as the origin
	cmd := exec.Command("git", "remote", "add", "origin")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

// Flow: 1. Create Dir 2. Apply the Template from the spec 3. Create the repo in GH 4. Set the repo url returned to the project

// StartServer TODO:
func StartServer() {
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	// testTemplate()
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
