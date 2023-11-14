package engi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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
		Middlewares() []Middleware
	}

	// Service - provides basic service methods.
	Service struct {
		middlewares []request.Middleware

		handlers map[string]*pathfinder.PathFinder

		marshaler types.Marshaler
		responser types.Responser

		logger *slog.Logger

		api  ServiceAPI
		path string
	}

	RouteByPath func(*Service, string) error
	Middleware  func(*Service)
	Routes      map[string]RouteByPath
)

func NewService(engine *Engine, api ServiceAPI, path string) *Service {
	return &Service{
		handlers: make(map[string]*pathfinder.PathFinder),

		logger:    engine.logger,
		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,

		api:  api,
		path: path,
	}
}

func (srv *Service) Middlewares() []Middleware {
	if middlewaresAPI, ok := srv.api.(MiddlewaresAPI); ok {
		return middlewaresAPI.Middlewares()
	}

	return nil
}

// add - creates route with custom method and path.
func (srv *Service) add(
	method, path string,
	route Route,
	middlewares ...request.Middleware,
) error {
	if _, ok := srv.handlers[method]; !ok {
		srv.handlers[method] = pathfinder.NewPathFinder()
	}

	for _, middleware := range middlewares {
		validator, ok := middleware.(request.ParamsValidator)
		if !ok {
			continue
		}

		if err := validator.Validate(path); err != nil {
			return fmt.Errorf("%w, service: %s", err, srv.api.Prefix())
		}
	}

	srv.handlers[method].Add(path, srv.handle(route, middlewares...))

	return nil
}

// Serve should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// Serve call.
func (srv *Service) Serve(uri string, r *http.Request, w http.ResponseWriter) {
	uri, _ = strings.CutPrefix(strings.Trim(uri, "/"), srv.api.Prefix())

	if err := srv.serve(r.Context(),
		request.New(r),
		response.New(w, srv.marshaler, srv.responser),
		r.Method, uri,
	); err != nil {
		srv.logger.Error(err.Error())
	}
}

func (srv *Service) serve(
	ctx context.Context,
	request *request.Request,
	response *response.Response,
	method,
	uri string,
) error {
	for _, middleware := range srv.middlewares {
		var obj = middleware.Handle(request, response.ResponseWriter())
		if obj == nil {
			continue
		}

		switch obj.Code {
		case http.StatusOK:
			return response.OK(obj.Code)
		default:
			return response.Error(obj.Code, obj.ErrorString)
		}
	}

	if _, ok := srv.handlers[method]; !ok {
		return response.NotFound(ErrMethodNotAppliable.Error())
	}

	var err = srv.handlers[method].Handle(ctx, request, response, strings.Trim(uri, "/"))
	if err == nil {
		return nil
	}

	if errors.Is(err, pathfinder.ErrNotHandled) {
		return response.NotFound(ErrPathNotFound.Error())
	}

	return response.InternalServerError(err.Error())
}

func (srv *Service) handle(
	route Route,
	middlewares ...request.Middleware,
) pathfinder.Handler {
	return func(ctx context.Context, request *request.Request, response *response.Response) error {
		for _, middleware := range middlewares {
			var obj = middleware.Handle(request, response.ResponseWriter())
			if obj != nil {
				return response.Error(obj.Code, obj.ErrorString)
			}
		}

		if err := route(ctx, request, response); err != nil {
			if err := response.InternalServerError(err.Error()); err != nil {
				return err
			}
		}

		return nil
	}
}
