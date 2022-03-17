package service

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi"
)

// Example service.
type RequestAPI struct {
	webapi.ServiceAPI

	// Bad practice in case of concurrency,
	// but useful in this example.
	Request    Body
	SubRequest Body
}

type Body struct {
	Field string `json:"field"`
}

func NewRequestAPI() webapi.ServiceAPI {
	return &RequestAPI{
		ServiceAPI: webapi.NewService("request"),
	}
}

func (api *RequestAPI) Routers() map[string]webapi.RouterFunc {
	return map[string]webapi.RouterFunc{
		"get": api.GET(
			api.GetByID,
			api.WithInt("id"),
		),
		"create": api.POST(
			api.Create,
			api.WithBody(&Body{}),
		),
		"create/sub-request": api.POST(
			api.CreateSubRequest,
			api.WithBody(&api.SubRequest),
		),
		"filter": api.GET(
			api.Filter,
			api.WithBool("bool"),
			api.WithInt("int"),
			api.WithFloat("float"),
			api.WithString("str"),
			api.WithTime("time", "2006-01-02 15:04"),
		),
	}
}

func (api *RequestAPI) Create(ctx *webapi.Context) error {
	if body := ctx.QueryParams.Body(); body != nil {
		api.Request = *body.(*Body)
	}

	return ctx.Response.Created()
}

func (api *RequestAPI) CreateSubRequest(ctx *webapi.Context) error {
	return ctx.Response.WithJSON(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", api.SubRequest),
	)
}

func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
	var id = ctx.QueryParams.Integer("id")

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.Response.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(ctx *webapi.Context) error {
	var (
		i     = ctx.QueryParams.Integer("int")
		str   = ctx.QueryParams.String("str")
		t     = ctx.QueryParams.Time("time", "2006-01-02 15:04")
		b     = ctx.QueryParams.Bool("bool")
		float = ctx.QueryParams.Float("float")
	)

	return ctx.Response.OK(fmt.Sprintf(
		"filtered by id: '%d' and field: %s, time: %s, isAssignable: %t and float: %f",
		i, str, t.Format("15:04 02/01/2006"), b, float,
	))
}
