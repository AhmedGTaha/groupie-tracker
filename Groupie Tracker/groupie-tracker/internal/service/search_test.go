package service

import (
	"errors"
	"testing"

	"groupie-tracker/internal/models"
)

type fakeAPIClient struct {
	artists   []models.Artist
	locations models.LocationIndex
	dates     models.DateIndex
	relations models.RelationIndex
	err       error
}

func (f fakeAPIClient) FetchArtists() ([]models.Artist, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.artists, nil
}

func (f fakeAPIClient) FetchLocations() (models.LocationIndex, error) {
	if f.err != nil {
		return models.LocationIndex{}, f.err
	}
	return f.locations, nil
}

func (f fakeAPIClient) FetchDates() (models.DateIndex, error) {
	if f.err != nil {
		return models.DateIndex{}, f.err
	}
	return f.dates, nil
}

func (f fakeAPIClient) FetchRelations() (models.RelationIndex, error) {
	if f.err != nil {
		return models.RelationIndex{}, f.err
	}
	return f.relations, nil
}

func TestFindArtistByIDReturnsDetails(t *testing.T) {
	svc := New(fakeAPIClient{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
		locations: models.LocationIndex{
			Index: []models.Location{{ID: 1, Locations: []string{"london-uk"}}},
		},
		dates: models.DateIndex{
			Index: []models.Date{{ID: 1, Dates: []string{"01-01-2026"}}},
		},
		relations: models.RelationIndex{
			Index: []models.Relation{{ID: 1, DatesLocations: map[string][]string{"london-uk": {"01-01-2026"}}}},
		},
	})

	details, err := svc.FindArtistByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if details.Artist.Name != "Queen" {
		t.Fatalf("expected Queen, got %q", details.Artist.Name)
	}
	if details.Locations[0] != "london-uk" {
		t.Fatalf("expected location london-uk, got %q", details.Locations[0])
	}
	if details.DatesLocations["london-uk"][0] != "01-01-2026" {
		t.Fatalf("expected relation date, got %q", details.DatesLocations["london-uk"][0])
	}
}

func TestFindArtistByIDReturnsNotFound(t *testing.T) {
	svc := New(fakeAPIClient{
		artists: []models.Artist{{ID: 1, Name: "Queen"}},
	})

	_, err := svc.FindArtistByID(99)
	if !errors.Is(err, ErrArtistNotFound) {
		t.Fatalf("expected ErrArtistNotFound, got %v", err)
	}
}

func TestSearchArtistsMatchesNameMemberAlbumYearAndLocation(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen", Members: []string{"Freddie Mercury"}, CreationDate: 1970, FirstAlbum: "14-12-1973"},
		{ID: 2, Name: "Daft Punk", Members: []string{"Thomas Bangalter"}, CreationDate: 1993, FirstAlbum: "Homework"},
	}
	locations := models.LocationIndex{
		Index: []models.Location{
			{ID: 1, Locations: []string{"london-uk"}},
			{ID: 2, Locations: []string{"paris-france"}},
		},
	}

	tests := []struct {
		name string
		q    string
		want string
	}{
		{name: "artist name", q: "queen", want: "Queen"},
		{name: "member", q: "thomas", want: "Daft Punk"},
		{name: "album", q: "homework", want: "Daft Punk"},
		{name: "year", q: "1970", want: "Queen"},
		{name: "location", q: "paris", want: "Daft Punk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchArtists(artists, locations, tt.q)
			if len(got) != 1 || got[0].Name != tt.want {
				t.Fatalf("expected %s, got %#v", tt.want, got)
			}
		})
	}
}

func TestSearchArtistsEmptyQueryReturnsAll(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "Queen"},
		{ID: 2, Name: "Daft Punk"},
	}

	got := SearchArtists(artists, models.LocationIndex{}, "")
	if len(got) != 2 {
		t.Fatalf("expected all artists, got %d", len(got))
	}
}
