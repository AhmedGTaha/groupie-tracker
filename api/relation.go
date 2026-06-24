package api

import (
	"encoding/json"

	"groupie-tracker/models"
)

// API endpoint for relation data
const relationAPIURL = "https://groupietrackers.herokuapp.com/api/relation"

// FetchRelations gets relation data from the API and returns it as Go structs
func FetchRelations() ([]models.Relation, error) {

	// Fetch raw JSON bytes from the API
	body, err := fetchAPIData(relationAPIURL)
	if err != nil {
		return nil, err
	}

	// Decode raw JSON bytes into Relation structs
	relations, err := decodeRelations(body)
	if err != nil {
		return nil, err
	}

	return relations, nil
}

// decodeRelations converts raw JSON bytes into a slice of Relation structs
func decodeRelations(data []byte) ([]models.Relation, error) {

	// The API wraps relations inside an "index" field
	var response struct {
		Index []models.Relation `json:"index"`
	}

	// Decode JSON into the response struct
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response.Index, nil
}