package parameter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/types"
	"github.com/KlyuchnikovV/engi/placing"
	"github.com/KlyuchnikovV/engi/response"
)

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, place placing.Placing, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractParam(key, place, r, opts, func(request string) (interface{}, error) {
			return strconv.ParseBool(request)
		})
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, place placing.Placing, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractParam(key, place, r, opts, func(p string) (interface{}, error) {
			result, err := strconv.ParseInt(p, request.IntBase, request.BitSize)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest, "Parameter '%s' not of type int (got: '%s')", key, p)
			}

			return result, err
		})
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, place placing.Placing, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractParam(key, place, r, opts, func(p string) (interface{}, error) {
			result, err := strconv.ParseFloat(p, request.BitSize)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest, "Parameter '%s' not of type float (got: '%s')", key, p)
			}

			return result, err
		})
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, place placing.Placing, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractParam(key, place, r, opts, func(request string) (interface{}, error) {
			return request, nil
		})
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, place placing.Placing, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractParam(key, place, r, opts, func(request string) (interface{}, error) {
			result, err := time.Parse(layout, request)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest,
					"could not parse '%s' request to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		})
	}
}

// Body - takes pointer to structure and saves casted request body into context and pointer.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func Body(pointer interface{}, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		unmarshal, err := request.GetUnmarshaler(r)
		if err != nil {
			return err
		}

		return request.ExtractBody(r, unmarshal, pointer, opts)
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func CustomBody(unmarshal types.Unmarshaler, pointer interface{}, opts ...request.Option) request.HandlerParams {
	return func(r *request.Request, _ http.ResponseWriter) error {
		return request.ExtractBody(r, unmarshal, pointer, opts)
	}
}

func Description(desc string) request.HandlerParams {
	return func(request *request.Request, _ http.ResponseWriter) error {
		request.Description = desc
		return nil
	}
}

func Header(key string) request.HandlerParams {
	return func(request *request.Request, _ http.ResponseWriter) error {
		header, ok := request.GetRequest().Header[key]
		if !ok || len(header) == 0 {
			return fmt.Errorf("no '%s' header found", key)
		}

		return nil
	}
}
