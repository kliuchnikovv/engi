package services

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/entity"
	"github.com/KlyuchnikovV/webapi/options"
	"github.com/KlyuchnikovV/webapi/parameter"
	"github.com/KlyuchnikovV/webapi/placing"
	"github.com/KlyuchnikovV/webapi/validate"
)

// Example service.
type RequestAPI struct{}

func (api *RequestAPI) Prefix() string {
	return "request"
}

func (api *RequestAPI) Middlewares() []webapi.Middleware {
	return []webapi.Middleware{
		webapi.UseCORS(webapi.AllowedOrigins("*")),
	}
}

func (api *RequestAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"get": webapi.GET(
			api.GetByID,
			parameter.Integer("id", placing.InQuery,
				options.Description("ID of request."),
				validate.AND(validate.Greater(1), validate.Less(10)),
			),
		),
		"create": webapi.POST(
			api.Create,
			parameter.Description("Creates new request"),
			parameter.Body(new(entity.RequestBody),
				options.Description("Body description"),
			),
		),
		"create/sub-request": webapi.POST(
			api.CreateSubRequest,
			parameter.Body([]entity.RequestBody{}),
		),
		"filter": webapi.GET(
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

func (api *RequestAPI) Create(ctx webapi.Context) error {
	if body := ctx.Body(); body != nil {
		return ctx.OK(body)
	}

	return ctx.Created()
}

func (api *RequestAPI) CreateSubRequest(ctx webapi.Context) error {
	return ctx.Object(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", []entity.RequestBody{}),
	)
}

func (api *RequestAPI) GetByID(ctx webapi.Context) error {
	var id = ctx.Integer("id", placing.InQuery)

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(ctx webapi.Context) error {
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
