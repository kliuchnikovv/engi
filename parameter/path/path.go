package path

import (
	"github.com/KlyuchnikovV/engi/internal/middlewares"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter"
	"github.com/KlyuchnikovV/engi/parameter/placing"
)

// Bool - mandatory boolean Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Bool'.
func Bool(key string, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return parameter.Bool(key, placing.InPath, opts...)
}

// Integer - queries mandatory integer Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Integer'.
func Integer(key string, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return parameter.Integer(key, placing.InPath, opts...)
}

// Float - mandatory floating point number Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.Float'.
func Float(key string, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return parameter.Float(key, placing.InPath, opts...)
}

// String - mandatory string Parameter from request by 'key'.
//
// Result can be retrieved from context using 'context.QueryParams.String'.
func String(key string, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return parameter.String(key, placing.InPath, opts...)
}

// Time - mandatory time Parameter from request by 'key' using 'layout'.
//
// Result can be retrieved from context using 'context.QueryParams.Time'.
func Time(key, layout string, opts ...request.Option) func(middlewares *middlewares.Middlewares) {
	return parameter.Time(key, layout, placing.InPath, opts...)
}
