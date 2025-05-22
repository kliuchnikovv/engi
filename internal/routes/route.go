package routes

import (
	"context"
	"net/http"
	"sort"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/types"
)

type Route struct {
	Path    string
	handler Handler

	middlewares []Middleware

	Marshaler types.Marshaler
	Responser types.Responser

	// auth   func(r *http.Request, w http.ResponseWriter) error
	// Body   Option
	// Params map[placing.Placing]map[string]Option
	// other  []Option

	// CORS options
	// AllowedHeaders []string
	// AllowedMethods []string
	// AllowedOrigins []string
}

func NewRoute(
	path string,
	handler Handler,
	marshaler types.Marshaler,
	responser types.Responser,
	middlewares ...Middleware,
	// options ...Middleware,
) (*Route, error) {
	var route = Route{
		Path:        path,
		handler:     handler,
		Marshaler:   marshaler,
		Responser:   responser,
		middlewares: middlewares,
		// auth: func(r *http.Request, w http.ResponseWriter) error {
		// 	return nil
		// },
		// Params: make(map[placing.Placing]map[string]Middleware),
	}

	sort.Slice(route.middlewares, func(i, j int) bool {
		return route.middlewares[i].Priority() < route.middlewares[j].Priority()
	})

	// for _, option := range options {
	// 	if err := option.Bind(&route); err != nil {
	// 		return nil, err
	// 	}
	// }

	return &route, nil
}

func (route *Route) Handle(
	ctx context.Context,
	request *request.Request,
	writer http.ResponseWriter,
) error {
	var response = response.New(writer,
		route.Marshaler,
		route.Responser,
	)

	for _, middleware := range route.middlewares {
		if err := middleware.Handle(ctx, request, response); err != nil {
			return response.BadRequest(err.Error())
		}
	}

	return route.handler(
		ctx,
		request,
		response,
	)
}

// func (route *Route) newRequest(
// 	r *http.Request,
// 	path string,
// ) *request.Request {
// 	var (
// 		request    = request.New(r)
// 		pathPieces = strings.Split(path, "/")
// 	)
// 	if !parameterRegexp.MatchString(route.Path) {
// 		return request
// 	}

// 	for i, paramName := range strings.Split(route.Path, "/") {
// 		if !parameterRegexp.MatchString(paramName) {
// 			continue
// 		}

// 		paramName = strings.Trim(paramName, "{}")

// 		_, ok := route.Params[placing.InPath][strings.Trim(paramName, "{}")]
// 		if !ok {
// 			panic(
// 				fmt.Sprintf("in-path parameter not found: %s", paramName),
// 			)
// 		}

// 		request.AddInPathParameter(paramName, pathPieces[i])
// 	}

// 	return request
// }
