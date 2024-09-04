package engi

import (
	"context"
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
)

type (
	Request  request.Requester
	Response response.Responser
	Route    func(ctx context.Context, request Request, response Response) error
)

// GET - implements GET api method call.
func GET(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodGet,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// PUT - implements PUT api method call.
func PUT(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodPut,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// HEAD - implements HEAD api method call.
func HEAD(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodHead,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// POST - implements POST api method call.
func POST(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodPost,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// PATCH - implements PATCH api method call.
func PATCH(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodPatch,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// TRACE - implements TRACE api method call.
func TRACE(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodTrace,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// DELETE - implements DELETE api method call.
func DELETE(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodDelete,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// CONNECT - implements CONNECT api method call.
func CONNECT(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodConnect,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

// OPTIONS - implements OPTIONS api method call.
func OPTIONS(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, path string) {
		srv.addRoute(
			http.MethodOptions,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}
