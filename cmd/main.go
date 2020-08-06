package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const version = "1.0.0"
const timeout = 120 * time.Second

func main() {
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/repository", repoHandler)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		fmt.Println("received shutdown signal:", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("http server shutdown: %s", err.Error())
		}
	}()

	fmt.Printf("Starting server at port 8080\n")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("http server:%s", err.Error())
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`HEALTHY ` + version))
}

func repoHandler(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println("%s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		setRepoURL(path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Successfully created new repo:` + repo.GetName()))
	default:
		fmt.Fprintf(w, "Only POST is supported")
	}
}

func createProject(p string) {
	// Check if the directory exists first
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.Mkdir(p, 0755)
		if err != nil {
			fmt.Println("%s", err.Error())
		}
	}
}

func initRepo(p string) {
	// initalize the repo as a git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = p
	// TODO: return some kind of progress to the user
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("%s", err.Error())
	}
}

func setRepoURL(p string) {
	// TODO: take the github repo url and set it as the origin
	cmd := exec.Command("git", "remote", "add", "origin", "")
	cmd.Dir = p
	// TODO: return some kind of progress to the user
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("%s", err.Error())
	}
}

// Flow: 1. Create Dir 2. Apply the Template from the spec 3. Create the repo in GH 4. Set the repo url returned to the project
