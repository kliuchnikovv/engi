package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo"
)

type (
	HandlerFunc func(*Context) error
	RouterFunc  func(group *echo.Group, path string) *echo.Route

	API interface {
		// Bind - производит привязку методов к путям апи.
		Bind(routers map[string]RouterFunc)

		// Routers - метод который должен быть реализован апи-обработчиком
		// и возвращает отображение обработчиков на имя функции.
		Routers() map[string]RouterFunc

		RegisterHandlers(g *echo.Group) error

		GET(HandlerFunc, ...HandlerFunc) RouterFunc
		PUT(HandlerFunc, ...HandlerFunc) RouterFunc
		HEAD(HandlerFunc, ...HandlerFunc) RouterFunc
		POST(HandlerFunc, ...HandlerFunc) RouterFunc
		PATCH(HandlerFunc, ...HandlerFunc) RouterFunc
		TRACE(HandlerFunc, ...HandlerFunc) RouterFunc
		DELETE(HandlerFunc, ...HandlerFunc) RouterFunc
		CONNECT(HandlerFunc, ...HandlerFunc) RouterFunc
		OPTIONS(HandlerFunc, ...HandlerFunc) RouterFunc

		WithBody(interface{}) HandlerFunc
		WithBool(key string) HandlerFunc
		WithInt(key string) HandlerFunc
		WithFloat(key string) HandlerFunc
		WithString(key string) HandlerFunc
		WithTime(key, layout string) HandlerFunc
	}

	BasicAPI struct {
		prefix  string
		routers map[string]RouterFunc
	}
)

func New(prefix string) API {
	return &BasicAPI{
		prefix: prefix,
	}
}

// Bind - производит привязку методов к путям апи.
func (api *BasicAPI) Bind(routers map[string]RouterFunc) {
	// TODO: remake - can be rewrited in service
	api.routers = routers
}

func (api *BasicAPI) Routers() map[string]RouterFunc {
	return nil
}

func (api *BasicAPI) RegisterHandlers(r *echo.Group) error {
	var group = r.Group(
		fmt.Sprintf(
			"/%s", strings.Trim(api.prefix, "/"),
		),
	)

	for path, register := range api.routers {
		register(group, fmt.Sprintf(
			"/%s", strings.Trim(path, "/"),
		))
	}

	return nil
}

func (api *BasicAPI) route(handler HandlerFunc, middlewares ...HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var context = NewContext(ctx)

		for _, middleware := range middlewares {
			if err := middleware(context); err != nil {
				var echoErr = new(echo.HTTPError)

				errors.As(err, &echoErr)

				return context.Response.Error(echoErr.Code, fmt.Errorf("%v", echoErr.Message))
			}
		}

		if err := handler(context); err != nil {
			return err
		}

		return nil
	}
}

func (api *BasicAPI) CONNECT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.CONNECT(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) DELETE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.DELETE(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) GET(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.GET(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) HEAD(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.HEAD(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) PATCH(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.PATCH(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) POST(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.POST(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) PUT(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.PUT(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) TRACE(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.TRACE(
			path,
			api.route(handler, middlewares...),
		)
	}
}

func (api *BasicAPI) OPTIONS(handler HandlerFunc, middlewares ...HandlerFunc) RouterFunc {
	return func(group *echo.Group, path string) *echo.Route {
		return group.OPTIONS(
			path,
			api.route(handler, middlewares...),
		)
	}
}

// func (api *BasicAPI) Use(middlewares ...HandlerFunc) {

// }

// TODO: debug middleware
// 	for _, route := range e.Echo.Routes() {
// 	fmt.Printf("Method: %s, path: %s\n", route.Method, route.Path)
// }

// for key, parameter := range ctx.Query.Parameters() {
// 	fmt.Printf("Key: %s, param: %s\n", key, parameter)
// }
