package webapi

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/KlyuchnikovV/webapi/options"
	"github.com/KlyuchnikovV/webapi/types"
)

var parameterRegexp = regexp.MustCompile("{[a-zA-Z]*}")

type (
	ServiceAPI interface {
		// Prefix - prefix of all paths for this service.
		Prefix() string

		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		Routers() map[string]RouterByPath
	}

	MiddlewaresAPI interface {
		Middlewares() []Middleware
	}

	// Service - provides basic service methods.
	Service struct {
		middlewares []options.HandlerParams

		handlers map[string]map[string]Handler

		marshaler types.Marshaler
		responser types.Responser

		log *types.Log

		api ServiceAPI
	}
)

func NewService(engine *Engine, api ServiceAPI) *Service {
	return &Service{
		handlers: make(map[string]map[string]Handler),

		log:       engine.log,
		marshaler: engine.responseMarshaler,
		responser: engine.responseObject,

		api: api,
	}
}

// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
func (api *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = NewContext(w, r, api.marshaler, api.responser)

	for _, middleware := range api.middlewares {
		err := middleware(ctx.Request, ctx.Response.Response())
		if err == nil {
			continue
		}

		var response types.ResponseObject

		errors.As(err, &response)

		switch response.Code {
		case http.StatusOK:
			err = ctx.Response.OK(response.Code)
		default:
			err = ctx.Response.Error(response.Code, response.ErrorString)
		}

		if err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	if _, ok := api.handlers[r.Method]; !ok {
		if err := ctx.Response.NotFound("method '%s' not appliable for '%s'", r.Method, r.URL.Path); err != nil {
			api.log.SendErrorf(err.Error())
		}

		return
	}

	if handler, ok := api.handlers[r.Method][r.URL.Path]; ok {
		handler(ctx)
		return
	}

	for path, handler := range api.handlers[r.Method] {
		if regexp.MustCompile(path).MatchString(r.URL.Path) {
			handler(ctx)
			return
		}
	}

	if err := ctx.Response.NotFound("path '%s' not found for method '%s'", r.URL.Path, r.Method); err != nil {
		api.log.SendErrorf(err.Error())
	}
}

// add - creates route with custom method and path.
func (api *Service) add(
	method, path string,
	route Route,
	middlewares ...options.HandlerParams,
) {
	if _, ok := api.handlers[method]; !ok {
		api.handlers[method] = make(map[string]Handler)
	}

	if _, ok := api.handlers[method][path]; ok {
		api.log.SendErrorf("method '%s' with path '%s' already defined", method, path)
		return
	}

	if strings.ContainsAny(path, "{}") {
		middlewares = append([]options.HandlerParams{parseInPathParameters(path)}, middlewares...)
		path = parameterRegexp.ReplaceAllString(path, "[a-zA-Z0-9]+")
	}

	api.handlers[method][path] = api.handle(route, middlewares...)
}

func (api *Service) handle(route Route, middlewares ...options.HandlerParams) Handler {
	return func(ctx *Context) {
		var response types.ResponseObject

		for _, middleware := range middlewares {
			if err := middleware(ctx.Request, ctx.Response.Response()); err != nil {
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
func GET(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodGet, path, route, middlewares...)
	}
}

// PUT - implements PUT api method call.
func PUT(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPut, path, route, middlewares...)
	}
}

// HEAD - implements HEAD api method call.
func HEAD(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodHead, path, route, middlewares...)
	}
}

// POST - implements POST api method call.
func POST(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPost, path, route, middlewares...)
	}
}

// PATCH - implements PATCH api method call.
func PATCH(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodPatch, path, route, middlewares...)
	}
}

// TRACE - implements TRACE api method call.
func TRACE(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodTrace, path, route, middlewares...)
	}
}

// DELETE - implements DELETE api method call.
func DELETE(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodDelete, path, route, middlewares...)
	}
}

// CONNECT - implements CONNECT api method call.
func CONNECT(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodConnect, path, route, middlewares...)
	}
}

// OPTIONS - implements OPTIONS api method call.
func OPTIONS(route Route, middlewares ...options.HandlerParams) RouterByPath {
	return func(api *Service, path string) {
		api.add(http.MethodOptions, path, route, middlewares...)
	}
}

func parseInPathParameters(pathTemplate string) options.HandlerParams {
	var templateParams = strings.Split(pathTemplate, "/")

	return func(request *options.Request, _ http.ResponseWriter) error {
		var (
			path       = request.Request().URL.Path
			pathParams = strings.Split(path, "/")
		)

		if len(pathParams) < len(templateParams) {
			return fmt.Errorf(
				"number of path params is less than in template (got: %d, want: %d)",
				len(pathParams), len(templateParams),
			)
		}

		for i, template := range templateParams {
			if pathParams[i] == template {
				continue
			}

			request.AddInPathParameter(template[1:len(template)-1], pathParams[i])
		}

		return nil
	}
}
