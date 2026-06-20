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
	"groupie-tracker/internal/service"
)

type fakeService struct {
	artists []models.Artist
	details models.ArtistDetails
	summary service.Summary
	err     error
}

func (f fakeService) AllArtists() ([]models.Artist, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.artists, nil
}

func (f fakeService) FindArtistByID(id int) (models.ArtistDetails, error) {
	if f.err != nil {
		return models.ArtistDetails{}, f.err
	}
	if f.details.Artist.ID != id {
		return models.ArtistDetails{}, service.ErrArtistNotFound
	}
	return f.details, nil
}

func (f fakeService) SearchArtists(query string) ([]models.Artist, error) {
	if f.err != nil {
		return nil, f.err
	}
	if strings.TrimSpace(query) == "" {
		return f.artists, nil
	}
	return []models.Artist{{ID: 2, Name: "Daft Punk"}}, nil
}

func (f fakeService) Summary() (service.Summary, error) {
	if f.err != nil {
		return service.Summary{}, f.err
	}
	return f.summary, nil
}

func useFakeService(t *testing.T, fake fakeService) {
	t.Helper()

	oldApp := app
	app = fake
	t.Cleanup(func() {
		app = oldApp
	})
}

func useTemplatePaths(t *testing.T, homeTemplate string, artistTemplate string, errorTemplate string) {
	t.Helper()

	tempDir := t.TempDir()
	oldHomeTemplatePath := homeTemplatePath
	oldArtistTemplatePath := artistTemplatePath
	oldErrorTemplatePath := errorTemplatePath

	homeTemplatePath = writeTestTemplate(t, tempDir, "index.html", homeTemplate)
	artistTemplatePath = writeTestTemplate(t, tempDir, "artist.html", artistTemplate)
	errorTemplatePath = writeTestTemplate(t, tempDir, "error.html", errorTemplate)

	t.Cleanup(func() {
		homeTemplatePath = oldHomeTemplatePath
		artistTemplatePath = oldArtistTemplatePath
		errorTemplatePath = oldErrorTemplatePath
	})
}

func writeTestTemplate(t *testing.T, dir string, name string, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test template: %v", err)
	}
	return path
}

func TestHomeHandlerRendersArtists(t *testing.T) {
	useFakeService(t, fakeService{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `{{.Title}} {{range .Artists}}{{.Name}}{{end}}`, `{{.Details.Artist.Name}}`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Queen") {
		t.Fatalf("expected rendered artist, got %q", rec.Body.String())
	}
}

func TestHomeHandlerReturnsNotFoundForUnknownPath(t *testing.T) {
	useTemplatePaths(t, `home`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHomeHandlerReturnsServerErrorWhenServiceFails(t *testing.T) {
	useFakeService(t, fakeService{err: errors.New("service failed")})
	useTemplatePaths(t, `home`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestSearchHandlerEmptyQueryShowsAllArtists(t *testing.T) {
	useFakeService(t, fakeService{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `{{range .Artists}}{{.Name}}{{end}}`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()

	SearchHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Queen") {
		t.Fatalf("expected all artists for empty query, got %q", rec.Body.String())
	}
}

func TestArtistHandlerRendersSelectedArtist(t *testing.T) {
	useFakeService(t, fakeService{
		details: models.ArtistDetails{Artist: models.Artist{ID: 2, Name: "Daft Punk"}},
	})
	useTemplatePaths(t, `home`, `{{.Details.Artist.Name}}`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist?id=2", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Daft Punk") {
		t.Fatalf("expected selected artist, got %q", rec.Body.String())
	}
}

func TestArtistHandlerAcceptsPathID(t *testing.T) {
	useFakeService(t, fakeService{
		details: models.ArtistDetails{Artist: models.Artist{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `home`, `{{.Details.Artist.Name}}`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestArtistHandlerReturnsBadRequestForMissingID(t *testing.T) {
	useTemplatePaths(t, `home`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestArtistHandlerReturnsBadRequestForInvalidID(t *testing.T) {
	useTemplatePaths(t, `home`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist?id=abc", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestArtistHandlerReturnsNotFoundForUnknownArtist(t *testing.T) {
	useFakeService(t, fakeService{
		details: models.ArtistDetails{Artist: models.Artist{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `home`, `artist`, `{{.Status}} {{.Message}}`)

	req := httptest.NewRequest(http.MethodGet, "/artist?id=99", nil)
	rec := httptest.NewRecorder()

	ArtistHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestAPITestHandlerReturnsSummary(t *testing.T) {
	useFakeService(t, fakeService{
		summary: service.Summary{
			ArtistCount:   1,
			LocationCount: 2,
			DateCount:     3,
			RelationCount: 4,
			FirstArtist:   "Queen",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api-test", nil)
	rec := httptest.NewRecorder()

	APITestHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Artists: 1") {
		t.Fatalf("expected summary, got %q", rec.Body.String())
	}
}

func TestNewRouterRoutesArtistPathRequests(t *testing.T) {
	useFakeService(t, fakeService{
		details: models.ArtistDetails{Artist: models.Artist{ID: 1, Name: "Queen"}},
	})
	useTemplatePaths(t, `home`, `{{.Details.Artist.Name}}`, `{{.Status}} {{.Message}}`)

	router := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
