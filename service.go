package engi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
	"github.com/KlyuchnikovV/engi/internal/types"
)

var (
	ErrMethodNotAppliable = fmt.Errorf("method not appliable for path")
	ErrPathNotFound       = fmt.Errorf("path not found for method")
)

type (
	RouteByPath func(*Service, string)
	Routes      map[string]RouteByPath
	Middleware  routes.Option

	ServiceDefinition interface {
		// Prefix - prefix of all paths for this service.
		Prefix() string

		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		Routers() Routes
	}

	MiddlewaresAPI interface {
		Middlewares() []Middleware
	}

	// Service - provides basic service methods.
	Service struct {
		// handlers map[string]*pathfinder.PathFinder
		routes routes.Routes

		marshaler types.Marshaler // TODO: remove from here
		responser types.Responser

		logger *slog.Logger

		api  ServiceDefinition
		path string
	}
)

func NewService(engine *Engine, api ServiceDefinition, path string) *Service {
	slog.New(engine.logger.Handler().WithAttrs([]slog.Attr{
		slog.String("service", api.Prefix()),
	}))

	return &Service{
		// handlers: make(map[string]*pathfinder.PathFinder),
		routes: routes.New(),

		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,

		api:  api,
		path: path,

		logger: slog.New(engine.logger.Handler().WithAttrs([]slog.Attr{
			slog.String("service", api.Prefix()),
		})),
	}
}

func (srv *Service) Middlewares() []Middleware {
	if middlewaresAPI, ok := srv.api.(MiddlewaresAPI); ok {
		return middlewaresAPI.Middlewares()
	}

	return nil
}

func (srv *Service) addRoute(
	method,
	path string,
	route routes.Handler,
	options ...Middleware,
) {
	var middlewares []routes.Option

	for _, option := range srv.Middlewares() {
		middlewares = append(middlewares, option)
	}

	for _, option := range options {
		middlewares = append(middlewares, option)
	}

	srv.routes.Add(
		method,
		path,
		func(ctx context.Context, request *request.Request, response *response.Response) error {
			return route(ctx, request, response)
		},
		srv.marshaler,
		srv.responser,
		middlewares...,
	)
}

// add - creates route with custom method and path.
// func (srv *Service) add(
// 	method, path string,
// 	route Route,
// 	options ...Register,
// ) error {
// 	if _, ok := srv.handlers[method]; !ok {
// 		srv.handlers[method] = pathfinder.NewPathFinder()
// 	}

// 	var middlewares = middlewares.New()
// 	for _, middleware := range srv.Middlewares() {
// 		middleware(middlewares)
// 	}

// 	for _, middleware := range options {
// 		middleware(middlewares)
// 	}

// 	srv.handlers[method].Add(path, srv.handleEndpoint(
// 		route,
// 		middlewares,
// 	))

// 	return nil
// }

// func (srv *Service) handleEndpoint(
// 	route Route,
// 	middlewares *middlewares.Middlewares,
// ) pathfinder.Handler {
// 	return func(ctx context.Context, request *request.Request, response *response.Response) error {
// 		if err := middlewares.Handle(request, response.ResponseWriter()); err != nil {
// 			return response.Error(err.Code, err.ErrorString)
// 		}

// 		if err := route(ctx, request, response); err != nil {
// 			if err := response.InternalServerError(err.Error()); err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	}
// }

func (srv *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var uri, _ = strings.CutPrefix(r.URL.Path, srv.path)

	srv.logger.Debug("got request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	if err := srv.routes.Handle(r, w, uri); err != nil {
		panic(err)
	}

	// if _, ok := srv.handlers[r.Method]; !ok {
	// 	srv.logger.Error(
	// 		response.NotFound(
	// 			ErrMethodNotAppliable.Error(),
	// 		).Error(),
	// 	)
	// }

	// var handler = srv.handlers[r.Method].Handle(request, strings.Trim(uri, "/"))
	// if handler == nil {
	// 	srv.logger.Error(
	// 		response.NotFound(
	// 			ErrPathNotFound.Error(),
	// 		).Error(),
	// 	)
	// }

	// if err := handler(r.Context(), request, response); err != nil {
	// 	srv.logger.Error(err.Error())
	// }

	srv.logger.Debug("request handled")
}
