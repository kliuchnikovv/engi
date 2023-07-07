package webapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/webapi/internal"
	"github.com/KlyuchnikovV/webapi/internal/context"
	"github.com/KlyuchnikovV/webapi/internal/request"
	"github.com/KlyuchnikovV/webapi/internal/types"
	"github.com/KlyuchnikovV/webapi/logger"
	"github.com/KlyuchnikovV/webapi/response"
)

type (
	ServiceAPI interface {
		// Prefix - prefix of all paths for this service.
		Prefix() string

		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		Routers() map[string]RouterByPath
	}

	MiddlewaresAPI interface {
		Middlewares() []Middleware
	}

	// Service - provides basic service methods.
	Service struct {
		middlewares []request.HandlerParams

		handlers map[string]internal.HanlderNode

		marshaler types.Marshaler
		responser types.Responser

		log *logger.Log

		api  ServiceAPI
		path string
	}

	RouterByPath func(*Service, string)
	Middleware   func(*Service)
)

func NewService(engine *Engine, api ServiceAPI, path string) *Service {
	return &Service{
		handlers: make(map[string]internal.HanlderNode),

		log:       engine.log,
		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,

		api:  api,
		path: path,
	}
}

// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
func (api *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = context.NewContext(w, r, api.marshaler, api.responser)

	for _, middleware := range api.middlewares {
		err := middleware(ctx.Request, ctx.Response.GetResponse())
		if err == nil {
			continue
		}

		var response response.AsObject

		errors.As(err, &response)

		switch response.Code {
		case http.StatusOK:
			err = ctx.Response.OK(response.Code)
		default:
			err = ctx.Response.Error(response.Code, response.ErrorString)
		}

		if err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	if _, ok := api.handlers[r.Method]; !ok {
		if err := ctx.Response.NotFound("method '%s' not appliable for '%s'", r.Method, r.URL.Path); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	if !api.handlers[r.Method].Handle(r.URL.Path, ctx) {
		if err := ctx.Response.NotFound("path '%s' not found for method '%s'", r.URL.Path, r.Method); err != nil {
			api.log.SendErrorf(err.Error())
		}
	}
}

// add - creates route with custom method and path.
func (api *Service) add(
	method, path string,
	route Route,
	middlewares ...request.HandlerParams,
) {
	if _, ok := api.handlers[method]; !ok {
		api.handlers[method] = internal.NewStringHandler("", nil)
	}

	api.handlers[method].Add(
		api.handle(route, middlewares...),
		strings.Split(path, "/")...,
	)
}

func (api *Service) handle(route Route, middlewares ...request.HandlerParams) context.Handler {
	return func(ctx *context.Context) {
		var response response.AsObject

		for _, middleware := range middlewares {
			if err := middleware(ctx.Request, ctx.Response.GetResponse()); err != nil {
				errors.As(err, &response)

				if err := ctx.Response.Error(response.Code, response.ErrorString); err != nil {
					api.log.SendErrorf(err.Error())
				}

				return
			}
		}

		if err := route(ctx); err != nil {
			if err := ctx.Response.InternalServerError(err.Error()); err != nil {
				api.log.SendErrorf(err.Error())
			}
		}
	}
}

// GET - implements GET api method call.
func GET(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodGet, path, route, middlewares...)
	}
}

// PUT - implements PUT api method call.
func PUT(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPut, path, route, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func HEAD(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodHead, path, route, middlewares...)
	}
}

// POST - implements POST api method call.
func POST(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPost, path, route, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func PATCH(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPatch, path, route, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func TRACE(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodTrace, path, route, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func DELETE(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodDelete, path, route, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func CONNECT(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodConnect, path, route, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func OPTIONS(route Route, middlewares ...request.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodOptions, path, route, middlewares...)
	}
}
