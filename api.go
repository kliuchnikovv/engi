package webapi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type (
	ResultFunc  func(http.ResponseWriter, *http.Request, ResponseMarshaler, Responser)
	RouterFunc  func(path string)
	HandlerFunc func(*Context) error

	ServiceBase struct {
		logger log.Logger

		prefix string

		middlewares []HandlerFunc

		responseObject Responser
		marshaler      ResponseMarshaler

		routes map[string]map[string]ResultFunc
	}
)

func NewService(prefix string) ServiceAPI {
	return &ServiceBase{
		prefix: prefix,
		routes: make(map[string]map[string]ResultFunc),
	}
}

// TODO:
func (api *ServiceBase) bind(
	responseBinder ResponseMarshaler,
	responseObject Responser,
) {
	api.marshaler = responseBinder
	api.responseObject = responseObject
}

func (api *ServiceBase) PathPrefix() string {
	return api.prefix
}

func (api *ServiceBase) Routers() map[string]RouterFunc {
	return nil
}

func (api *ServiceBase) Use(middlewares ...HandlerFunc) {
	api.middlewares = append(api.middlewares, middlewares...)
}

func (api *ServiceBase) Add(method, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	if _, ok := api.routes[method]; !ok {
		api.routes[method] = make(map[string]ResultFunc)
	}

	if _, ok := api.routes[method][path]; ok {
		// TODO: errors channel
		api.logger.Fatalf("method on path already defined")
		return
	}

	api.routes[method][path] = func(
		response http.ResponseWriter,
		request *http.Request,
		responseMarshaler ResponseMarshaler,
		responseObject Responser,
	) {
		var (
			ctx  = NewContext(response, request, responseMarshaler, responseObject)
			resp ResponseObject
		)

		for _, middleware := range api.middlewares {
			if err := middleware(ctx); err != nil {
				errors.As(err, &resp)
				ctx.Error(resp.code, resp.ErrorString)

				return
			}
		}

		for _, middleware := range middlewares {
			if err := middleware(ctx); err != nil {
				errors.As(err, &resp)
				ctx.Error(
					resp.code, resp.ErrorString,
				)

				return
			}
		}

		if err := handler(ctx); err != nil {
			ctx.InternalServerError(err.Error())

			return
		}
	}
}

// GET - implements GET api method call.
func (api *ServiceBase) GET(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodGet, path, handler, middlewares...)
	}
}

// PUT - implements PUT api method call.
func (api *ServiceBase) PUT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPut, path, handler, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func (api *ServiceBase) HEAD(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodHead, path, handler, middlewares...)
	}
}

// POST - implements POST api method call.
func (api *ServiceBase) POST(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPost, path, handler, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func (api *ServiceBase) PATCH(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodPatch, path, handler, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func (api *ServiceBase) TRACE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodTrace, path, handler, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func (api *ServiceBase) DELETE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodDelete, path, handler, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func (api *ServiceBase) CONNECT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodConnect, path, handler, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func (api *ServiceBase) OPTIONS(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(path string) {
		api.Add(http.MethodOptions, path, handler, middlewares...)
	}
}

func (api *ServiceBase) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, path: %s\n", r.Method, r.URL.Path)

	if _, ok := api.routes[r.Method]; !ok {
		http.Error(w, fmt.Sprintf("method '%s' not appliable for '%s'", r.Method, r.URL.Path), http.StatusNotFound)
		return
	}

	route, ok := api.routes[r.Method][r.URL.Path]
	if !ok {
		for method, routes := range api.routes {
			for route := range routes {
				log.Printf("Method: %s, route: %s\n", method, route)
			}
		}

		http.Error(w, "not found", http.StatusNotFound)

		return
	}

	route(w, r, api.marshaler, api.responseObject)
}
