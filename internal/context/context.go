package context

import (
	"context"
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type (
	Handler func(*Context) error

	// Context - provides methods for extracting data from query and response back.
	Context struct {
		context.Context // TODO: ignore
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
		Context:  r.Context(),
		Request:  request.New(r),
		Response: response.New(writer, responseMarshaler, responseObject),
	}
}
