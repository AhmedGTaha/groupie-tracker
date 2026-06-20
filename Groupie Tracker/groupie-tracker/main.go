package main

import (
	"log"
	"net/http"
	"os"

	"groupie-tracker/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server running at http://localhost:%s\n", port)

	err := http.ListenAndServe(addr, handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}
