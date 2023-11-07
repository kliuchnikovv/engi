package query

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/placing"
	"github.com/KlyuchnikovV/engi/response"
)

type queryParameter struct {
	key     string
	options []request.Option
	extract func(string) (interface{}, error)
}

func (p *queryParameter) Handle(r *request.Request, w http.ResponseWriter) error {
	return request.ExtractParam(p.key, placing.InQuery, r, p.options, p.extract)
}

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, opts ...request.Option) request.Middleware {
	return &queryParameter{
		key:     key,
		options: opts,
		extract: func(request string) (interface{}, error) {
			return strconv.ParseBool(request)
		},
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, opts ...request.Option) request.Middleware {
	return &queryParameter{
		key:     key,
		options: opts,
		extract: func(p string) (interface{}, error) {
			result, err := strconv.ParseInt(p, request.IntBase, request.BitSize)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest, "Parameter '%s' not of type int (got: '%s')", key, p)
			}

			return result, err
		},
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, opts ...request.Option) request.Middleware {
	return &queryParameter{
		key:     key,
		options: opts,
		extract: func(p string) (interface{}, error) {
			result, err := strconv.ParseFloat(p, request.BitSize)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest, "Parameter '%s' not of type float (got: '%s')", key, p)
			}

			return result, err
		},
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, opts ...request.Option) request.Middleware {
	return &queryParameter{
		key:     key,
		options: opts,
		extract: func(request string) (interface{}, error) {
			return request, nil
		},
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, opts ...request.Option) request.Middleware {
	return &queryParameter{
		key:     key,
		options: opts,
		extract: func(request string) (interface{}, error) {
			result, err := time.Parse(layout, request)
			if err != nil {
				return nil, response.NewError(http.StatusBadRequest,
					"could not parse '%s' request to datetime using '%s' layout", key, layout,
				)
			}

			return result, err
		},
	}
}
