package param

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/webapi/options"
	"github.com/KlyuchnikovV/webapi/types"
)

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, place options.Placing, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractParam(key, place, request, opts, func(options string) (interface{}, error) {
			return strconv.ParseBool(options)
		})
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, place options.Placing, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractParam(key, place, request, opts, func(p string) (interface{}, error) {
			result, err := strconv.ParseInt(p, options.IntBase, options.BitSize)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest, "Parameter '%s' not of type int", key)
			}

			return result, err
		})
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, place options.Placing, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractParam(key, place, request, opts, func(p string) (interface{}, error) {
			result, err := strconv.ParseFloat(p, options.BitSize)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest, "Parameter '%s' not of type float", key)
			}

			return result, err
		})
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, place options.Placing, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractParam(key, place, request, opts, func(options string) (interface{}, error) {
			return options, nil
		})
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, place options.Placing, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractParam(key, place, request, opts, func(options string) (interface{}, error) {
			result, err := time.Parse(layout, options)
			if err != nil {
				return nil, types.NewErrorResponse(http.StatusBadRequest,
					"could not parse '%s' options to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		})
	}
}

// Body - takes pointer to structure and saves casted request body into context and pointer.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func Body(pointer interface{}, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		unmarshal, err := options.GetUnmarshaler(request)
		if err != nil {
			return err
		}

		return options.ExtractBody(request, unmarshal, pointer, opts)
	}
}

// CustomBody - takes unmarshaler and pointer to structure and saves casted request body into context.
//
// Result can be retrieved from context using 'context.QueryParams.Body'.
func CustomBody(unmarshal types.Unmarshaler, pointer interface{}, opts ...options.Option) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		return options.ExtractBody(request, unmarshal, pointer, opts)
	}
}

func Description(desc string) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		request.Description = desc
		return nil
	}
}

func Header(key string) options.HandlerParams {
	return func(request *options.Request, _ http.ResponseWriter) error {
		header, ok := request.Request().Header[key]
		if !ok || len(header) == 0 {
			return fmt.Errorf("no '%s' header found", key)
		}

		return nil
	}
}
