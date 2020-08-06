package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
	switch r.Method {
	case "POST":
		token := os.Getenv("GITHUB_AUTH_TOKEN")
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Successfully created new repo:` + repo.GetName()))
	default:
		fmt.Fprintf(w, "Only POST is supported")
	}
}
