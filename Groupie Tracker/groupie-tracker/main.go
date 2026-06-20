package main

import (
	"log"
	"net/http"

	"groupie-tracker/internal/handlers"
)

func main() {
	log.Println("Server running at http://localhost:3000")

	err := http.ListenAndServe(":3000", handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}