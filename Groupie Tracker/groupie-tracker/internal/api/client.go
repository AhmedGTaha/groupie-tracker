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

// getJSON sends a GET request to an API endpoint and decodes the JSON response.
func (c *Client) getJSON(endpoint string, target any) error {
	// endpoint is a path like "/artists"; BaseURL is the shared API root.
	resp, err := c.HTTPClient.Get(c.BaseURL + endpoint)
	if err != nil {
		return err
	}

	// Always close the response body when we are done reading it.
	defer resp.Body.Close()

	// Non-200 responses usually mean the API rejected the request or had a problem.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	// target is passed as any so this helper can decode artists, locations, dates, or relations.
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

// FetchArtists calls the /artists endpoint and converts the JSON response into Go structs.
func (c *Client) FetchArtists() ([]models.Artist, error) {
	var artists []models.Artist

	err := c.getJSON("/artists", &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// FetchLocations calls the /locations endpoint.
func (c *Client) FetchLocations() (models.LocationIndex, error) {
	var locations models.LocationIndex

	err := c.getJSON("/locations", &locations)
	if err != nil {
		return models.LocationIndex{}, err
	}

	return locations, nil
}

// FetchDates calls the /dates endpoint.
func (c *Client) FetchDates() (models.DateIndex, error) {
	var dates models.DateIndex

	err := c.getJSON("/dates", &dates)
	if err != nil {
		return models.DateIndex{}, err
	}

	return dates, nil
}

// FetchRelations calls the /relation endpoint.
func (c *Client) FetchRelations() (models.RelationIndex, error) {
	var relations models.RelationIndex

	err := c.getJSON("/relation", &relations)
	if err != nil {
		return models.RelationIndex{}, err
	}

	return relations, nil
}
