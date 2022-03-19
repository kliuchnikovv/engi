package webapi

import (
	"net/http"
)

type Binder interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// Context - provides methods for extracting data from query and response back.
type Context struct {
	*Request
	*Response
}

// NewContext - returns new Context instance from 'echo.Context'.
func NewContext(
	response http.ResponseWriter,
	request *http.Request,
	responseMarshaler ResponseMarshaler,
	responseObject Responser,
) *Context {
	return &Context{
		Request:  NewRequest(request),
		Response: NewResponse(response, responseMarshaler, responseObject),
	}
}

// PathParameter - retrieves path parameter by its name.
// func (ctx *Context) PathParameter(key string) string {
// 	// ctx.request.URL.Query()
// 	return ctx.Context.Param(key)
// }

// // AllPathParameters - return 'name - value' pairs of all path parameters.
// func (ctx *Context) AllPathParameters() map[string]string {
// 	var (
// 		names  = ctx.Context.ParamNames()
// 		result = make(map[string]string, len(names))
// 	)

// 	for _, name := range names {
// 		result[name] = ctx.Context.Param(name)
// 	}

// 	return result
// }

// // Body - returns query body.
// // Body must be requested by 'api.WithBody(pointer)'.
// func (ctx *Context) Body() interface{} {
// 	return ctx.body
// }
