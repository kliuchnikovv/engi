package services

import (
	"fmt"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/entity"
	"github.com/KlyuchnikovV/webapi/options"
	"github.com/KlyuchnikovV/webapi/param"
	"github.com/KlyuchnikovV/webapi/validate"
)

// Example service.
type NotesAPI struct{}

type Func func(*webapi.Engine) error

func (api *NotesAPI) Prefix() string {
	return "notes"
}

func (api *NotesAPI) Middlewares() []webapi.Middleware {
	return []webapi.Middleware{
		webapi.UseCORS(webapi.AllowedOrigins("*")),
	}
}

func (api *NotesAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"create": webapi.POST(
			api.Create,
			param.Body(new(entity.NotesRequest)),
			param.Description("creates new note"),
		),
		"get": webapi.GET(
			api.GetByID,
			param.Integer("id", options.InQuery,
				validate.AND(validate.Greater(1), validate.Less(10)),
				options.Description("ID of request."),
			),
		),
		"{object}/{id}": webapi.GET(
			api.GetByIDFromPath,
			param.Integer("id", options.InPath),
			param.String("object", options.InPath),
		),
	}
}

func (api *NotesAPI) Create(ctx *webapi.Context) error {
	if body := ctx.Request.Body(); body != nil {
		return ctx.OK(body)
	}

	return ctx.Created()
}

func (api *NotesAPI) GetByID(ctx *webapi.Context) error {
	var id = ctx.Request.Integer("id", options.InQuery)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	// TODO: failed
	return ctx.OK(struct {
		Message string `json:"message" description:"Response message"`
	}{
		Message: fmt.Sprintf("got id: '%d'", id),
	})
}

func (api *NotesAPI) GetByIDFromPath(ctx *webapi.Context) error {
	var (
		id     = ctx.Request.Integer("id", options.InPath)
		object = ctx.Request.String("object", options.InPath)
	)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.OK(fmt.Sprintf("got id for object '%s': '%d'", object, id))
}
