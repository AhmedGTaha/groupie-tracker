package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchArtistsReturnsArtists(t *testing.T) {
	// httptest.Server starts a real local HTTP server just for this test.
	// That lets us test HTTP code without calling the real external API.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTPClient: &http.Client{
			Timeout: time.Second,
		},
	}

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

func TestFetchArtistsReturnsErrorForBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}

	_, err := client.FetchArtists()
	if err == nil {
		t.Fatal("expected an error for non-200 status, got nil")
	}
}

func TestFetchArtistsReturnsErrorForInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}

	_, err := client.FetchArtists()
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}
