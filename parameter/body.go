package parameter

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/middlewares"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/types"
	"github.com/KlyuchnikovV/engi/response"
)

// Body - takes pointer to structure and saves casted request body into context and pointer.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func Body(pointer interface{}, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			unmarshaler, err := request.GetUnmarshaler(r)
			if err != nil {
				return response.AsError(http.StatusInternalServerError, err.Error())
			}

			return request.ExtractBody(r, unmarshaler, pointer, opts)
		})
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func CustomBody(
	unmarshaler types.Unmarshaler,
	pointer interface{},
	opts ...request.Option,
) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractBody(r, unmarshaler, pointer, opts)
		})
	}
}
