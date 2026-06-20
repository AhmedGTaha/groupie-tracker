package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"groupie-tracker/internal/models"
)

// BaseURL stores the real Groupie Tracker API address.
// Tests can use a fake URL instead, so they do not depend on the internet.
const BaseURL = "https://groupietrackers.herokuapp.com/api"

// Client groups the settings needed to talk to the API.
// Keeping these values in a struct makes the code easier to test and reuse.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a client configured for the real API.
func NewClient() *Client {
	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			// Timeout prevents the app from waiting forever if the API is slow.
			Timeout: 10 * time.Second,
		},
	}
}

// FetchArtists calls the /artists endpoint and converts the JSON response into Go structs.
func (c *Client) FetchArtists() ([]models.Artist, error) {
	// Get sends an HTTP GET request to the API.
	resp, err := c.HTTPClient.Get(c.BaseURL + "/artists")
	if err != nil {
		return nil, err
	}

	// defer runs after this function finishes.
	// Closing the body avoids leaking network resources.
	defer resp.Body.Close()

	// The API should return 200 OK for a successful request.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var artists []models.Artist

	// Decode reads JSON from the response body and fills the artists slice.
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
