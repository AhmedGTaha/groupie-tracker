# Groupie Tracker

Groupie Tracker is a Go web application that consumes the public
[Groupie Trackers API](https://groupietrackers.herokuapp.com/api) and displays
artist, concert location, date, and relation data in a user-friendly website.

The project is built with Go's standard library and uses server-side HTML
templates, static CSS, and a small JavaScript layer for browser interactions.

## Features

- Display a list of artists and bands from the API.
- Show artist details such as image, members, start year, and first album.
- Display concert locations and dates using the relation data.
- Provide a client-server interaction such as search or filtering.
- Handle missing data, bad routes, and API errors without crashing.
- Include unit tests for API, handler, and search logic.

## Project Structure

```text
groupie-tracker/
|-- go.mod
|-- main.go
|-- README.md
|-- .gitignore
|-- internal/
|   |-- api/
|   |   |-- client.go
|   |   `-- client_test.go
|   |-- models/
|   |   `-- models.go
|   |-- handlers/
|   |   |-- handlers.go
|   |   `-- handlers_test.go
|   `-- service/
|       |-- search.go
|       `-- search_test.go
|-- templates/
|   |-- layout.html
|   |-- index.html
|   |-- artist.html
|   |-- error.html
|   `-- partials/
|       `-- artist_card.html
|-- static/
|   |-- css/
|   |   `-- style.css
|   |-- js/
|   |   `-- app.js
|   `-- img/
|       `-- .gitkeep
`-- docs/
    `-- README.md
```

## Main Packages

- `internal/api`: API client for fetching data from Groupie Trackers.
- `internal/models`: Shared data structures used across the application.
- `internal/handlers`: HTTP handlers and template rendering logic.
- `internal/service`: Application logic such as search and filtering.
- `templates`: HTML pages and reusable partials.
- `static`: CSS, JavaScript, and image assets.
- `docs`: Project brief and supporting documentation.

## Requirements

- Go 1.22 or newer
- Internet access when fetching data from the Groupie Trackers API

Only standard Go packages should be used for this project.

## Getting Started

Clone or open the repository, then move into the project folder:

```bash
cd groupie-tracker
```

Initialize or download dependencies if needed:

```bash
go mod tidy
```

Run the application:

```bash
go run .
```

Then open the local server URL shown in the terminal, commonly:

```text
http://localhost:8080
```

## Testing

Run all tests with:

```bash
go test ./...
```

Test files are planned for:

- API client behavior
- HTTP handlers
- Search and filtering logic

## API Source

The application uses:

```text
https://groupietrackers.herokuapp.com/api
```

The API contains:

- `artists`: Artist and band information.
- `locations`: Concert locations.
- `dates`: Concert dates.
- `relation`: Links artists with their locations and dates.

## Error Handling

The server should return friendly error pages for invalid routes, missing
artists, failed API requests, and unexpected server errors. The application
should not panic or crash during normal use.

## Status

This repository currently contains the planned project structure. Application
logic, templates, styles, and tests will be implemented inside the existing
folders.
