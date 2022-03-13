package main

import (
	"fmt"

	webapi "github.com/KlyuchnikovV/webapi/api"
)

// Example service
type RequestAPI struct {
	webapi.API

	Body Body
}

type Body struct {
	Field string `json:"field"`
}

func NewRequestAPI() webapi.API {
	return &RequestAPI{
		API: webapi.New("request"),
	}
}

func (api *RequestAPI) Routers() map[string]webapi.RouterFunc {
	return map[string]webapi.RouterFunc{
		":id": api.GET(
			api.GetByID,
		),
		"create": api.POST(
			api.Create,
			api.WithBody(&Body{}),
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
	if body, _ := ctx.QueryParams.Body(); body != nil {
		api.Body = *body.(*Body)
	}

	return ctx.Response.Created()
}

func (api *RequestAPI) GetByID(ctx *webapi.Context) error {
	return ctx.Response.OK(api.Body)
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
