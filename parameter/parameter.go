package parameter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/engi/internal/middlewares"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/response"
)

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, place placing.Placing, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractParam(key, place, r, opts,
				func(request string) (interface{}, error) {
					return strconv.ParseBool(request)
				},
			)
		})
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, place placing.Placing, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractParam(key, place, r, opts,
				func(p string) (interface{}, error) {
					result, err := strconv.ParseInt(p, request.IntBase, request.BitSize)
					if err != nil {
						return nil, response.AsError(http.StatusBadRequest, "Parameter '%s' not of type int (got: '%s')", key, p)
					}

					return result, err
				},
			)
		})
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, place placing.Placing, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractParam(key, place, r, opts,
				func(p string) (interface{}, error) {
					result, err := strconv.ParseFloat(p, request.BitSize)
					if err != nil {
						return nil, response.AsError(http.StatusBadRequest, "Parameter '%s' not of type float (got: '%s')", key, p)
					}

					return result, err
				},
			)
		})
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, place placing.Placing, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractParam(key, place, r, opts,
				func(request string) (interface{}, error) {
					return request, nil
				},
			)
		})
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, place placing.Placing, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddParams(func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
			return request.ExtractParam(key, place, r, opts,
				func(request string) (interface{}, error) {
					result, err := time.Parse(layout, request)
					if err != nil {
						return nil, response.AsError(http.StatusBadRequest,
							"could not parse '%s' request to datetime using '%s' layout", key, layout,
						)
					}

					return result, err
				},
			)
		})
	}
}
