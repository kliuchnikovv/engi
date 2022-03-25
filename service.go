package webapi

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/webapi/param"
	"github.com/KlyuchnikovV/webapi/types"
)

type (
	ServiceAPI interface {
		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		//	Must be implemented by Service
		Routers() map[string]RouterByPath

		// PathPrefix - prefix of all paths for this service.
		PathPrefix() string

		// ServeHTTP should write reply headers and data to the ResponseWriter
		// and then return. Returning signals that the request is finished; it
		// is not valid to use the ResponseWriter or read from the
		// Request.Body after or concurrently with the completion of the
		// ServeHTTP call.
		ServeHTTP(http.ResponseWriter, *http.Request)
	}

	// Service - provides basic service methods.
	Service struct {
		prefix string

		middlewares []param.HandlersOption

		handlers map[string]map[string]Handler

		marshaler types.Marshaler
		responser types.Responser

		log *types.Log
	}
)

func NewService(engine *Engine, prefix string) *Service {
	return &Service{
		prefix:   strings.Trim(prefix, "/"),
		handlers: make(map[string]map[string]Handler),

		log:       engine.log,
		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,
	}
}

// PathPrefix - prefix of all paths for this service.
func (api *Service) PathPrefix() string {
	return api.prefix
}

// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
func (api *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, path: %s\n", r.Method, r.URL.Path)

	var ctx = NewContext(w, r, api.marshaler, api.responser)

	if _, ok := api.handlers[r.Method]; !ok {
		if err := ctx.Response.NotFound("method '%s' not appliable for '%s'", r.Method, r.URL.Path); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	handler, ok := api.handlers[r.Method][r.URL.Path]
	if !ok {
		if err := ctx.Response.NotFound("path '%s' not found for method '%s'", r.URL.Path, r.Method); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	handler(ctx)
}

// Add - creates route with custom method and path.
func (api *Service) Add(
	method, path string,
	route Route,
	middlewares ...param.HandlersOption,
) {
	if _, ok := api.handlers[method]; !ok {
		api.handlers[method] = make(map[string]Handler)
	}

	if _, ok := api.handlers[method][path]; ok {
		api.log.SendErrorf("method '%s' with path '%s' already defined", method, path)
		return
	}

	api.handlers[method][path] = api.handle(route, middlewares...)
}

func (api *Service) handle(route Route, middlewares ...param.HandlersOption) Handler {
	return func(ctx *Context) {
		var response types.ResponseObject

		for _, middleware := range api.middlewares {
			if err := middleware(ctx.Request); err != nil {
				errors.As(err, &response)

				if err := ctx.Response.Error(response.Code, response.ErrorString); err != nil {
					api.log.SendErrorf(err.Error())
				}

				return
			}
		}

		for _, middleware := range middlewares {
			if err := middleware(ctx.Request); err != nil {
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
func (api *Service) GET(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodGet, path, route, middlewares...)
	}
}

// PUT - implements PUT api method call.
func (api *Service) PUT(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodPut, path, route, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func (api *Service) HEAD(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodHead, path, route, middlewares...)
	}
}

// POST - implements POST api method call.
func (api *Service) POST(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodPost, path, route, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func (api *Service) PATCH(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodPatch, path, route, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func (api *Service) TRACE(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodTrace, path, route, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func (api *Service) DELETE(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodDelete, path, route, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func (api *Service) CONNECT(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodConnect, path, route, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func (api *Service) OPTIONS(route Route, middlewares ...param.HandlersOption) RouterByPath {
	return func(path string) {
		api.Add(http.MethodOptions, path, route, middlewares...)
	}
}
