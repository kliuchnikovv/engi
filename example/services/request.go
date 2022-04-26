package services

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/entity"
	"github.com/KlyuchnikovV/webapi/param"
)

// Example service.
type RequestAPI struct {
	webapi.Service
}

func NewRequestAPI(engine *webapi.Engine) *RequestAPI {
	return &RequestAPI{
		Service: *webapi.NewService(engine, "request"),
	}
}

// Should contain only return statement for documentation.
func (api *RequestAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"get": api.GET(
			api.GetByID,
			param.QueryInteger("id",
				param.Description("ID of request."),
				param.AND(param.Greater(1), param.Less(10)),
			),
		),
		"create": api.POST(
			api.Create,
			param.Body(&entity.RequestBody{}),
		),
		"create/sub-request": api.POST(
			api.CreateSubRequest,
			param.Body([]entity.RequestBody{}),
		),
		"filter": api.GET(
			api.Filter,
			param.QueryBool("bool"),
			param.QueryFloat("float", param.NotEmpty),
			param.QueryString("str",
				param.AND(param.NotEmpty, param.Greater(2)),
			),
			param.QueryInteger("int"),
			param.QueryTime("time",
				"2006-01-02 15:04",
				param.Description("Filter by time field."),
			),
		),
	}
}

func (api *RequestAPI) Create(ctx *webapi.Context) error {
	if body := ctx.Body(); body != nil {
		return ctx.Response.OK(body)
	}

	ctx.Response.Created()

	return nil
}

func (api *RequestAPI) CreateSubRequest(ctx *webapi.Context) error {
	return ctx.Response.WithJSON(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", []entity.RequestBody{}),
	)
}

func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
	var id = ctx.Request.QueryInteger("id")

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.Response.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(ctx *webapi.Context) error {
	var (
		i     = ctx.Request.QueryInteger("int")
		str   = ctx.Request.QueryString("str")
		t     = ctx.Request.QueryTime("time", "2006-01-02 15:04")
		b     = ctx.Request.QueryBool("bool")
		float = ctx.Request.QueryFloat("float")
	)

	return ctx.Response.OK(fmt.Sprintf(
		"filtered by id: '%d' and field: %s, time: %s, isAssignable: %t and float: %f",
		i, str, t.Format("15:04 02/01/2006"), b, float,
	))
}
