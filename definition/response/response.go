package response

import (
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type (
	Responser types.Responser
	Marshaler types.Marshaler
)

type responserObject struct {
	responser types.Responser
}

func (object *responserObject) Bind(route *routes.Route) error {
	route.Responser = object.responser
	return nil
}

func (object *responserObject) Handle(*request.Request, *response.Response) error {
	return nil
}

func (object *responserObject) Docs(*routes.Route) {
	panic("not implemented")
}

type marshalerObject struct {
	marshaler types.Marshaler
}

func (object *marshalerObject) Bind(route *routes.Route) error {
	route.Marshaler = object.marshaler
	return nil
}

func (object *marshalerObject) Handle(*request.Request, *response.Response) error {
	return nil
}

func (object *marshalerObject) Docs(*routes.Route) {
	panic("not implemented")
}
