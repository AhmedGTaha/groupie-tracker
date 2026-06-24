package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFetchAPIData_Success(t *testing.T) {
    // Create a fake server that returns a 200 with a body
    expectedBody := `[{"id":1,"name":"Test Artist"}]`
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(expectedBody))
    }))
    defer server.Close()

    // Call fetchAPIData with the fake server's URL
    body, err := fetchAPIData(server.URL)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if string(body) != expectedBody {
        t.Errorf("got body %q, want %q", string(body), expectedBody)
    }
}

func TestFetchAPIData_Non200(t *testing.T) {
    // Create a fake server that returns a 500 error
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    _, err := fetchAPIData(server.URL)
    if err == nil {
        t.Fatal("expected an error for non-200 response, got nil")
    }
}

func TestFetchAPIData_InvalidURL(t *testing.T) {
    // Pass a clearly invalid URL (no server)
    _, err := fetchAPIData("http://127.0.0.1:0/nothing")
    if err == nil {
        t.Fatal("expected an error for invalid URL, got nil")
    }
}

func TestFetchArtists_Success(t *testing.T) {
    // Fake JSON like the real API returns
    jsonResponse := `[
        {
            "id": 1,
            "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
            "name": "Queen",
            "members": ["Freddie Mercury", "Brian May"],
            "creationDate": 1970,
            "firstAlbum": "1973"
        },
        {
            "id": 2,
            "image": "https://groupietrackers.herokuapp.com/api/images/scorpions.jpeg",
            "name": "Scorpions",
            "members": ["Klaus Meine", "Rudolf Schenker"],
            "creationDate": 1965,
            "firstAlbum": "1972"
        }
    ]`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(jsonResponse))
    }))
    defer server.Close()

    // Override the API URL temporarily (we'll refactor later)
    originalURL := artistsAPIURL
    artistsAPIURL = server.URL
    defer func() { artistsAPIURL = originalURL }()

    artists, err := FetchArtists()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(artists) != 2 {
        t.Errorf("expected 2 artists, got %d", len(artists))
    }
    if artists[0].Name != "Queen" {
        t.Errorf("expected first artist name Queen, got %s", artists[0].Name)
    }
    if len(artists[1].Members) != 2 {
        t.Errorf("expected 2 members for Scorpions, got %d", len(artists[1].Members))
    }
}

func TestFetchArtists_Non200(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    originalURL := artistsAPIURL
    artistsAPIURL = server.URL
    defer func() { artistsAPIURL = originalURL }()

    _, err := FetchArtists()
    if err == nil {
        t.Fatal("expected an error for 500 response, got nil")
    }
}

func TestFetchArtists_InvalidJSON(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("not valid json"))
    }))
    defer server.Close()

    originalURL := artistsAPIURL
    artistsAPIURL = server.URL
    defer func() { artistsAPIURL = originalURL }()

    _, err := FetchArtists()
    if err == nil {
        t.Fatal("expected a decode error, got nil")
    }
}

func TestGetArtistByID_Found(t *testing.T) {
    jsonResponse := `[
        {"id": 1, "name": "Queen"},
        {"id": 2, "name": "Scorpions"}
    ]`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(jsonResponse))
    }))
    defer server.Close()

    originalURL := artistsAPIURL
    artistsAPIURL = server.URL
    defer func() { artistsAPIURL = originalURL }()

    artist, err := GetArtistByID(2)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if artist.Name != "Scorpions" {
        t.Errorf("expected Scorpions, got %s", artist.Name)
    }
    if artist.ID != 2 {
        t.Errorf("expected ID 2, got %d", artist.ID)
    }
}

func TestGetArtistByID_NotFound(t *testing.T) {
    jsonResponse := `[{"id": 1, "name": "Queen"}]`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(jsonResponse))
    }))
    defer server.Close()

    originalURL := artistsAPIURL
    artistsAPIURL = server.URL
    defer func() { artistsAPIURL = originalURL }()

    _, err := GetArtistByID(99)
    if err == nil {
        t.Fatal("expected an error for non-existent ID, got nil")
    }
}