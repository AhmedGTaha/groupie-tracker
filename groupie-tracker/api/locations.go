package api

import (
	"encoding/json"

	"groupie-tracker/models"
)

// API endpoint for locations data
const locationsAPIURL = "https://groupietrackers.herokuapp.com/api/locations"

// FetchLocations gets location data from the API and returns it as Go structs
func FetchLocations() ([]models.Location, error) {

	// Fetch raw JSON bytes from the API
	body, err := fetchAPIData(locationsAPIURL)
	if err != nil {
		return nil, err
	}

	// Decode raw JSON bytes into Location structs
	locations, err := decodeLocations(body)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func decodeLocations(data []byte) ([]models.Location, error) {
	var response struct {
		index []models.Location
	}

	err :=json.Unmarshal(data, &response)
		if err != nil {
		return nil, err
	}

	return response.index, nil
}