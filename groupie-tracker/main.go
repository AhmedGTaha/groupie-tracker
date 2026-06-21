package main

import (
	"fmt"           // write text responses
	"log"           // print server messages
	"net/http"      // tools to build the server
	"html/template" // load and render HTML templates
)

// A handler is a function that runs when the browser visits a route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter is used to send a response back to the browser
	// r *http.Request contains information about the request, like the URL path

	// This checks if the user visited a wrong path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// This sends text back to the browser
	fmt.Fprint(w, "Server is running")
}

// It works the same way as homeHandler, but it is responsible for the /artists page
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter is used to send a response back to the browser
	// r *http.Request contains information about the request, like the URL path

	// This checks if the user visited a wrong path
	if r.URL.Path != "/artists" {
		http.NotFound(w, r)
		return
	}

	// This sends text back to the browser
	fmt.Fprint(w, "Artists page")
}

func main() {
	// mux = multiplexer
	// A mux is a router it decides which function should handle each route
	mux := http.NewServeMux()

	// This allows the server to serve files from the static folder
	// This registers a new route. Any URL that starts with "/static/"
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/artists", artistsHandler)

	log.Println("Server started at http://localhost:8080")

	// err is an error var used to store errors returned by http funcs
	// This starts the server on port 8080
	err := http.ListenAndServe(":8080", mux)

	// If the server cannot start, print error
	if err != nil {
		log.Fatal(err)
	}
}
