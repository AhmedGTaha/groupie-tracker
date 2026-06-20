package main

import (
	"log"
	"net/http"

	"groupie-tracker/internal/handlers"
)

func main() {
	log.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}
