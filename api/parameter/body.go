package parameter

import (
	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type BodyParameter struct {
	pointer interface{}

	unmarshaler types.Unmarshaler
	options     []request.Option
}

func (body *BodyParameter) Bind(route *routes.Route) error {
	route.Body = body

	return nil
}

func (body *BodyParameter) Handle(r *request.Request, response *response.Response) error {
	var (
		err         error
		unmarshaler = body.unmarshaler
	)

	if unmarshaler == nil {
		unmarshaler, err = request.GetUnmarshaler(r)
		if err != nil {
			return response.InternalServerError(err.Error())
		}
	}

	return request.ExtractBody(r, unmarshaler, body.pointer, body.options)
}

func (body *BodyParameter) Docs(route *routes.Route) {
	panic("not implemented")
}

// Body - takes pointer to structure and saves casted request body into context and pointer.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func Body(pointer interface{}, options ...request.Option) engi.Middleware {
	return &BodyParameter{
		pointer: pointer,
		options: options,
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func CustomBody(
	unmarshaler types.Unmarshaler,
	pointer interface{},
	options ...request.Option,
) engi.Middleware {
	return &BodyParameter{
		pointer:     pointer,
		unmarshaler: unmarshaler,
		options:     options,
	}
}
