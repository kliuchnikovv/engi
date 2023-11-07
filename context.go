package engi

import (
	"context"
	"net/http"
	"time"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
)

type Context interface {
	// Basic go context

	context.Context

	// Request.

	// Headers - returns request headers.
	Headers() map[string][]string
	// All - returns all parsed parameters.
	All() map[placing.Placing]map[string]string
	// GetParameter - returns parameter value from defined place.
	GetParameter(string, placing.Placing) string
	// GetRequest - return http.Request object associated with request.
	GetRequest() *http.Request
	// Body - returns request body.
	// Body must be requested by 'api.Body(pointer)' or 'api.CustomBody(unmarshaler, pointer)'.
	Body() interface{}
	// Bool - returns boolean parameter.
	// Mandatory parameter should be requested by 'api.Bool'.
	// Otherwise, parameter will be obtained by key and its value will be checked for truth.
	Bool(string, placing.Placing) bool
	// Integer - returns integer parameter.
	// Mandatory parameter should be requested by 'api.Integer'.
	// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
	Integer(string, placing.Placing) int64
	// Float - returns floating point number parameter.
	// Mandatory parameter should be requested by 'api.Float'.
	// Otherwise, parameter will be obtained by key and its value will be converted to float64.
	Float(string, placing.Placing) float64
	// String - returns String parameter.
	// Mandatory parameter should be requested by 'api.String'.
	// Otherwise, parameter will be obtained by key.
	String(string, placing.Placing) string
	// Time - returns date-time parameter.
	// Mandatory parameter should be requested by 'api.Time'.
	// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
	Time(key string, layout string, paramPlacing placing.Placing) time.Time

	// Responses.

	// ResponseWriter - returns http.ResponseWriter associated with request.
	ResponseWriter() http.ResponseWriter
	// Object - responses with provided custom code and body.
	// Body will be marshaled using service-defined object and marshaler.
	Object(code int, payload interface{}) error
	// WithourContent - responses with provided custom code and no body.
	WithoutContent(code int) error
	// Error - responses custom error with provided code and formatted string message.
	Error(code int, format string, args ...interface{}) error
	// OK - writes payload into json's 'result' field with 200 http code.
	OK(payload interface{}) error
	// Created - responses with 201 http code and no content.
	Created() error
	// NoContent - responses with 204 http code and no content.
	NoContent() error
	// BadRequest - responses with 400 code and provided formatted string message.
	BadRequest(format string, args ...interface{}) error
	// Forbidden - responses with 403 error code and provided formatted string message.
	Forbidden(format string, args ...interface{}) error
	// NotFound - responses with 404 error code and provided formatted string message.
	NotFound(format string, args ...interface{}) error
	// MethodNotAllowed - responses with 405 error code and provided formatted string message.
	MethodNotAllowed(format string, args ...interface{}) error
	// InternalServerError - responses with 500 error code and provided formatted string message.
	InternalServerError(format string, args ...interface{}) error
}

type Route func(Context) error

// GET - implements GET api method call.
func GET(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodGet, path, route, middlewares...)
	}
}

// PUT - implements PUT api method call.
func PUT(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPut, path, route, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func HEAD(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodHead, path, route, middlewares...)
	}
}

// POST - implements POST api method call.
func POST(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPost, path, route, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func PATCH(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPatch, path, route, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func TRACE(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodTrace, path, route, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func DELETE(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodDelete, path, route, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func CONNECT(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodConnect, path, route, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func OPTIONS(route Route, middlewares ...request.Middleware) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodOptions, path, route, middlewares...)
	}
}
