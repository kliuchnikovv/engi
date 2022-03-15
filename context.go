package webapi

import (
	"github.com/labstack/echo"
)

// Context - provides methods for extracting data from query and response back.
type Context struct {
	Context echo.Context

	QueryParams QueryParams
	Response    Response

	bodyRequested bool
	body          interface{}

	requestedParams map[string]struct{}
	queryParameters map[string]interface{}
}

// NewContext - returns new Context instance from 'echo.Context'.
func NewContext(ctx echo.Context) *Context {
	var context = &Context{
		Context:         ctx,
		requestedParams: map[string]struct{}{},
		queryParameters: map[string]interface{}{},
	}

	context.QueryParams = NewQueryParams(context)
	context.Response = NewResponse(context)

	return context
}

// PathParameter - retrieves path parameter by its name.
func (ctx *Context) PathParameter(key string) string {
	return ctx.Context.Param(key)
}

// AllPathParameters - return 'name - value' pairs of all path parameters.
func (ctx *Context) AllPathParameters() map[string]string {
	var (
		names  = ctx.Context.ParamNames()
		result = make(map[string]string, len(names))
	)

	for _, name := range names {
		result[name] = ctx.Context.Param(name)
	}

	return result
}

// Body - returns query body.
// Body must be requested by 'api.WithBody(pointer)'.
func (ctx *Context) Body() interface{} {
	return ctx.body
}
