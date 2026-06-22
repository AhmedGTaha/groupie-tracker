package main

import (
	"html/template" // load and render HTML templates
	"log"           // print server messages
	"net/http"      // tools to build the server
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

	// Go to template file -> Pars home.html as template
	// It returns the template and any errors
	tmpl, err := template.ParseFiles("templates/home.html")

	// If home.html is missing or has a problem, the server should not crash
	if err != nil {
		log.Println("template parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// This sends the HTML page to the browser and nil = not passing data into the HTML
	err = tmpl.Execute(w, nil)

	if err != nil {
		// If Go fails while sending the page, show a server error instead of crashing
		log.Println("template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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

	// Go to templates folder -> parse artists.html as a template
	// It returns the template and any errors
	tmpl, err := template.ParseFiles("templates/artists.html")

	// If artists.html is missing or has a problem, the server should not crash
	if err != nil {
		log.Println("template parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// This sends the HTML page to the browser and nil = not passing data into the HTML
	err = tmpl.Execute(w, nil)

	if err != nil {
		// If Go fails while sending the page, show a server error instead of crashing
		log.Println("template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
