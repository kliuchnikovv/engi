package api

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
)

type (
	// HandlerFunc - describes any handler.
	HandlerFunc func(*Context) error

	// RouterFunc - describes handler prepared for registration.
	RouterFunc func(group *echo.Group, path string) *echo.Route

	// ServiceAPI - describes service methods for handlers creation and registering.
	ServiceAPI interface {
		// PathPrefix - prefix of all paths for this service.
		PathPrefix() string
		// Routers returns the handlers and their relative paths (relative to the service) for registration.
		Routers() map[string]RouterFunc

		// GET - implements GET api method call.
		GET(HandlerFunc, ...HandlerFunc) RouterFunc
		// PUT - implements PUT api method call.
		PUT(HandlerFunc, ...HandlerFunc) RouterFunc
		// HEAD - implements HEAD api method call.
		HEAD(HandlerFunc, ...HandlerFunc) RouterFunc
		// POST - implements POST api method call.
		POST(HandlerFunc, ...HandlerFunc) RouterFunc
		// PATCH - implements PATCH api method call.
		PATCH(HandlerFunc, ...HandlerFunc) RouterFunc
		// TRACE - implements TRACE api method call.
		TRACE(HandlerFunc, ...HandlerFunc) RouterFunc
		// DELETE - implements DELETE api method call.
		DELETE(HandlerFunc, ...HandlerFunc) RouterFunc
		// CONNECT - implements CONNECT api method call.
		CONNECT(HandlerFunc, ...HandlerFunc) RouterFunc
		// OPTIONS - implements OPTIONS api method call.
		OPTIONS(HandlerFunc, ...HandlerFunc) RouterFunc

		// WithBody - takes pointer to structure and saves casted request body into context.
		// Result can be retrieved from context using 'context.QueryParams.Body()'.
		WithBody(interface{}) HandlerFunc
		// WithBool - queries mandatory boolean parameter from request by 'key'.
		// Result can be retrieved from context using 'context.QueryParams.Bool(key)'.
		WithBool(key string) HandlerFunc
		// WithInt - queries mandatory integer parameter from request by 'key'.
		// Result can be retrieved from context using 'context.QueryParams.Integer(key)'.
		WithInt(key string) HandlerFunc
		// WithFloat - queries mandatory floating point number parameter from request by 'key'.
		// Result can be retrieved from context using 'context.QueryParams.Float(key)'.
		WithFloat(key string) HandlerFunc
		// WithString - queries mandatory string parameter from request by 'key'.
		// Result can be retrieved from context using 'context.QueryParams.String(key)'.
		WithString(key string) HandlerFunc
		// WithTime - queries mandatory time parameter from request by 'key' using 'layout'.
		// Result can be retrieved from context using 'context.QueryParams.Time(key, layout)'.
		WithTime(key, layout string) HandlerFunc
	}

	// ServiceBase - provides service methods for handlers creation and registering.
	ServiceBase struct {
		prefix string
	}
)

func New(prefix string) ServiceAPI {
	return &ServiceBase{
		prefix: prefix,
	}
}

// PathPrefix - prefix of all paths for this service.
func (api *ServiceBase) PathPrefix() string {
	return api.prefix
}

// Routers returns the handlers and their relative paths (relative to the service) for registration.
func (api *ServiceBase) Routers() map[string]RouterFunc {
	return nil
}

// route - creates 'webapi.Context' from 'echo.Context' and wraps handler function calls.
func (api *ServiceBase) route(handler HandlerFunc, middlewares ...HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var context = NewContext(ctx)

		for _, middleware := range middlewares {
			if err := middleware(context); err != nil {
				var echoErr = new(echo.HTTPError)

				errors.As(err, &echoErr)

				return context.Response.Error(echoErr.Code, fmt.Errorf("%v", echoErr.Message))
			}
		}

		return handler(context)
	}
}

// GET - implements GET api method call.
func (api *ServiceBase) GET(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.GET(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// PUT - implements PUT api method call.
func (api *ServiceBase) PUT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.PUT(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// HEAD - implements HEAD api method call.
func (api *ServiceBase) HEAD(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.HEAD(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// POST - implements POST api method call.
func (api *ServiceBase) POST(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.POST(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// PATCH - implements PATCH api method call.
func (api *ServiceBase) PATCH(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.PATCH(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// TRACE - implements TRACE api method call.
func (api *ServiceBase) TRACE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.TRACE(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// DELETE - implements DELETE api method call.
func (api *ServiceBase) DELETE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.DELETE(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// CONNECT - implements CONNECT api method call.
func (api *ServiceBase) CONNECT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.CONNECT(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// OPTIONS - implements OPTIONS api method call.
func (api *ServiceBase) OPTIONS(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.OPTIONS(
			path,
			api.route(handler, middlewares...),
		)
	}
}
