# Engi

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/kliuchnikovv/engi/go.yml?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/kliuchnikovv/engi?style=for-the-badge)](https://goreportcard.com/report/github.com/kliuchnikovv/engi)
![GitHub gso.mod Go version](https://img.shields.io/github/go-mod/go-version/kliuchnikovv/engi?style=for-the-badge)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/kliuchnikovv/engi)
[![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2Fkliuchnikovv%2Fengi&label=Views&labelColor=%23697689&countColor=%23555555)](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2Fkliuchnikovv%2Fengi)
![GitHub](https://img.shields.io/github/license/kliuchnikovv/engi?style=for-the-badge)


## A web framework that prioritizes developer usability.

### Description
This framework forces developer to write more structured, human-centric code.

### Installation

```sh
go get github.com/kliuchnikovv/engi
```
### Example of usage

The idea of this framework is to create **services**, each of which works with one model.

```golang
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

The handler described as a **relative** path to the handler wrapped in a request method (`POST`, `GET` ...<!--(godoc link?)-->)
with additional middleware functions, including those for requesting mandatory parameters:

```golang
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

Further, when requesting, all the necessary parameters will be checked for the presence and type (if the required parameter is missing, `BadRequest` error will be returned) and then will be available for use in handlers through the context `ctx`. <!--(godoc link?)-->

Also, through the context `ctx`<!--(godoc link?)-->, you can form a result or an error using predefined functions for the most used answers:

```golang
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

```golang
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

Workable example of this api you can found [here](https://github.com/kliuchnikovv/engi-example)
