package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient starts a fake HTTP API server and returns a Client pointed at it.
// This lets tests exercise real HTTP requests without depending on the public API.
func newTestClient(t *testing.T, handler http.Handler) (*Client, func()) {
	t.Helper()

	server := httptest.NewServer(handler)
	client := &Client{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}

	return client, server.Close
}

func TestFetchArtistsReturnsArtists(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The client should call the /artists endpoint.
		if r.URL.Path != "/artists" {
			t.Fatalf("expected path /artists, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"image": "queen.jpg",
				"name": "Queen",
				"members": ["Freddie Mercury", "Brian May"],
				"creationDate": 1970,
				"firstAlbum": "14-12-1973"
			}
		]`))
	}))
	defer cleanup()

	artists, err := client.FetchArtists()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(artists) != 1 {
		t.Fatalf("expected 1 artist, got %d", len(artists))
	}

	if artists[0].Name != "Queen" {
		t.Fatalf("expected artist name Queen, got %q", artists[0].Name)
	}
}

func TestFetchLocationsReturnsLocations(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/locations" {
			t.Fatalf("expected path /locations, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"index": [
				{
					"id": 1,
					"locations": ["london-uk", "paris-france"]
				}
			]
		}`))
	}))
	defer cleanup()

	locations, err := client.FetchLocations()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(locations.Index) != 1 {
		t.Fatalf("expected 1 location record, got %d", len(locations.Index))
	}

	if locations.Index[0].Locations[0] != "london-uk" {
		t.Fatalf("expected first location london-uk, got %q", locations.Index[0].Locations[0])
	}
}

func TestFetchDatesReturnsDates(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dates" {
			t.Fatalf("expected path /dates, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"index": [
				{
					"id": 1,
					"dates": ["01-01-2026"]
				}
			]
		}`))
	}))
	defer cleanup()

	dates, err := client.FetchDates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(dates.Index) != 1 {
		t.Fatalf("expected 1 date record, got %d", len(dates.Index))
	}

	if dates.Index[0].Dates[0] != "01-01-2026" {
		t.Fatalf("expected first date 01-01-2026, got %q", dates.Index[0].Dates[0])
	}
}

func TestFetchRelationsReturnsRelations(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/relation" {
			t.Fatalf("expected path /relation, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"index": [
				{
					"id": 1,
					"datesLocations": {
						"london-uk": ["01-01-2026"]
					}
				}
			]
		}`))
	}))
	defer cleanup()

	relations, err := client.FetchRelations()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(relations.Index) != 1 {
		t.Fatalf("expected 1 relation record, got %d", len(relations.Index))
	}

	if relations.Index[0].DatesLocations["london-uk"][0] != "01-01-2026" {
		t.Fatalf("expected relation date 01-01-2026, got %q", relations.Index[0].DatesLocations["london-uk"][0])
	}
}

func TestFetchArtistsReturnsErrorForBadStatus(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Any non-200 status should become an error.
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer cleanup()

	_, err := client.FetchArtists()
	if err == nil {
		t.Fatal("expected an error for non-200 status, got nil")
	}
}

func TestFetchArtistsReturnsErrorForInvalidJSON(t *testing.T) {
	client, cleanup := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json`))
	}))
	defer cleanup()

	_, err := client.FetchArtists()
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}
