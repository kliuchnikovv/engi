package services

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/example/entity"
	"github.com/KlyuchnikovV/engi/options"
	"github.com/KlyuchnikovV/engi/parameter"
	"github.com/KlyuchnikovV/engi/placing"
	"github.com/KlyuchnikovV/engi/validate"
)

// Example service.
type RequestAPI struct{}

func (api *RequestAPI) Prefix() string {
	return "request"
}

func (api *RequestAPI) Middlewares() []engi.Middleware {
	return []engi.Middleware{
		engi.UseCORS(engi.AllowedOrigins("*")),
	}
}

func (api *RequestAPI) Routers() map[string]engi.RouterByPath {
	return map[string]engi.RouterByPath{
		"get": engi.GET(
			api.GetByID,
			parameter.Integer("id", placing.InQuery,
				options.Description("ID of request."),
				validate.AND(validate.Greater(1), validate.Less(10)),
			),
		),
		"create": engi.POST(
			api.Create,
			parameter.Description("Creates new request"),
			parameter.Body(new(entity.RequestBody),
				options.Description("Body description"),
			),
		),
		"create/sub-request": engi.POST(
			api.CreateSubRequest,
			parameter.Body([]entity.RequestBody{}),
		),
		"filter": engi.GET(
			api.Filter,
			parameter.Bool("bool", placing.InQuery),
			parameter.Float("float", placing.InQuery, validate.NotEmpty),
			parameter.Integer("int", placing.InQuery),
			parameter.String("str", placing.InQuery,
				validate.AND(validate.NotEmpty, validate.Greater(2)),
			),
			parameter.Time("time", "2006-01-02 15:04", placing.InQuery,
				options.Description("Filter by time field."),
			),
		),
	}
}

func (api *RequestAPI) Create(ctx engi.Context) error {
	if body := ctx.Body(); body != nil {
		return ctx.OK(body)
	}

	return ctx.Created()
}

func (api *RequestAPI) CreateSubRequest(ctx engi.Context) error {
	return ctx.Object(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", []entity.RequestBody{}),
	)
}

func (api *RequestAPI) GetByID(ctx engi.Context) error {
	var id = ctx.Integer("id", placing.InQuery)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(ctx engi.Context) error {
	var (
		i     = ctx.Integer("int", placing.InQuery)
		str   = ctx.String("str", placing.InQuery)
		t     = ctx.Time("time", "2006-01-02 15:04", placing.InQuery)
		b     = ctx.Bool("bool", placing.InQuery)
		float = ctx.Float("float", placing.InQuery)
	)

	return ctx.OK(fmt.Sprintf(
		"filtered by id: '%d' and field: %s, time: %s, isAssignable: %t and float: %f",
		i, str, t.Format("15:04 02/01/2006"), b, float,
	))
}
