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
const timeout = 15 * time.Second

// HealthHandler TODO:
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`HEALTHY ` + version))
}

// RepoHandler TODO
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
		err := initRepo(m.ProjectPath)
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

func initRepo(p string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func initMod(p string) error {
	cmd := exec.Command("go", "mod", "init")
	cmd.Dir = p
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func setRepoURL(p string, url string) error {
	cmd := exec.Command("git", "remote", "add", "origin", url)
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
