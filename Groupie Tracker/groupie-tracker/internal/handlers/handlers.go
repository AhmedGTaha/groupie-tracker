package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/service"
)

// appService is the behavior handlers need from the service layer.
// Tests can provide a fake implementation with these same methods.
type appService interface {
	AllArtists() ([]models.Artist, error)
	FindArtistByID(id int) (models.ArtistDetails, error)
	SearchArtists(query string) ([]models.Artist, error)
	Summary() (service.Summary, error)
}

// app uses the real API client during normal server runs.
var app appService = service.New(api.NewClient())

// Template paths live in variables so tests can point handlers at temporary files.
var homeTemplatePath = "templates/index.html"
var artistTemplatePath = "templates/artist.html"
var errorTemplatePath = "templates/error.html"

// NewRouter wires URL paths to handler functions.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Static files are served directly from the static folder.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/search", SearchHandler)

	// Support both detail URL styles:
	// /artist?id=1 and /artist/1
	mux.HandleFunc("/artist", ArtistHandler)
	mux.HandleFunc("/artist/", ArtistHandler)

	// This route is useful while learning/debugging the API client.
	mux.HandleFunc("/api-test", APITestHandler)

	return mux
}

// renderTemplate loads an HTML template file and writes the rendered page.
func renderTemplate(w http.ResponseWriter, filePath string, data any) error {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

func renderError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	data := models.ErrorPageData{
		Title:   http.StatusText(status),
		Status:  status,
		Message: message,
	}

	if err := renderTemplate(w, errorTemplatePath, data); err != nil {
		http.Error(w, message, status)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// The "/" route can catch other paths too, so guard it manually.
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "The page you requested does not exist.")
		return
	}

	artists, err := app.AllArtists()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to load artists. Please try again later.")
		return
	}

	data := models.HomePageData{
		Title:   "Groupie Tracker",
		Artists: artists,
	}

	if err := renderTemplate(w, homeTemplatePath, data); err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to render the home page.")
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	artists, err := app.SearchArtists(query)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to search artists. Please try again later.")
		return
	}

	data := models.HomePageData{
		Title:   "Search - Groupie Tracker",
		Query:   query,
		Artists: artists,
	}

	if err := renderTemplate(w, homeTemplatePath, data); err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to render search results.")
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := artistIDFromRequest(r)
	if !ok {
		renderError(w, http.StatusBadRequest, "Please choose an artist from the home page.")
		return
	}

	details, err := app.FindArtistByID(id)
	if errors.Is(err, service.ErrArtistNotFound) {
		renderError(w, http.StatusNotFound, "Artist not found.")
		return
	}
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to load artist details. Please try again later.")
		return
	}

	data := models.ArtistPageData{
		Title:   details.Artist.Name,
		Details: details,
	}

	if err := renderTemplate(w, artistTemplatePath, data); err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to render the artist page.")
	}
}

func APITestHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := app.Summary()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to fetch API summary.")
		return
	}

	fmt.Fprintf(w, "Artists: %d\n", summary.ArtistCount)
	fmt.Fprintf(w, "Locations: %d\n", summary.LocationCount)
	fmt.Fprintf(w, "Dates: %d\n", summary.DateCount)
	fmt.Fprintf(w, "Relations: %d\n", summary.RelationCount)

	if summary.FirstArtist != "" {
		fmt.Fprintf(w, "\nFirst artist: %s\n", summary.FirstArtist)
	}
	if summary.FirstLocation != 0 {
		fmt.Fprintf(w, "First location ID: %d\n", summary.FirstLocation)
	}
	if summary.FirstDate != 0 {
		fmt.Fprintf(w, "First date ID: %d\n", summary.FirstDate)
	}
	if summary.FirstRelation != 0 {
		fmt.Fprintf(w, "First relation ID: %d\n", summary.FirstRelation)
	}
}

func artistIDFromRequest(r *http.Request) (int, bool) {
	idText := r.URL.Query().Get("id")

	if idText == "" && len(r.URL.Path) > len("/artist/") {
		idText = strings.TrimPrefix(r.URL.Path, "/artist/")
	}

	if idText == "" {
		return 0, false
	}

	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
