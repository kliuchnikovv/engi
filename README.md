# WebApi

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/KlyuchnikovV/webapi/Go?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/KlyuchnikovV/webapi?style=for-the-badge)](https://goreportcard.com/report/github.com/KlyuchnikovV/webapi)
![GitHub](https://img.shields.io/github/license/KlyuchnikovV/webapi?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/KlyuchnikovV/webapi?style=for-the-badge)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/KlyuchnikovV/webapi)
![Visitors](https://visitor-badges.glitch.me?username=KlyuchnikovV&repo=webapi&style=for-the-badge&label=Views)


## A web framework that prioritizes developer usability.

### Description
This framework aims to write more structured, human-centric code.

### Installation

```sh
go get github.com/KlyuchnikovV/webapi
```
### Example of usage

The idea of this framework is to create **services**, each of which works with one model.

```golang
type RequestAPI struct {
    webapi.API
}

func NewRequestAPI() webapi.API {
    return &RequestAPI{
        // 'request' string is a prefix to the query
        // so full path to GetByID handler will be '/api/request/get'
		API: webapi.New("request"),
	}
}
```

For each service, the `Routers` method defined, which gives handlers upon registration.

The handler described as a **relative** path to the handler wrapped in a request method (`POST`, `GET` ...<!--(godoc link?)-->)
with additional middleware functions, including those for requesting mandatory query parameters.

```golang
func (api *RequestAPI) Routers() map[string]webapi.RouterFunc {
	return map[string]webapi.RouterFunc{
		"get":    api.GET(api.GetByID, api.WithInt("id")),
		"create": api.POST(api.Create, api.Body(&Body{})),
	}
}
```

Further, when requesting, all the necessary parameters will be checked for the presence and type (if the required parameter is missing, BadRequest will be returned) and then will be available for use in handlers through the context `ctx.QueryParams`. <!--(godoc link?)-->

Also, through the context `ctx.Response`<!--(godoc link?)-->, you can form a result or an error using predefined functions.

```golang
func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
    var id = ctx.QueryParams.Integer("id")

    // Do something with id (we will check it)
    if id < 0 {
        return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
    }

    return ctx.Response.OK(fmt.Sprintf("got id: '%d'", id))
}
```

As a result, to create an application, it remains to register the service and run the api.

```golang
func main() {
    w := webapi.New()

    w.RegisterServices(
        service.NewRequestAPI(),
    )

    w.Start("localhost:8080")
}
```

Workable example of this api you can found [here](https://github.com/KlyuchnikovV/webapi/tree/main/example)
