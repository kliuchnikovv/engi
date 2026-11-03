# Engi

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kliuchnikovv/engi/go.yml?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/kliuchnikovv/engi?style=for-the-badge)](https://goreportcard.com/report/github.com/kliuchnikovv/engi)
![GitHub gso.mod Go version](https://img.shields.io/github/go-mod/go-version/kliuchnikovv/engi?style=for-the-badge)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/kliuchnikovv/engi)
[![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2Fkliuchnikovv%2Fengi&label=Views&labelColor=%23697689&countColor=%23555555)](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2Fkliuchnikovv%2Fengi)
![GitHub](https://img.shields.io/github/license/kliuchnikovv/engi?style=for-the-badge)

## A web framework that prioritizes developer usability.

### Description

Engi is a lightweight Go framework designed to give developers the most intuitive and
declarative API for building HTTP services. Its key idea is structured services,
automatic parameter validation, convenient response wrappers, and an embedded
documentation layer — so that the framework, rather than the development team, handles
as much boilerplate code as possible.

### Who is this for

Engi is built for Go developers who write HTTP APIs and want the routing, parameter
handling and response plumbing to stay out of their handlers.

It fits well when you want:

- **Declarative routes** — a handler declares the parameters it needs, and the framework
  parses and validates them before the handler runs.
- **A small dependency footprint** — the only non-standard-library dependencies are the
  OpenTelemetry packages, and tracing is opt-in.
- **Structure that scales with the service** — a service owns its prefix, its routes and
  its middleware, so a large API stays split into small, self-contained units.
- **Batteries, but replaceable** — CORS, Basic/Bearer/API-key auth and response
  wrapping ship with the framework, and every one of them is an ordinary middleware you
  can swap for your own.

It is probably *not* the right choice if you need gRPC or WebSocket support today, if
you want a large third-party middleware ecosystem, or if you want a framework that
opinionates hard on project layout.

### Features

- **Trie-based router** — static segments, `{param}` segments and a trailing `*wildcard`,
  with static segments taking precedence over parameterized ones.
- **Service-oriented design** — every service owns a prefix, a set of routes, and its
  own middleware.
- **Automatic parameter handling** — parameters are declared next to the route, parsed
  into the right Go type, and checked with composable validators.
- **Middleware pipeline** — middlewares are ordered by priority, service-level
  middlewares compose with route-level ones.
- **Response wrappers** — `response.OK`, `response.NotFound` and friends; JSON and XML
  marshaling.
- **CORS and auth middleware** — Basic, Bearer token, API key, and per-service or
  per-route configuration.
- **OpenTelemetry tracing** — the HTTP server is instrumented out of the box.
- **Graceful shutdown** — `SIGINT`/`SIGTERM` handled with a bounded drain.

### Installation

```sh
go get github.com/kliuchnikovv/engi
```

Engi requires Go 1.22 or newer.

### Quick start

A complete service, from definition to a running server:

```go
package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/middlewares"
	"github.com/kliuchnikovv/engi/definition/parameter/path"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/definition/response"
	"github.com/kliuchnikovv/engi/definition/validate"
)

// Note is the model this service works with.
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// NotesAPI serves the notes resource.
type NotesAPI struct{}

// Prefix is the URL prefix of every route of this service.
func (api *NotesAPI) Prefix() string {
	return "notes"
}

// Routers declares the routes of this service and the parameters they require.
func (api *NotesAPI) Routers() engi.Routes {
	return engi.Routes{
		engi.GET("{id}"): engi.Handle(
			api.Get,
			path.Integer("id", validate.Greater(0)),
			middlewares.Description("get note by id"),
		),
	}
}

// Get handles GET /{prefix}/notes/{id}.
func (api *NotesAPI) Get(
	ctx context.Context,
	request engi.Request,
	response engi.Response,
) error {
	// The parameter was declared in Routers, so it is already parsed and validated.
	id := request.Integer("id", placing.InPath)

	return response.OK(Note{ID: int(id), Title: "a note"})
}

func main() {
	engine := engi.New(":8080",
		engi.WithPrefix("api"),
		engi.ResponseAsJSON(response.AsIs),
		engi.WithLogger(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		)),
	)

	if err := engine.RegisterServices(&NotesAPI{}); err != nil {
		log.Fatal(err)
	}

	// Start blocks until the server is shut down.
	if err := engine.Start(); err != nil {
		log.Fatal(err)
	}
}
```

Run it and call the route:

```sh
go run main.go
curl http://localhost:8080/api/notes/42
```

```
{"id":42,"title":"a note"}
```

A request for a non-numeric id is rejected before the handler runs, because `id` was
declared with `path.Integer`:

```sh
curl -i http://localhost:8080/api/notes/abc
```

```
HTTP/1.1 400 Bad Request
```

### How it works

The idea of this framework is to create **services**, each of which works with one model.

```go
type NotesAPI struct {
  notesStore store.NoteStore
}

func (api *NotesAPI) Prefix() string {
  // All requests to NotesAPI will have prefix "/notes" in their url path, e.g. /{engi-prefix}/notes/{api-route}
  return "notes"
}
```

Each service must implement 2 methods: `Prefix` and `Routers`:

- `Prefix` gives route prefix and serves as name of your service;
- `Routers` defines handlers, their paths and their mandatory parameters;

The handler described as a **relative** path to the handler wrapped in a request method
(`POST`, `GET` ...) with additional middleware functions, including those for requesting
mandatory parameters:

```go
func (api *NotesAPI) Routers() engi.Routes {
  return engi.Routes{
    ...
    engi.GET("{id}"): engi.Handle( // Using GET method to get a note with {id} using path "{url}/notes/{id}".
      api.Get,                                   // Handler to handle this request.
      path.Integer("id", validate.Greater(0)),   // Path parameter to be parsed into integer.
      middlewares.Description("get note by id"), // Description of this route for documentation purposes.
    ),
    ...
  }
}
```

Further, when requesting, all the necessary parameters will be checked for the presence and type (if the required parameter is missing, `BadRequest` error will be returned) and then will be available for use in handlers through the context `ctx`.

Also, through the context `ctx`, you can form a result or an error using predefined functions for the most used answers:

```go
func (api *NotesAPI) Get(
  ctx context.Context,    // Standart Golang context
  request engi.Request,   // Request - contains all metadata about request itself and parsed parameters that you described in Routers method.
  response engi.Response, // Response - contains wrappers for easy and painless response creation without any formatting and manipulating with headers and statuses.
) error {
  // Extract variable from path and cast in to integer.
  // If there is no possibility to cast, BadRequest error will be returned automatically.
  var id = request.Integer("id", placing.InPath)

  note, err := api.notesStore.GetByID(ctx, id)
  if err != nil {
    // In order to return error, use predefined wrappers.
    // They will automatically wrap error and marshal it using settings from middlewares and engine.
    return response.NotFound(err.Error())
  }

  // In order to return result, use predefined wrappers.
  // They will automatically wrap object and marshal it using settings from middlewares and engine.
  return response.OK(note)
}
```

As a result, to create an application, it remains to create server with `engi.New` passing tcp address and global (for every handler) prefix, register service and start the api.

```go
func main() {
  // 1. Create engine
  var engine = engi.New(":8080", // Defines address to listen.
    engi.WithPrefix("api"),    // Defines global prefix for all routes.
    engi.ResponseAsJSON(       // Defines all responses to be marshaled as JSON objects.
      response.AsIs,         // All responses will use no wrappers and will be sent as is.
    ),
    engi.WithLogger(slog.NewTextHandler(os.Stdout,
      &slog.HandlerOptions{
        Level: slog.LevelDebug,
      },
    )),
  )

  // 2. Register services
  if err := engine.RegisterServices(
    services.NewNotesAPI(*store.NewNoteStore(db)),
  ); err != nil {
    return err
  }

  // 3. Start server - blocking call
  if err := engine.Start(); err != nil {
    log.Fatal(err)
  }
}
```

### Middleware

Middlewares are declared next to the route they belong to, or at the service level, and
are ordered by priority:

```go
import (
  "github.com/kliuchnikovv/engi/definition/middlewares/auth"
  "github.com/kliuchnikovv/engi/definition/middlewares/cors"
)

// Service-level middleware applies to every route of the service.
func (api *NotesAPI) Middlewares() []engi.Middleware {
  return []engi.Middleware{
    cors.AllowedOrigins("*"),
    cors.AllowedHeaders("Content-Type", "Authorization"),
    cors.AllowedMethods("GET", "POST", "PUT", "DELETE"),
    auth.BearerToken(os.Getenv("API_TOKEN")),
  }
}
```

Available out of the box:

| Package | Provides |
|---|---|
| `definition/middlewares/cors` | `AllowedOrigins`, `AllowedHeaders`, `AllowedMethods` |
| `definition/middlewares/auth` | `NoAuth`, `Basic`, `BearerToken`, `BearerFunc`, `APIKey` |
| `definition/middlewares` | `Description` — attaches documentation to a route |

### Parameters

Parameters can come from the path, the query string, a header or a cookie, and are
declared with the type they should be parsed into:

```go
import (
  "github.com/kliuchnikovv/engi/definition/parameter/path"
  "github.com/kliuchnikovv/engi/definition/parameter/query"
)

engi.GET("{id}"): engi.Handle(
  api.Get,
  path.Integer("id", validate.Greater(0)),
  query.String("search"),
),
```

Validators compose with `AND` and `OR`:

```go
validate.NotEmpty,
validate.Greater(0),
validate.Less(100),
```

### Running the example

The [engi-example](https://github.com/kliuchnikovv/engi-example) repository contains a
working application with SQLite-backed storage, plus the same API implemented with
Gorilla Mux for comparison.

```sh
git clone https://github.com/kliuchnikovv/engi-example.git
cd engi-example
go run main.go
```

It starts two servers:

- `http://localhost:8080` — the API implemented with Engi
- `http://localhost:8081` — the same API implemented with Gorilla Mux

```sh
curl http://localhost:8080/notes
curl http://localhost:8081/tasks
```

### Roadmap

Planned work is tracked in [ROADMAP.md](ROADMAP.md) and mirrored as GitHub issues.
The near-term focus is correctness and stability, followed by OpenAPI documentation
generation, test coverage, and benchmarks.

## License

Released under the [MIT License](LICENSE).
