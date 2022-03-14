package api

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
