package app

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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

// Start TODO:
func Start() {
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	testTemplate()
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/repository", RepoHandler)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		fmt.Println("received shutdown signal:", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println(err.Error())
		}
	}()

	fmt.Printf("Starting server at port 8080\n")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err.Error())
		return
	}
}

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
		// var m
		// decoder := json.NewDecoder(r.Body)
		// err := decoder.Decode(&m)
		// if err != nil {
		// 	panic(err)
		// }

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

		// We don't care about the resp pagination, we might care about the rate for the future though
		repo, _, err := client.Repositories.Create(ctx, "", repoCfg)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		setRepoURL(path)

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(repo)
	default:
		fmt.Fprintf(w, "Only POST is supported")
	}
}

// TODO: deletion candidate
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

// Example TODO:
type Example struct {
	ServicePackage string
	Copyright      bool
	DefaultImports string
}

func testTemplate() {
	var one = Example{
		"main", true, `"fmt", "os"`,
	}

	templateText, err := ioutil.ReadFile("./generator/templates/main.gotmpl")
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("Main").Parse(string(templateText))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	f, err := os.Create("./generator/templates/testing.go")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	// Run the template to verify the output.
	err = tmpl.Execute(f, one)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	f.Close()
}

// Flow: 1. Create Dir 2. Apply the Template from the spec 3. Create the repo in GH 4. Set the repo url returned to the project
