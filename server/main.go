//go:build !js

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	fmt.Printf("Serving on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
