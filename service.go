package engi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/middlewares"
	"github.com/KlyuchnikovV/engi/internal/pathfinder"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

var (
	ErrMethodNotAppliable = fmt.Errorf("method not appliable for path")
	ErrPathNotFound       = fmt.Errorf("path not found for method")
)

type (
	ServiceAPI interface {
		// Prefix - prefix of all paths for this service.
		Prefix() string

		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		Routers() Routes
	}

	MiddlewaresAPI interface {
		Middlewares() []Register
	}

	// Service - provides basic service methods.
	Service struct {
		handlers map[string]*pathfinder.PathFinder

		marshaler types.Marshaler
		responser types.Responser

		logger *slog.Logger

		api  ServiceAPI
		path string
	}

	RouteByPath func(*Service, string) error
	Routes      map[string]RouteByPath
)

func NewService(engine *Engine, api ServiceAPI, path string) *Service {
	slog.New(engine.logger.Handler().WithAttrs([]slog.Attr{
		slog.String("service", api.Prefix()),
	}))

	return &Service{
		handlers: make(map[string]*pathfinder.PathFinder),

		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,

		api:  api,
		path: path,

		logger: slog.New(engine.logger.Handler().WithAttrs([]slog.Attr{
			slog.String("service", api.Prefix()),
		})),
	}
}

func (srv *Service) Middlewares() []Register {
	if middlewaresAPI, ok := srv.api.(MiddlewaresAPI); ok {
		return middlewaresAPI.Middlewares()
	}

	return nil
}

// add - creates route with custom method and path.
func (srv *Service) add(
	method, path string,
	route Route,
	options ...Register,
) error {
	if _, ok := srv.handlers[method]; !ok {
		srv.handlers[method] = pathfinder.NewPathFinder()
	}

	var middlewares = middlewares.New()
	for _, middleware := range srv.Middlewares() {
		middleware(middlewares)
	}

	for _, middleware := range options {
		middleware(middlewares)
	}

	srv.handlers[method].Add(path, srv.handleEndpoint(
		route,
		middlewares,
	))

	return nil
}

func (srv *Service) Serve(
	w http.ResponseWriter, r *http.Request, uri string,
) error {
	srv.logger.Debug("got request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	var (
		request  = request.New(r)
		response = response.New(w, srv.marshaler, srv.responser)
	)

	if _, ok := srv.handlers[r.Method]; !ok {
		return response.NotFound(ErrMethodNotAppliable.Error())
	}

	var handler = srv.handlers[r.Method].Handle(request, strings.Trim(uri, "/"))
	if handler == nil {
		return response.NotFound(ErrPathNotFound.Error())
	}

	return handler(r.Context(), request, response)
}

func (srv *Service) handleEndpoint(
	route Route,
	middlewares *middlewares.Middlewares,
) pathfinder.Handler {
	return func(ctx context.Context, request *request.Request, response *response.Response) error {
		if err := middlewares.Handle(request, response.ResponseWriter()); err != nil {
			return response.Error(err.Code, err.ErrorString)
		}

		if err := route(ctx, request, response); err != nil {
			if err := response.InternalServerError(err.Error()); err != nil {
				return err
			}
		}

		return nil
	}
}
