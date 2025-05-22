package engi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
	"github.com/kliuchnikovv/engi/internal/types"
)

var (
	ErrMethodNotAppliable = fmt.Errorf("method not appliable for path")
	ErrPathNotFound       = fmt.Errorf("path not found for method")
)

type (
	RouteByPath func(*Service, string, string) error
	Routes      map[RouteMethodPair]RouteByPath
	Middleware  routes.Middleware

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
		routes routes.Routes

		marshaler types.Marshaler // TODO: remove from here
		responser types.Responser

		logger *slog.Logger

		api  ServiceDefinition
		path string
	}

	RouteMethodPair struct {
		method string
		path   string
	}
)

func NewService(engine *Engine, api ServiceDefinition, path string) *Service {
	slog.New(engine.logger.Handler().WithAttrs([]slog.Attr{
		slog.String("service", api.Prefix()),
	}))

	return &Service{
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
) error {
	var middlewares []routes.Middleware

	for _, option := range srv.Middlewares() {
		middlewares = append(middlewares, option)
	}

	for _, option := range options {
		middlewares = append(middlewares, option)
	}

	return srv.routes.Add(
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

func (srv *Service) Serve(w http.ResponseWriter, r *http.Request) error {
	var uri, _ = strings.CutPrefix(r.URL.Path, srv.path)

	srv.logger.Debug("got request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	if err := srv.routes.Handle(context.Background(), r, w, r.Method, uri); err != nil {
		return err
	}

	return nil
}
