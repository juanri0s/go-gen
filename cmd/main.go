package main

import (
	"fmt"
	"net/http"
)

var version = "1.0.0"

func main() {
	// TODO: sample, remove
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.ListenAndServe(":8080", nil)
	// TODO: health check and graceful shutdown
}
