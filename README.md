# WebApi

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/KlyuchnikovV/webapi/go.yml?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/KlyuchnikovV/webapi?style=for-the-badge)](https://goreportcard.com/report/github.com/KlyuchnikovV/webapi)
![GitHub gso.mod Go version](https://img.shields.io/github/go-mod/go-version/KlyuchnikovV/webapi?style=for-the-badge)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/KlyuchnikovV/webapi)
[![Visitors](https://api.visitorbadge.io/api/visitors?path=https%3A%2F%2Fgithub.com%2FKlyuchnikovV%2Fwebapi&label=Views&labelColor=%23697689&countColor=%23555555)](https://visitorbadge.io/status?path=https%3A%2F%2Fgithub.com%2FKlyuchnikovV%2Fwebapi)
![GitHub](https://img.shields.io/github/license/KlyuchnikovV/webapi?style=for-the-badge)


## A web framework that prioritizes developer usability.

### Description
This framework forces developer to write more structured, human-centric code.

### Installation

```sh
go get github.com/KlyuchnikovV/webapi
```
### Example of usage

The idea of this framework is to create **services**, each of which works with one model.

```golang
type RequestAPI struct{}

func (api *RequestAPI) Prefix() string {
	return "request"
}
```

Each service must implement 2 methods: `Prefix` and `Routers`:

- `Prefix` gives route prefix and serves as name of your service;
- `Routers` defines handlers, their paths and their mandatory parameters;

The handler described as a **relative** path to the handler wrapped in a request method (`POST`, `GET` ...<!--(godoc link?)-->)
with additional middleware functions, including those for requesting mandatory parameters:

```golang
func (api *RequestAPI) Routers() map[string]webapi.RouterFunc {
	return map[string]webapi.RouterFunc{
		"get": webapi.GET(
			api.GetByID,
			parameter.Integer("id", placing.InQuery,
				options.Description("ID of the request."),
				validate.AND(validate.Greater(1), validate.Less(10)),
			),
		),
	}
}
```

Further, when requesting, all the necessary parameters will be checked for the presence and type (if the required parameter is missing, `BadRequest` error will be returned) and then will be available for use in handlers through the context `ctx`. <!--(godoc link?)-->

Also, through the context `ctx`<!--(godoc link?)-->, you can form a result or an error using predefined functions for the most used answers:

```golang
func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
    var id = ctx.Integer("id", placing.InQuery)

    // Do something with id
    if id == 5 {
        return ctx.BadRequest("id can't be '%d'", id)
    }

    return ctx.OK(
		fmt.Sprintf("got id: '%d'", id),
	)
}
```

As a result, to create an application, it remains to create server with `webapi.New` passing tcp address and global (for every handler) prefix, register service and start the api.

```golang
func main() {
   	w := webapi.New(
		":8080",
		webapi.WithPrefix("api"),
		// Define all responses as JSON object
		webapi.ResponseAsJSON(
			// Define all responses use Result field to wrap response and Error field to wrap error
			new(response.AsObject),
		),
	)

	if err := w.RegisterServices(
		new(services.RequestAPI),
	); err != nil {
		log.Fatal(err)
	}

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
```

Workable example of this api you can found [here](https://github.com/KlyuchnikovV/webapi/tree/main/example)
