package webapi

import (
	"net/http"

	"github.com/KlyuchnikovV/webapi/options"
	"github.com/KlyuchnikovV/webapi/types"
)

type (
	RouterByPath func(*Service, string)
	Route        func(*Context) error
	Handler      func(*Context)
	Middleware   func(*Service)

	// Context - provides methods for extracting data from query and response back.
	Context struct {
		*options.Request
		*options.Response
	}
)

func NewContext(
	response http.ResponseWriter,
	request *http.Request,
	responseMarshaler types.Marshaler,
	responseObject types.Responser,
) *Context {
	return &Context{
		Request:  options.NewRequest(request),
		Response: options.NewResponse(response, responseMarshaler, responseObject),
	}
}
