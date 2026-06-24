package models

// Date represents the dates endpoint
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}