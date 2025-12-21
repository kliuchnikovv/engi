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

Engi is a lightweight Go framework designed to give developers the most intuitive and declarative API for building HTTP services. Its key idea is structured services, automatic parameter validation, convenient response wrappers, and an embedded documentation layer—so that the framework, rather than the development team, handles as much boilerplate code as possible.

### Installation

```sh
go get github.com/kliuchnikovv/engi
```
### Example of usage

The idea of this framework is to create **services**, each of which works with one model.

```golang
type NotesAPI struct {\n  notesStore store.NoteStore\n}\n
func (api *NotesAPI) Prefix() string {\n  // All requests to NotesAPI will have prefix \"/notes\" in their url path, e.g. /{engi-prefix}/notes/{api-route}\n  return \"notes\"\n}\n```

Each service must implement 2 methods: `Prefix` and `Routers`:

- `Prefix` gives route prefix and serves as name of your service;
- `Routers` defines handlers, their paths and their mandatory parameters;

The handler described as a **relative** path to the handler wrapped in a request method (`POST`, `GET` ...<!--(godoc link?)-->)
with additional middleware functions, including those for requesting mandatory parameters:

```golang
func (api *NotesAPI) Routers() engi.Routes {\n  return engi.Routes{\n    ...\n    engi.GET(\"{id}\"): engi.Handle( // Using GET method to get a note with {id} using path \"{url}/notes/{id}\".
      api.Get,                                   // Handler to handle this request.\n      path.Integer(\"id\", validate.Greater(0)),   // Path parameter to be parsed into integer.\n      middlewares.Description(\"get note by id\"), // Description of this route for documentation purposes.\n    ),\n    ...\n  }\n}\n```

Further, when requesting, all the necessary parameters will be checked for the presence and type (if the required parameter is missing, `BadRequest` error will be returned) and then will be available for use in handlers through the context `ctx`. <!--(godoc link?)-->

Also, through the context `ctx`<!--(godoc link?)-->, you can form a result or an error using predefined functions for the most used answers:

```golang
func (api *NotesAPI) Get(\n  ctx context.Context,    // Standart Golang context\n  request engi.Request,   // Request - contains all metadata about request itself and parsed parameters that you described in Routers method.\n  response engi.Response, // Response - contains wrappers for easy and painless response creation without any formatting and manipulating with headers and statuses.\n) error {\n  // Extract variable from path and cast in to integer.\n  // If there is no possibility to cast, BadRequest error will be returned automatically.\n  var id = request.Integer(\"id\", placing.InPath) \n\n  note, err := api.notesStore.GetByID(ctx, id)\n  if err != nil {\n    // In order to return error, use predefined wrappers.\n    // They will automatically wrap error and marshal it using settings from middlewares and engine.\n    return response.NotFound(err.Error())\n  }\n\n  // In order to return result, use predefined wrappers.\n  // They will automatically wrap object and marshal it using settings from middlewares and engine.\n  return response.OK(note)\n}\n```

As a result, to create an application, it remains to create server with `engi.New` passing tcp address and global (for every handler) prefix, register service and start the api.

```golang
func main() {\n  // 1. Create engine\n  var engine = engi.New(\":8080\", // Defines address to listen.\n    engi.WithPrefix(\"api\"),    // Defines global prefix for all routes.\n    engi.ResponseAsJSON(       // Defines all responses to be marshaled as JSON objects.\n      response.AsIs,         // All responses will use no wrappers and will be sent as is.\n    ),\n    engi.WithLogger(slog.NewTextHandler(os.Stdout,\n      &slog.HandlerOptions{\n        Level: slog.LevelDebug,\n      },\n    )),\n  )\n
  // 2. Register services\n  if err := engine.RegisterServices(\n    services.NewNotesAPI(*store.NewNoteStore(db)),\n  ); err != nil {\n    return err\n  }\n
  // 3. Start server - blocking call\n  if err := engine.Start(); err != nil {\n    log.Fatal(err)\n  }\n}\n```

 Workable example of this api you can found [here](https://github.com/kliuchnikovv/engi-example)