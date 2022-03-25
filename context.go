package webapi

import (
	"net/http"

	"github.com/KlyuchnikovV/webapi/param"
	"github.com/KlyuchnikovV/webapi/types"
)

type (
	RouterByPath func(string)
	Route        func(*Context) error
	Handler      func(*Context)

	// Context - provides methods for extracting data from query and response back.
	Context struct {
		*param.Request
		*param.Response
	}
)

func NewContext(
	response http.ResponseWriter,
	request *http.Request,
	responseMarshaler types.Marshaler,
	responseObject types.Responser,
) *Context {
	return &Context{
		Request:  param.NewRequest(request),
		Response: param.NewResponse(response, responseMarshaler, responseObject),
	}
}
