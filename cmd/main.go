package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

const version = "1.0.0"
const timeout = 120 * time.Second

func main() {
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
