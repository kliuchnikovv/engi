package webapi

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type (
	ServiceAPI interface {
		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		//	Must be implemented by Service
		Routers() map[string]RouterFunc

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

		middlewares []HandlerFunc

		routes map[string]map[string]ResultFunc

		marshaler MarshalerFunc
		responser Responser

		log *Log
	}
)

func NewService(engine *Engine, prefix string) *Service {
	return &Service{
		prefix: strings.Trim(prefix, "/"),
		routes: make(map[string]map[string]ResultFunc),

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

	if _, ok := api.routes[r.Method]; !ok {
		if err := ctx.NotFound("method '%s' not appliable for '%s'", r.Method, r.URL.Path); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	route, ok := api.routes[r.Method][r.URL.Path]
	if !ok {
		if err := ctx.NotFound("path '%s' not found for method '%s'", r.URL.Path, r.Method); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	route(ctx)
}

// Add - creates route with custom method and path.
func (api *Service) Add(
	method, path string,
	handler HandlerFunc,
	middlewares ...HandlerFunc,
) {
	if _, ok := api.routes[method]; !ok {
		api.routes[method] = make(map[string]ResultFunc)
	}

	if _, ok := api.routes[method][path]; ok {
		api.log.SendErrorf("method '%s' with path '%s' already defined", method, path)
		return
	}

	api.routes[method][path] = api.route(handler, middlewares...)
}

func (api *Service) route(handler HandlerFunc, middlewares ...HandlerFunc) ResultFunc {
	return func(ctx *Context) {
		var response ResponseObject

		for _, middleware := range api.middlewares {
			if err := middleware(ctx); err != nil {
				errors.As(err, &response)

				if err := ctx.Error(response.code, response.ErrorString); err != nil {
					api.log.SendErrorf(err.Error())
				}

				return
			}
		}

		for _, middleware := range middlewares {
			if err := middleware(ctx); err != nil {
				errors.As(err, &response)

				if err := ctx.Error(response.code, response.ErrorString); err != nil {
					api.log.SendErrorf(err.Error())
				}

				return
			}
		}

		if err := handler(ctx); err != nil {
			if err := ctx.InternalServerError(err.Error()); err != nil {
				api.log.SendErrorf(err.Error())
			}
		}
	}
}

// GET - implements GET api method call.
func (api *Service) GET(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodGet, path, handler, middlewares...)
	}
}

// PUT - implements PUT api method call.
func (api *Service) PUT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPut, path, handler, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func (api *Service) HEAD(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodHead, path, handler, middlewares...)
	}
}

// POST - implements POST api method call.
func (api *Service) POST(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPost, path, handler, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func (api *Service) PATCH(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPatch, path, handler, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func (api *Service) TRACE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodTrace, path, handler, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func (api *Service) DELETE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodDelete, path, handler, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func (api *Service) CONNECT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodConnect, path, handler, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func (api *Service) OPTIONS(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodOptions, path, handler, middlewares...)
	}
}
