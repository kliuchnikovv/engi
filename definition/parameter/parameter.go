package parameter

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
)

var numRegexp = regexp.MustCompile("[0-9]")

type Parameter struct {
	key     string
	placing placing.Placing
	options []request.Option

	typeName string
	regexp   string
	parse    func(string) (any, error)
}

func (parameter Parameter) Name() string {
	return parameter.key
}

func (parameter Parameter) Regexp() string {
	return parameter.regexp
}

// func (parameter Parameter) Bind(route *routes.Route) error {
// 	// if _, ok := route.Params[parameter.placing]; !ok {
// 	// 	route.Params[parameter.placing] = make(map[string]routes.Middleware)
// 	// }

// 	// route.Params[parameter.placing][parameter.key] = parameter

// 	return nil
// }

func (parameter Parameter) Handle(
	ctx context.Context,
	r *request.Request,
	response *response.Response,
) error {
	return request.ExtractParam(r,
		parameter.key,
		parameter.placing,
		parameter.options,
		parameter.parse,
	)
}

func (parameter Parameter) Docs(route *routes.Route) {
	panic("not implemented")
}

func (parameter Parameter) Priority() int {
	return 100
}

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, place placing.Placing, options ...request.Option) engi.Middleware {
	return &Parameter{
		key:      key,
		placing:  place,
		options:  options,
		typeName: "bool",
		regexp:   `(1|t|T|TRUE|true|True|0|f|F|FALSE|false|False)`,
		parse: func(request string) (any, error) {
			return strconv.ParseBool(request)
		},
	}
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, place placing.Placing, options ...request.Option) engi.Middleware {
	return &Parameter{
		key:      key,
		placing:  place,
		options:  options,
		typeName: "int64",
		regexp:   `((\+|-)?\d+)`,
		parse: func(p string) (interface{}, error) {
			result, err := strconv.ParseInt(p, request.IntBase, request.BitSize)
			if err != nil {
				return nil, err //response.BadRequest("Parameter '%s' not of type int (got: '%s')", key, p)
			}

			return result, nil
		},
	}
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, place placing.Placing, options ...request.Option) engi.Middleware {
	return &Parameter{
		key:      key,
		placing:  place,
		options:  options,
		typeName: "float64",
		regexp:   `((+|-)\d+(\.\d+)?)`,
		parse: func(p string) (interface{}, error) {
			result, err := strconv.ParseFloat(p, request.BitSize)
			if err != nil {
				return nil, err //response.BadRequest("Parameter '%s' not of type float (got: '%s')", key, p)
			}

			return result, nil
		},
	}
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, place placing.Placing, options ...request.Option) engi.Middleware {
	return &Parameter{
		key:      key,
		placing:  place,
		options:  options,
		typeName: "string",
		regexp:   `(.+)`,
		parse: func(p string) (interface{}, error) {
			return p, nil
		},
	}
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, place placing.Placing, options ...request.Option) engi.Middleware {
	return &Parameter{
		key:      key,
		placing:  place,
		options:  options,
		typeName: "time",
		regexp:   numRegexp.ReplaceAllString(layout, `\d`),
		parse: func(request string) (interface{}, error) {
			result, err := time.Parse(layout, request)
			if err != nil {
				return nil, err //response.BadRequest("could not parse '%s' request to datetime using '%s' layout", key, layout)
			}

			return result, nil
		},
	}
}
