# Groupie Tracker

Groupie Tracker is a Go web application that uses the public Groupie Trackers API to display artists, bands, concert locations, concert dates, and location/date relations.

The app uses only Go's standard library: `net/http`, `html/template`, `encoding/json`, `sync`, and other built-in packages.

## Features

- Home page with artist cards from the real API
- Artist detail page with image, members, creation date, first album, locations, dates, and relation table
- Server-side search at `/search?q=...`
- Friendly error page for bad routes and invalid artist IDs
- Static CSS served from `/static/`
- API client with timeout, status checks, JSON decoding, and error returns
- Service layer with simple in-memory caching
- Unit tests for API client, service search/detail logic, and HTTP handlers

## Routes

- `GET /` - all artists
- `GET /search?q=query` - filtered artists
- `GET /artist?id=1` - artist details
- `GET /artist/1` - artist details alternate URL
- `GET /api-test` - debug API summary
- `GET /static/...` - static files

## Project Structure

```text
groupie-tracker/
|-- go.mod
|-- main.go
|-- internal/
|   |-- api/
|   |   |-- client.go
|   |   `-- client_test.go
|   |-- handlers/
|   |   |-- handlers.go
|   |   `-- handlers_test.go
|   |-- models/
|   |   `-- models.go
|   `-- service/
|       |-- search.go
|       `-- search_test.go
|-- templates/
|   |-- index.html
|   |-- artist.html
|   |-- error.html
|   `-- layout.html
`-- static/
    |-- css/
    |   `-- style.css
    `-- js/
        `-- app.js
```

## Run

```bash
go run .
```

Then open:

```text
http://localhost:8080
```

## Test

```bash
go test ./...
```

## API

The app fetches data from:

```text
https://groupietrackers.herokuapp.com/api
```

Used endpoints:

- `/artists`
- `/locations`
- `/dates`
- `/relation`

If the API is down or returns invalid data, the app returns a friendly error response instead of crashing.
