# Groupie Tracker

Groupie Tracker is a Go web application that uses the public Groupie Trackers API to display artists, bands, concert locations, concert dates, and location/date relations.

The project uses only Go standard library packages such as `net/http`, `html/template`, `encoding/json`, `sync`, and `time`. There are no external backend or frontend dependencies.

## Features

- Home page with artist cards loaded from the real API
- Artist detail page with image, members, creation date, first album, locations, dates, and relation table
- Server-side search at `/search?q=...`
- Friendly error page for bad routes, invalid artist IDs, API failures, and template errors
- Static CSS served from `/static/`
- API client with timeout, status checks, JSON decoding, and safe response body closing
- Service layer with simple in-memory caching
- Unit tests for the API client, service search/detail logic, and HTTP handlers

## Routes

- `GET /` - show all artists
- `GET /search?q=query` - show filtered artists
- `GET /artist?id=1` - show one artist detail page
- `GET /artist/1` - alternate artist detail URL
- `GET /api-test` - debug API summary
- `GET /static/...` - static files

## Project Structure

```text
groupie-tracker/
|-- go.mod
|-- main.go
|-- README.md
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
|   |-- layout.html
|   `-- partials/
|       `-- artist_card.html
|-- static/
|   |-- css/
|   |   `-- style.css
|   |-- img/
|   |   `-- .gitkeep
|   `-- js/
|       `-- app.js
`-- docs/
    |-- README.md
    `-- Groupie_Tracker_SRS_and_Learning_Plan.docx
```

## Run

Default port:

```bash
go run .
```

Then open:

```text
http://localhost:8080
```

If port `8080` is already in use, choose another port:

PowerShell:

```powershell
$env:PORT="3000"
go run .
```

macOS/Linux:

```bash
PORT=3000 go run .
```

Then open:

```text
http://localhost:3000
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

If the API is down, times out, returns a non-200 status, or sends invalid JSON, the app returns a friendly error response instead of crashing.

## Notes

- Search is a normal GET form, so it works without JavaScript.
- `/api-test` is a debug route and can be removed later if the final submission should not expose it.
- `templates/layout.html`, `templates/partials/artist_card.html`, and `static/js/app.js` are present for future refactoring or enhancement, but the current app works without JavaScript.
