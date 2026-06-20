package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"groupie-tracker/internal/models"
)

type fakeArtistClient struct {
	artists []models.Artist
	err     error
}

// FetchArtists makes fakeArtistClient satisfy the artistFetcher interface.
func (f fakeArtistClient) FetchArtists() ([]models.Artist, error) {
	return f.artists, f.err
}

// useFakeClient replaces the real API client during one test.
// t.Cleanup restores the original value after the test finishes.
func useFakeClient(t *testing.T, client fakeArtistClient) {
	t.Helper()

	oldNewAPIClient := newAPIClient
	newAPIClient = func() artistFetcher {
		return client
	}

	t.Cleanup(func() {
		newAPIClient = oldNewAPIClient
	})
}

// useTemplatePaths points the handlers at temporary test templates.
// This keeps tests independent from the real HTML files.
func useTemplatePaths(t *testing.T, homeTemplate string, artistTemplate string) {
	t.Helper()

	tempDir := t.TempDir()
	oldHomeTemplatePath := homeTemplatePath
	oldArtistTemplatePath := artistTemplatePath

	homeTemplatePath = writeTestTemplate(t, tempDir, "index.html", homeTemplate)
	artistTemplatePath = writeTestTemplate(t, tempDir, "artist.html", artistTemplate)

	t.Cleanup(func() {
		homeTemplatePath = oldHomeTemplatePath
		artistTemplatePath = oldArtistTemplatePath
	})
}

// writeTestTemplate creates one temporary template file for a test.
func writeTestTemplate(t *testing.T, dir string, name string, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test template: %v", err)
	}

	return path
}

func TestHomeHandlerRendersArtists(t *testing.T) {
	// Arrange: fake API data and a tiny template.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Daft Punk"},
		},
	})
	useTemplatePaths(t, `{{.Title}} {{range .Artists}}{{.Name}} {{end}}`, `{{.Artist.Name}}`)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Act: call the handler directly.
	HomeHandler(rec, req)

	// Assert: check status and response body.
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Groupie Tracker") || !strings.Contains(body, "Queen") {
		t.Fatalf("expected rendered home page, got %q", body)
	}
}

func TestHomeHandlerReturnsNotFoundForUnknownPath(t *testing.T) {
	// HomeHandler should only serve exactly "/".
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHomeHandlerReturnsServerErrorWhenAPIFails(t *testing.T) {
	// Simulate an API failure so the error path is tested.
	useFakeClient(t, fakeArtistClient{
		err: errors.New("api failed"),
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestArtistHandlerRendersSelectedArtist(t *testing.T) {
	// This tests the query-string style URL: /artist?id=2.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Daft Punk"},
		},
	})
	useTemplatePaths(t, `{{.Title}}`, `{{.Title}} {{.Artist.Name}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist?id=2", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Daft Punk") {
		t.Fatalf("expected selected artist page, got %q", rec.Body.String())
	}
}

func TestArtistHandlerRendersSelectedArtistFromPathID(t *testing.T) {
	// This tests the path style URL: /artist/1.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{
			{ID: 1, Name: "Queen"},
		},
	})
	useTemplatePaths(t, `{{.Title}}`, `{{.Artist.Name}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Queen") {
		t.Fatalf("expected selected artist page, got %q", rec.Body.String())
	}
}

func TestArtistHandlerReturnsBadRequestForMissingID(t *testing.T) {
	// A detail request needs either ?id=1 or /artist/1.
	req := httptest.NewRequest(http.MethodGet, "/artist", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestArtistHandlerReturnsBadRequestForInvalidID(t *testing.T) {
	// strconv.Atoi cannot convert "abc" into a number, so this should be 400.
	req := httptest.NewRequest(http.MethodGet, "/artist?id=abc", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestArtistHandlerReturnsNotFoundForUnknownArtist(t *testing.T) {
	// The ID is valid, but no artist in the fake data has ID 99.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})

	req := httptest.NewRequest(http.MethodGet, "/artist?id=99", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestAPITestHandlerReturnsArtistSummary(t *testing.T) {
	// /api-test prints a short text summary of fetched artists.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{
			{ID: 1, Name: "Queen"},
			{ID: 2, Name: "Daft Punk"},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api-test", nil)
	rec := httptest.NewRecorder()

	APITestHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Fetched 2 artists") || !strings.Contains(body, "First artist: Queen") {
		t.Fatalf("expected api summary, got %q", body)
	}
}

func TestNewRouterRoutesArtistRequests(t *testing.T) {
	// This verifies the router sends /artist?id=1 to ArtistHandler.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `{{.Title}}`, `{{.Artist.Name}}`)

	router := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/artist?id=1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestNewRouterRoutesArtistPathRequests(t *testing.T) {
	// This verifies the router sends /artist/1 to ArtistHandler.
	useFakeClient(t, fakeArtistClient{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `{{.Title}}`, `{{.Artist.Name}}`)

	router := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
