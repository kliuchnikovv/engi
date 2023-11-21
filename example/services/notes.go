package services

import (
	"context"
	"fmt"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/example/entity"
	"github.com/KlyuchnikovV/engi/parameter"
	"github.com/KlyuchnikovV/engi/parameter/path"
	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/validate"
)

// Example service.
type NotesAPI struct{}

func (api *NotesAPI) Prefix() string {
	return "notes"
}

func (api *NotesAPI) Middlewares() []engi.Register {
	return []engi.Register{
		engi.UseCORS(engi.AllowedOrigins("*")),
		engi.UseAuthorization(engi.BasicAuth("Dave", "IsCrazy")),
	}
}

func (api *NotesAPI) Routers() engi.Routes {
	return engi.Routes{
		"create": engi.POST(api.Create,
			parameter.Body(new(entity.NotesRequest)),
			engi.UseAuthorization(engi.BasicAuth("Dave", "NotCrazy")),
		),
		"get/{id}": engi.GET(api.GetByID,
			path.Integer("id",
				validate.AND(validate.Greater(1), validate.Less(10)),
			),
			engi.UseAuthorization(engi.BearerAuth(func(s string) bool { return s == "token" })),
		),
		"{object}/{id}": engi.GET(api.GetByIDFromPath,
			path.Integer("id"),
			path.String("object"),
		),
	}
}

func (api *NotesAPI) Create(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	if body := request.Body(); body != nil {
		return response.OK(body)
	}

	return response.Created()
}

func (api *NotesAPI) GetByID(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var id = request.Integer("id", placing.InPath)

	// Do something with id (we will check it)
	if id < 0 {
		return response.BadRequest("id can't be negative (got: %d)", id)
	}

	return response.OK(struct {
		Message string `description:"Response message" json:"message"`
	}{
		Message: fmt.Sprintf("got id: '%d'", id),
	})
}

func (api *NotesAPI) GetByIDFromPath(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var (
		id     = request.Integer("id", placing.InPath)
		object = request.String("object", placing.InPath)
	)

	// Do something with id (we will check it)
	if id < 0 {
		return response.BadRequest("id can't be negative (got: %d)", id)
	}

	return response.OK(fmt.Sprintf("got id for object '%s': '%d'", object, id))
}
