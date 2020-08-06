package main

import (
	"fmt"
	"net/http"
	"time"
)

const version = "1.0.0"
const timeout = 120 * time.Second

func main() {
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	// TODO: graceful shutdown
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/repository", repoHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("http server:%s", err.Error())
		return
	}

	// TODO: listen for signals
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`healthy`))
}

func repoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`created`))
	default:
		fmt.Fprintf(w, "Only GET and POST methods are supported")
	}
}
