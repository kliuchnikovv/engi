package api

import (
	"github.com/labstack/echo"
)

type Context struct {
	Context echo.Context

	QueryParams QueryParams
	Response    Response

	bodyRequested bool
	body          interface{}

	requestedParams map[string]struct{}
	queryParameters map[string]interface{}
}

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
