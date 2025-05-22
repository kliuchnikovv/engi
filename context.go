package engi

import (
	"context"
	"net/http"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
)

type (
	Request  request.Requester
	Response response.Responser
	Route    func(ctx context.Context, request Request, response Response) error
)

func Handle(route Route, middlewares ...Middleware) RouteByPath {
	return func(srv *Service, method, path string) error {
		return srv.addRoute(
			method,
			path,
			func(ctx context.Context, request *request.Request, response *response.Response) error {
				return route(ctx, request, response)
			},
			middlewares...,
		)
	}
}

func NewMethod(method, path string) RouteMethodPair {
	return RouteMethodPair{
		method: method,
		path:   path,
	}
}

// GET - implements GET api method call.
func GET(path string) RouteMethodPair {
	return NewMethod(http.MethodGet, path)
}

// PUT - implements PUT api method call.
func PUT(path string) RouteMethodPair {
	return NewMethod(http.MethodPut, path)
}

// HED - implements HEAD api method call.
func HED(path string) RouteMethodPair {
	return NewMethod(http.MethodHead, path)
}

// PST - implements POST api method call.
func PST(path string) RouteMethodPair {
	return NewMethod(http.MethodPost, path)
}

// PTC - implements PATCH api method call.
func PTC(path string) RouteMethodPair {
	return NewMethod(http.MethodPatch, path)
}

// TRC - implements TRACE api method call.
func TRC(path string) RouteMethodPair {
	return NewMethod(http.MethodTrace, path)
}

// DEL - implements DELETE api method call.
func DEL(path string) RouteMethodPair {
	return NewMethod(http.MethodDelete, path)
}

// CNT - implements CONNECT api method call.
func CNT(path string) RouteMethodPair {
	return NewMethod(http.MethodConnect, path)
}

// OPT - implements OPTIONS api method call.
func OPT(path string) RouteMethodPair {
	return NewMethod(http.MethodOptions, path)
}
