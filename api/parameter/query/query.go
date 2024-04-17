package query

import (
	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/api/parameter"
	"github.com/KlyuchnikovV/engi/api/parameter/placing"
	"github.com/KlyuchnikovV/engi/internal/request"
)

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, opts ...request.Option) engi.Middleware {
	return parameter.Bool(key, placing.InQuery, opts...)
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, opts ...request.Option) engi.Middleware {
	return parameter.Integer(key, placing.InQuery, opts...)
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, opts ...request.Option) engi.Middleware {
	return parameter.Float(key, placing.InQuery, opts...)
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, opts ...request.Option) engi.Middleware {
	return parameter.String(key, placing.InQuery, opts...)
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, opts ...request.Option) engi.Middleware {
	return parameter.Time(key, layout, placing.InQuery, opts...)
}
