package response

import (
	"context"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
	"github.com/kliuchnikovv/engi/internal/types"
)

type (
	Responser types.Responser
	Marshaler types.Marshaler
)

type responserObject struct {
	responser types.Responser
}

func (object *responserObject) Bind(route *routes.Route) error {
	return nil
}

func (object *responserObject) Handle(ctx context.Context, _ *request.Request, resp *response.Response) error {
	response.SetResponser(resp, object.responser)
	return nil
}

func (object *responserObject) Docs(*routes.Route) {
	panic("not implemented")
}

func (object *responserObject) Priority() int {
	return 0
}

type marshalerObject struct {
	marshaler types.Marshaler
}

func (object *marshalerObject) Bind(route *routes.Route) error {
	return nil
}

func (object *marshalerObject) Handle(ctx context.Context, _ *request.Request, resp *response.Response) error {
	response.SetMarshaler(resp, object.marshaler)
	return nil
}

func (object *marshalerObject) Docs(*routes.Route) {
	panic("not implemented")
}

func (object *marshalerObject) Priority() int {
	return 0
}
