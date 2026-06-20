package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

type HomePageData struct {
	Title   string
	Artists []models.Artist
}

// ArtistPageData is the data shape sent to artist.html.
// Title is used for the browser tab, and Artist is the selected artist to display.
type ArtistPageData struct {
	Title  string
	Artist models.Artist
}

// artistFetcher describes only the API method these handlers need.
// This keeps handlers easy to test because tests can provide a fake version.
type artistFetcher interface {
	FetchArtists() ([]models.Artist, error)
}

// newAPIClient creates the real API client during normal app runs.
// Tests replace this variable so they do not call the real internet API.
var newAPIClient = func() artistFetcher {
	return api.NewClient()
}

// Template paths live in variables so tests can point handlers at temporary templates.
var homeTemplatePath = "templates/index.html"
var artistTemplatePath = "templates/artist.html"

// NewRouter wires URL paths to handler functions.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// "/" shows the home page.
	mux.HandleFunc("/", HomeHandler)

	// Support both detail URL styles:
	// /artist?id=1 and /artist/1
	mux.HandleFunc("/artist", ArtistHandler)
	mux.HandleFunc("/artist/", ArtistHandler)

	// This route is useful while learning/debugging the API client.
	mux.HandleFunc("/api-test", APITestHandler)

	return mux
}

// renderTemplate loads an HTML template file and writes the rendered page.
func renderTemplate(w http.ResponseWriter, filePath string, data any) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// The "/" route can catch other paths too, so guard it manually.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	client := newAPIClient()

	artists, err := client.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	data := HomePageData{
		Title:   "Groupie Tracker",
		Artists: artists,
	}

	renderTemplate(w, homeTemplatePath, data)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// First support links like /artist?id=1.
	idText := r.URL.Query().Get("id")

	// Also support links like /artist/1.
	if idText == "" && len(r.URL.Path) > len("/artist/") {
		idText = r.URL.Path[len("/artist/"):]
	}

	if idText == "" {
		http.Error(w, "Missing artist id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	client := newAPIClient()

	artists, err := client.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	var selectedArtist models.Artist
	found := false

	// Search through all artists until the requested ID is found.
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	data := ArtistPageData{
		Title:  selectedArtist.Name,
		Artist: selectedArtist,
	}

	renderTemplate(w, artistTemplatePath, data)
}

func APITestHandler(w http.ResponseWriter, r *http.Request) {
	client := newAPIClient()

	artists, err := client.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(artists) == 0 {
		fmt.Fprintln(w, "No artists found")
		return
	}

	fmt.Fprintf(w, "Fetched %d artists\n", len(artists))
	fmt.Fprintf(w, "First artist: %s\n", artists[0].Name)
}
