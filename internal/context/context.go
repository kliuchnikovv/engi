package context

import (
	"net/http"

	"github.com/KlyuchnikovV/webapi/internal/request"
	"github.com/KlyuchnikovV/webapi/internal/response"
	"github.com/KlyuchnikovV/webapi/internal/types"
)

type (
	Handler func(*Context)

	// Context - provides methods for extracting data from query and response back.
	Context struct {
		*request.Request
		*response.Response
	}
)

func NewContext(
	writer http.ResponseWriter,
	r *http.Request,
	responseMarshaler types.Marshaler,
	responseObject types.Responser,
) *Context {
	return &Context{
		Request:  request.New(r),
		Response: response.New(writer, responseMarshaler, responseObject),
	}
}
