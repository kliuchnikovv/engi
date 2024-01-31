package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/example/entity"
	"github.com/KlyuchnikovV/engi/parameter"
	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/parameter/query"
	"github.com/KlyuchnikovV/engi/validate"
)

// Example service.
type RequestAPI struct{}

func (api *RequestAPI) Prefix() string {
	return "request"
}

func (api *RequestAPI) Middlewares() []engi.Register {
	return []engi.Register{
		engi.UseCORS(engi.AllowedOrigins("*")),
		engi.UseAuthorization(engi.BasicAuth("Dave", "IsCrazy")),
	}
}

func (api *RequestAPI) Routers() engi.Routes {
	return engi.Routes{
		"get": engi.GET(api.GetByID,
			query.Integer("id",
				validate.AND(validate.Greater(1), validate.Less(10)),
			),
		),
		"create": engi.POST(api.Create,
			parameter.Body(new(entity.RequestBody)),
		),
		"create/sub-request": engi.POST(api.CreateSubRequest,
			parameter.Body([]entity.RequestBody{}),
		),
		"filter": engi.GET(api.Filter,
			query.Bool("bool"),
			query.Float("float", validate.NotEmpty),
			query.Integer("int"),
			query.String("str",
				validate.AND(validate.NotEmpty, validate.Greater(2)),
			),
			query.Time("time", "2006-01-02 15:04"),
		),
	}
}

func (api *RequestAPI) Create(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	if body := request.Body(); body != nil {
		return response.OK(body)
	}

	return response.Created()
}

func (api *RequestAPI) CreateSubRequest(
	_ context.Context,
	_ engi.Request,
	response engi.Response,
) error {
	return response.Object(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", []entity.RequestBody{}),
	)
}

func (api *RequestAPI) GetByID(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var id = request.Integer("id", placing.InQuery)

	// Do something with id (we will check it)
	if id < 0 {
		return response.BadRequest("id can't be negative (got: %d)", id)
	}

	return response.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(
	_ context.Context,
	request engi.Request,
	response engi.Response,
) error {
	var (
		i     = request.Integer("int", placing.InQuery)
		str   = request.String("str", placing.InQuery)
		t     = request.Time("time", "2006-01-02 15:04", placing.InQuery)
		b     = request.Bool("bool", placing.InQuery)
		float = request.Float("float", placing.InQuery)
	)

	return response.OK(fmt.Sprintf(
		"filtered by id: '%d' and field: %s, time: %s, isAssignable: %t and float: %f",
		i, str, t.Format("15:04 02/01/2006"), b, float,
	))
}
