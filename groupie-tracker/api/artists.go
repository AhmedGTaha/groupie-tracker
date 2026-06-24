package api

import (
	"encoding/json" // decode JSON into Go structs
	"fmt"           // create custom error messages
	"io"            // read response body from the API
	"net/http"      // send HTTP requests
	"time"          // set timeout for API requests

	"groupie-tracker/models"
)

// API endpoint for artists data
const artistsAPIURL = "https://groupietrackers.herokuapp.com/api/artists"

// FetchArtists gets artist data from the API and returns it as Go structs
func FetchArtists() ([]models.Artist, error) {
	// Fetch raw JSON bytes from the API
	body, err := fetchAPIData(artistsAPIURL)
	if err != nil {
		return nil, err
	}

	// Decode raw JSON bytes into Artist structs
	artists, err := decodeArtists(body)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// GetArtistByID finds and returns one artist using its ID
func GetArtistByID(id int) (models.Artist, error) {

	// Fetch all artists from the API
	artists, err := FetchArtists()

	// If FetchArtists failed, return the error immediately
	if err != nil {
		return models.Artist{}, err
	}

	// Loop through every artist in the slice
	for _, artist := range artists {

		// Check if this artist's ID matches the requested ID
		if artist.ID == id {

			// Artist found, return it and no error
			return artist, nil
		}
	}

	// If we reach this point, no artist matched the ID
	return models.Artist{}, fmt.Errorf("artist not found")
}

// fetchAPIData sends a request to an API URL and returns the response body
func fetchAPIData(url string) ([]byte, error) {
	// Create a client with timeout so the server does not hang forever
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Send GET request to the API
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// Close the response body when the function finishes
	defer resp.Body.Close()

	// Check if the API returned a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read all data from the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// decodeArtists converts raw JSON bytes into a slice of Artist structs
func decodeArtists(data []byte) ([]models.Artist, error) {
	// Create a variable to store all decoded artists
	var artists []models.Artist

	// Decode the JSON data into the artists slice
	err := json.Unmarshal(data, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}