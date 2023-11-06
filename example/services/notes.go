package services

import (
	"fmt"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/example/entity"
	"github.com/KlyuchnikovV/engi/options"
	"github.com/KlyuchnikovV/engi/parameter"
	"github.com/KlyuchnikovV/engi/placing"
	"github.com/KlyuchnikovV/engi/validate"
)

// Example service.
type NotesAPI struct{}

type Func func(*engi.Engine) error

func (api *NotesAPI) Prefix() string {
	return "notes"
}

func (api *NotesAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		engi.UseCORS(engi.AllowedOrigins("*")),
	}
}

func (api *NotesAPI) Routers() map[string]engi.RouterByPath {
	return map[string]engi.RouterByPath{
		"create": engi.POST(
			api.Create,
			parameter.Body(new(entity.NotesRequest)),
			parameter.Description("creates new note"),
		),
		"get/{id}": engi.GET(
			api.GetByID,
			parameter.Integer("id", placing.InPath,
				options.Description("ID of request."),
				validate.AND(
					validate.Greater(1),
					validate.Less(10),
				),
			),
		),
		"{object}/{id}": engi.GET(
			api.GetByIDFromPath,
			parameter.Integer("id", placing.InPath),
			parameter.String("object", placing.InPath),
		),
	}
}

func (api *NotesAPI) Create(ctx engi.Context) error {
	if body := ctx.Body(); body != nil {
		return ctx.OK(body)
	}

	return ctx.Created()
}

func (api *NotesAPI) GetByID(ctx engi.Context) error {
	var id = ctx.Integer("id", placing.InPath)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.OK(struct {
		Message string `json:"message" description:"Response message"`
	}{
		Message: fmt.Sprintf("got id: '%d'", id),
	})
}

func (api *NotesAPI) GetByIDFromPath(ctx engi.Context) error {
	var (
		id     = ctx.Integer("id", placing.InPath)
		object = ctx.String("object", placing.InPath)
	)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.OK(fmt.Sprintf("got id for object '%s': '%d'", object, id))
}
