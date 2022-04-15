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
}

type Body struct {
	String       string      `json:"field"`
	Integer      int         `json:"integer"`
	Array        []Body      `json:"array"`
	SimpleArray  []string    `json:"simple_array"`
	ArrayOfArray [][]float32 `json:"array_of_array"`
	WithoutTag   float64
}

func NewRequestAPI(engine *webapi.Engine) webapi.ServiceAPI {
	return &RequestAPI{
		Service: *webapi.NewService(engine, "request"),
	}
}

// Should contain only return statement for documentation.
func (api *RequestAPI) Routers() map[string]webapi.RouterByPath {
	return map[string]webapi.RouterByPath{
		"get": api.GET(
			api.GetByID,
			param.WithInteger("id",
				param.Description("ID of request."),
				param.AND(param.Greater(1), param.Less(10)),
			),
		),
		"create": api.POST(
			api.Create,
			param.WithBody(&Body{}),
		),
		"create/sub-request": api.POST(
			api.CreateSubRequest,
			param.WithBody([]Body{}),
		),
		"filter": api.GET(
			api.Filter,
			param.WithBool("bool"),
			param.WithFloat("float", param.NotEmpty),
			param.WithString("str",
				param.AND(param.NotEmpty, param.Greater(2)),
			),
			param.WithInteger("int"),
			param.WithTime("time",
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
		fmt.Sprintf("sub request created with body %#v", []Body{}),
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
