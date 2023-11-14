package parameter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/response"
)

type parameter struct {
	key     string
	place   placing.Placing
	options []request.Option
	extract func(string) (interface{}, error)
}

func (p *parameter) Handle(r *request.Request, _ http.ResponseWriter) *response.AsObject {
	if err := request.ExtractParam(p.key, p.place, r, p.options, p.extract); err != nil {
		return response.AsError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, place placing.Placing, opts ...request.Option) request.Middleware {
	return &parameter{
		key:     key,
		options: opts,
		place:   place,
		extract: func(request string) (interface{}, error) {
			return strconv.ParseBool(request)
		},
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, place placing.Placing, opts ...request.Option) request.Middleware {
	return &parameter{
		key:     key,
		options: opts,
		place:   place,
		extract: func(p string) (interface{}, error) {
			result, err := strconv.ParseInt(p, request.IntBase, request.BitSize)
			if err != nil {
				return nil, response.AsError(http.StatusBadRequest, "Parameter '%s' not of type int (got: '%s')", key, p)
			}

			return result, err
		},
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, place placing.Placing, opts ...request.Option) request.Middleware {
	return &parameter{
		key:     key,
		options: opts,
		place:   place,
		extract: func(p string) (interface{}, error) {
			result, err := strconv.ParseFloat(p, request.BitSize)
			if err != nil {
				return nil, response.AsError(http.StatusBadRequest, "Parameter '%s' not of type float (got: '%s')", key, p)
			}

			return result, err
		},
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, place placing.Placing, opts ...request.Option) request.Middleware {
	return &parameter{
		key:     key,
		options: opts,
		place:   place,
		extract: func(request string) (interface{}, error) {
			return request, nil
		},
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, place placing.Placing, opts ...request.Option) request.Middleware {
	return &parameter{
		key:     key,
		options: opts,
		place:   place,
		extract: func(request string) (interface{}, error) {
			result, err := time.Parse(layout, request)
			if err != nil {
				return nil, response.AsError(http.StatusBadRequest,
					"could not parse '%s' request to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		},
	}
}

// func Description(desc string)

// func Description(desc string) request.Middleware {
// 	return func(request *request.Request, _ http.ResponseWriter) error {
// 		request.Description = desc
// 		return nil
// 	}
// }

// func Header(key string) request.Middleware {
// 	return func(request *request.Request, _ http.ResponseWriter) error {
// 		header, ok := request.GetRequest().Header[key]
// 		if !ok || len(header) == 0 {
// 			return fmt.Errorf("no '%s' header found", key)
// 		}

// 		return nil
// 	}
// }
