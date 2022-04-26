package service

import (
	"fmt"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/entity"
	"github.com/KlyuchnikovV/webapi/param"
)

// Example service.
type NotesAPI struct {
	webapi.Service
}

func NewNotesAPI(engine *webapi.Engine) *NotesAPI {
	return &NotesAPI{
		Service: *webapi.NewService(engine, "notes"),
	}
}

// Should contain only return statement for documentation.
func (api *NotesAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"{object}/{id}": api.GET(
			api.GetByIDFromPath,
			param.InPathInteger("id"),
			param.InPathString("object"),
		),
		"get": api.GET(
			api.GetByID,
			param.QueryInteger("id",
				param.Description("ID of request."),
				param.AND(param.Greater(1), param.Less(10)),
			),
		),
		"create": api.POST(
			api.Create,
			param.Body(&entity.NotesRequest{}),
		),
	}
}

func (api *NotesAPI) Create(ctx *webapi.Context) error {
	if body := ctx.Body(); body != nil {
		return ctx.Response.OK(body)
	}

	ctx.Response.Created()

	return nil
}

func (api *NotesAPI) GetByID(ctx *webapi.Context) error {
	var id = ctx.Request.QueryInteger("id")

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.Response.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *NotesAPI) GetByIDFromPath(ctx *webapi.Context) error {
	var (
		id     = ctx.Request.InPathInteger("id")
		object = ctx.Request.InPathString("object")
	)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.Response.OK(fmt.Sprintf("got id for object '%s': '%d'", object, id))
}
