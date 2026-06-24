// Test repo

package main

import (
	"fmt"           // format text messages
	"html/template" // load and render HTML templates
	"log"           // print server messages
	"net/http"      // tools to build the server
	"strconv"

	"groupie-tracker/api"    // functions for fetching API data
	"groupie-tracker/models" // data structures used in the app
)

// ArtistsPageData is the data we send to artists.html
type ArtistsPageData struct {
	StatusMessage string
	Artists       []models.Artist
}

// ArtistDetailsPageData is the data we send to artist-details.html
type ArtistDetailsPageData struct {
	Artist models.Artist
}

// A handler is a function that runs when the browser visits a route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter is used to send a response back to the browser
	// r *http.Request contains information about the request, like the URL path

	// This checks if the user visited a wrong path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Go to templates folder -> parse home.html as a template
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

	// Creates an empty slice of Artists
	artists := []models.Artist{}

	// This message will be sent to the HTML template
	statusMessage := ""

	fetchedArtists, err := api.FetchArtists()

	if err != nil {
		// If the API request fails, the page should still work
		log.Println("API fetching error:", err)
		statusMessage = "Could not load artist data right now."
	} else {
		// stores the fetched artists in the artists slice
		artists = fetchedArtists
		// message
		statusMessage = fmt.Sprintf("Loaded %d artists from the API.", len(artists))
	}

	// One object to be displayed in the artists page
	pageData := ArtistsPageData{
		StatusMessage: statusMessage,
		Artists:       artists,
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

	// This sends the HTML page to the browser
	// pageData is passed into the HTML template
	err = tmpl.Execute(w, pageData)

	if err != nil {
		// If Go fails while sending the page, show a server error instead of crashing
		log.Println("template execution error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// artistDetailsHandler is responsible for the /artist details page
func artistDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// This checks if the user visited a wrong path
	if r.URL.Path != "/artist" {
		http.NotFound(w, r)
		return
	}

	// Read the artist ID from the URL
	// Example: /artist?id=1
	artistID := r.URL.Query().Get("id")

	// If the id is missing, return a bad request error
	if artistID == "" {
		http.Error(w, "Missing artist id", http.StatusBadRequest)
		return
	}

	// Convert the URL id from string to int
	id, err := strconv.Atoi(artistID)

	// If conversion fails, the id is invalid
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	// Fetch the artist with the matching ID
	selectedArtist, err := api.GetArtistByID(id)

	// If the artist does not exist, return 404
	if err != nil {
		log.Println("artist lookup error:", err)
		http.NotFound(w, r)
		return
	}

	// Go to templates folder -> parse artist-details.html as a template
	// It returns the template and any errors
	tmpl, err := template.ParseFiles("templates/artist-details.html")

	// If artist-details.html is missing or has a problem, the server should not crash
	if err != nil {
		log.Println("template parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the data that will be sent to the template
	pageData := ArtistDetailsPageData{
		Artist: selectedArtist,
	}

	// This sends the HTML page to the browser
	err = tmpl.Execute(w, pageData)

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

	mux.HandleFunc("/artist", artistDetailsHandler)
	mux.HandleFunc("/artists", artistsHandler)
	mux.HandleFunc("/", homeHandler)

	log.Println("Server started at http://localhost:8080")

	// err is an error var used to store errors returned by http funcs
	// This starts the server on port 8080
	err := http.ListenAndServe(":8080", mux)

	// If the server cannot start, print error
	if err != nil {
		log.Fatal(err)
	}
}
