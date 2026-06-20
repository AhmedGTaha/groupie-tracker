package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"groupie-tracker/internal/models"
)

var ErrArtistNotFound = errors.New("artist not found")

type APIClient interface {
	FetchArtists() ([]models.Artist, error)
	FetchLocations() (models.LocationIndex, error)
	FetchDates() (models.DateIndex, error)
	FetchRelations() (models.RelationIndex, error)
}

type Summary struct {
	ArtistCount   int
	LocationCount int
	DateCount     int
	RelationCount int
	FirstArtist   string
	FirstLocation int
	FirstDate     int
	FirstRelation int
}

type Data struct {
	Artists   []models.Artist
	Locations models.LocationIndex
	Dates     models.DateIndex
	Relations models.RelationIndex
}

type Service struct {
	client APIClient
	ttl    time.Duration

	mu       sync.RWMutex
	cache    Data
	cachedAt time.Time
	loaded   bool
}

func New(client APIClient) *Service {
	return &Service{
		client: client,
		ttl:    5 * time.Minute,
	}
}

func (s *Service) AllArtists() ([]models.Artist, error) {
	data, err := s.LoadData()
	if err != nil {
		return nil, err
	}
	return data.Artists, nil
}

func (s *Service) LoadData() (Data, error) {
	s.mu.RLock()
	if s.loaded && time.Since(s.cachedAt) < s.ttl {
		data := s.cache
		s.mu.RUnlock()
		return data, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.loaded && time.Since(s.cachedAt) < s.ttl {
		return s.cache, nil
	}

	artists, err := s.client.FetchArtists()
	if err != nil {
		return Data{}, fmt.Errorf("fetch artists: %w", err)
	}

	locations, err := s.client.FetchLocations()
	if err != nil {
		return Data{}, fmt.Errorf("fetch locations: %w", err)
	}

	dates, err := s.client.FetchDates()
	if err != nil {
		return Data{}, fmt.Errorf("fetch dates: %w", err)
	}

	relations, err := s.client.FetchRelations()
	if err != nil {
		return Data{}, fmt.Errorf("fetch relations: %w", err)
	}

	s.cache = Data{
		Artists:   artists,
		Locations: locations,
		Dates:     dates,
		Relations: relations,
	}
	s.cachedAt = time.Now()
	s.loaded = true

	return s.cache, nil
}

func (s *Service) FindArtistByID(id int) (models.ArtistDetails, error) {
	data, err := s.LoadData()
	if err != nil {
		return models.ArtistDetails{}, err
	}

	artist, ok := findArtistByID(data.Artists, id)
	if !ok {
		return models.ArtistDetails{}, ErrArtistNotFound
	}

	locationByID := locationsByID(data.Locations.Index)
	dateByID := datesByID(data.Dates.Index)
	relationByID := relationsByID(data.Relations.Index)

	return models.ArtistDetails{
		Artist:         artist,
		Locations:      locationByID[id],
		Dates:          dateByID[id],
		DatesLocations: relationByID[id],
	}, nil
}

func (s *Service) SearchArtists(query string) ([]models.Artist, error) {
	data, err := s.LoadData()
	if err != nil {
		return nil, err
	}

	return SearchArtists(data.Artists, data.Locations, query), nil
}

func (s *Service) Summary() (Summary, error) {
	data, err := s.LoadData()
	if err != nil {
		return Summary{}, err
	}

	summary := Summary{
		ArtistCount:   len(data.Artists),
		LocationCount: len(data.Locations.Index),
		DateCount:     len(data.Dates.Index),
		RelationCount: len(data.Relations.Index),
	}

	if len(data.Artists) > 0 {
		summary.FirstArtist = data.Artists[0].Name
	}
	if len(data.Locations.Index) > 0 {
		summary.FirstLocation = data.Locations.Index[0].ID
	}
	if len(data.Dates.Index) > 0 {
		summary.FirstDate = data.Dates.Index[0].ID
	}
	if len(data.Relations.Index) > 0 {
		summary.FirstRelation = data.Relations.Index[0].ID
	}

	return summary, nil
}

func SearchArtists(artists []models.Artist, locationIndex models.LocationIndex, query string) []models.Artist {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return artists
	}

	locationByID := locationsByID(locationIndex.Index)
	results := make([]models.Artist, 0)

	for _, artist := range artists {
		if artistMatches(artist, locationByID[artist.ID], query) {
			results = append(results, artist)
		}
	}

	return results
}

func findArtistByID(artists []models.Artist, id int) (models.Artist, bool) {
	for _, artist := range artists {
		if artist.ID == id {
			return artist, true
		}
	}
	return models.Artist{}, false
}

func artistMatches(artist models.Artist, locations []string, query string) bool {
	if strings.Contains(strings.ToLower(artist.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
		return true
	}
	if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
		return true
	}

	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true
		}
	}

	for _, location := range locations {
		if strings.Contains(strings.ToLower(location), query) {
			return true
		}
	}

	return false
}

func locationsByID(locations []models.Location) map[int][]string {
	result := make(map[int][]string, len(locations))
	for _, location := range locations {
		result[location.ID] = location.Locations
	}
	return result
}

func datesByID(dates []models.Date) map[int][]string {
	result := make(map[int][]string, len(dates))
	for _, date := range dates {
		result[date.ID] = date.Dates
	}
	return result
}

func relationsByID(relations []models.Relation) map[int]map[string][]string {
	result := make(map[int]map[string][]string, len(relations))
	for _, relation := range relations {
		result[relation.ID] = relation.DatesLocations
	}
	return result
}
