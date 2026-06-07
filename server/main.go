//go:build !js

// Serveur HTTP statique pour l'interface web de la calculatrice RPN.
// Sert les fichiers du répertoire web/ sur le port spécifié par la variable
// d'environnement PORT (8080 par défaut).
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
