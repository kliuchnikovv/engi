package service

import (
	"fmt"
	"net/http"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/param"
)

// Example service.
type RequestAPI struct {
	webapi.Service

	// Bad practice in case of concurrency,
	// but useful in this example.
	Request    Body
	SubRequest Body
}

type Body struct {
	Field string `json:"field"`
}

func NewRequestAPI(engine *webapi.Engine) webapi.ServiceAPI {
	return &RequestAPI{
		Service: *webapi.NewService(engine, "request"),
	}
}

func (api *RequestAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"get": api.GET(
			api.GetByID,
			param.WithInteger("id",
				param.Description("id of request"),
				param.AND(param.Greater(1), param.Less(10)),
			),
		),
		"create": api.POST(
			api.Create,
			param.WithBody(&Body{}),
		),
		"create/sub-request": api.POST(
			api.CreateSubRequest,
			param.WithBody(&api.SubRequest),
		),
		"filter": api.GET(
			api.Filter,
			param.WithBool("bool"),
			param.WithFloat("float"),
			param.WithString("str"),
			param.WithInteger("int"),
			param.WithTime("time", "2006-01-02 15:04"),
		),
	}
}

func (api *RequestAPI) Create(ctx *webapi.Context) error {
	if body := ctx.Body(); body != nil {
		api.Request = *body.(*Body)
	}

	ctx.Response.Created()

	return nil
}

func (api *RequestAPI) CreateSubRequest(ctx *webapi.Context) error {
	return ctx.Response.WithJSON(http.StatusCreated,
		fmt.Sprintf("sub request created with body %#v", api.SubRequest),
	)
}

func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
	var id = ctx.Request.Integer("id")

	// Do something with id (we will check it)
	if id < 0 {
		return ctx.Response.BadRequest("id can't be negative (got: %d)", id)
	}

	return ctx.Response.OK(fmt.Sprintf("got id: '%d'", id))
}

func (api *RequestAPI) Filter(ctx *webapi.Context) error {
	var (
		i     = ctx.Request.Integer("int")
		str   = ctx.Request.String("str")
		t     = ctx.Request.Time("time", "2006-01-02 15:04")
		b     = ctx.Request.Bool("bool")
		float = ctx.Request.Float("float")
	)

	return ctx.Response.OK(fmt.Sprintf(
		"filtered by id: '%d' and field: %s, time: %s, isAssignable: %t and float: %f",
		i, str, t.Format("15:04 02/01/2006"), b, float,
	))
}
