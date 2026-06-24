# Groupie Tracker

Groupie Tracker is a Go web application that displays music artist and band information from the public [Groupie Trackers API](https://groupietrackers.herokuapp.com/api). It lets users browse artists, search by artist name, and view detail pages with members, albums, creation dates, concert locations, concert dates, and location/date relations.

## Features

- Home page with navigation to the artist directory.
- Artist directory populated from the external API.
- Case-insensitive artist name search using the `/artists?q=` query parameter.
- Artist detail page with:
  - artist image
  - members
  - creation date
  - first album
  - concert locations
  - concert dates
  - grouped location/date relations
- Server-side HTML rendering with Go templates.
- Static CSS served from the `static` directory.
- Unit tests for API fetching/decoding behavior and search matching.

## Tech Stack

- Go
- Go standard library only
- HTML templates
- CSS
- Groupie Trackers API

## Project Structure

```text
.
|-- api/                 # API fetchers and JSON decoders
|-- docs/                # Project requirement notes
|-- models/              # Data structs for artists, dates, locations, relations
|-- static/              # CSS assets
|-- templates/           # HTML templates
|-- go.mod               # Go module definition
|-- main.go              # HTTP server, routes, and page handlers
|-- main_test.go         # Tests for app helper functions
`-- README.md
```

## Requirements

- Go 1.26 as configured in `go.mod`
- Internet access while running the app, because artist data is fetched from `https://groupietrackers.herokuapp.com/api`

## Getting Started

Clone the repository, then run the server from the project root:

```bash
go run .
```

Open the app in your browser:

```text
http://localhost:8080
```

## Routes

| Route | Description |
| --- | --- |
| `/` | Home page |
| `/artists` | Artist directory |
| `/artists?q=queen` | Artist directory filtered by name |
| `/artist?id=1` | Detail page for a single artist |
| `/static/` | Static file assets |

## Running Tests

Run all tests with:

```bash
go test ./...
```

The test suite includes:

- API request success and error handling tests using `httptest`
- JSON decoding failure checks
- artist lookup checks
- case-insensitive search helper tests

## Data Source

The application fetches data from these Groupie Trackers endpoints:

- `https://groupietrackers.herokuapp.com/api/artists`
- `https://groupietrackers.herokuapp.com/api/locations`
- `https://groupietrackers.herokuapp.com/api/dates`
- `https://groupietrackers.herokuapp.com/api/relation`

## Notes

- The server listens on port `8080`.
- API requests use a 10-second timeout.
- If the external API cannot be reached, the artists page returns a friendly status message instead of crashing.
