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
func GET(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodGet, path, route, middlewares...)
	}
}

// PUT - implements PUT api method call.
func PUT(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPut, path, route, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func HEAD(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodHead, path, route, middlewares...)
	}
}

// POST - implements POST api method call.
func POST(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPost, path, route, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func PATCH(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodPatch, path, route, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func TRACE(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodTrace, path, route, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func DELETE(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodDelete, path, route, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func CONNECT(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodConnect, path, route, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func OPTIONS(route Route, middlewares ...Register) RouteByPath {
	return func(srv *Service, path string) error {
		return srv.add(http.MethodOptions, path, route, middlewares...)
	}
}
