package models

// Location represents the locations endpoint
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}