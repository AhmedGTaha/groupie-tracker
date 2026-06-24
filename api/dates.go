package api

import (
	"encoding/json"

	"groupie-tracker/models"
)

// API endpoint for dates data
const datesAPIURL = "https://groupietrackers.herokuapp.com/api/dates"

// FetchDates gets dates data from the API and returns it as Go structs
func FetchDates() ([]models.Date, error) {

	// Fetch raw JSON bytes from the API
	body, err := fetchAPIData(datesAPIURL)
	if err != nil {
		return nil, err
	}

	// Decode raw JSON bytes into Date structs
	dates, err := decodeDates(body)
	if err != nil {
		return nil, err
	}

	return dates, nil
}

// decodeDates converts raw JSON bytes into a slice of Date structs
func decodeDates(data []byte) ([]models.Date, error) {

	// The API wraps dates inside an "index" field
	var response struct {
		Index []models.Date `json:"index"`
	}

	// Decode JSON into the response struct
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Index, nil
}