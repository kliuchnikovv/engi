package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/definition/parameter/placing"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type Route struct {
	Path    string
	handler Handler

	marshaler types.Marshaler
	responser types.Responser

	// cors   request.Middleware
	auth   func(r *http.Request, w http.ResponseWriter) error
	Body   Option
	Params map[placing.Placing]map[string]Option
	other  []Option

	// CORS options
	AllowedHeaders []string
	AllowedMethods []string
	AllowedOrigins []string
}

func NewRoute(
	path string,
	handler Handler,
	marshaler types.Marshaler,
	responser types.Responser,
	options ...Option,
) (*Route, error) {
	var route = Route{
		Path:      path,
		handler:   handler,
		marshaler: marshaler,
		responser: responser,
		auth: func(r *http.Request, w http.ResponseWriter) error {
			return nil
		},
		Params: make(map[placing.Placing]map[string]Option),
	}

	for _, option := range options {
		if err := option.Bind(&route); err != nil {
			return nil, err
		}
	}

	return &route, nil
}

func (route *Route) Handle(
	r *http.Request,
	w http.ResponseWriter,
	path string,
) error {
	var (
		request  = route.newRequest(r, path)
		response = response.New(w,
			route.marshaler,
			route.responser,
		)
	)

	if route.Body != nil {
		if err := route.Body.Handle(request, response); err != nil {
			// return err
			return response.BadRequest(err.Error())
		}
	}

	for _, params := range route.Params {
		for _, param := range params {
			if err := param.Handle(request, response); err != nil {
				return response.BadRequest(err.Error())
			}
		}
	}

	for _, other := range route.other {
		if err := other.Handle(request, response); err != nil {
			// return err
			return response.BadRequest(err.Error())
		}
	}

	return route.handler(
		r.Context(),
		request,
		response,
	)
}

func (route *Route) SetAuth(
	auth func(*http.Request, http.ResponseWriter) error,
) {
	route.auth = auth
}

func (route *Route) newRequest(
	r *http.Request,
	path string,
) *request.Request {
	var (
		request    = request.New(r)
		pathPieces = strings.Split(path, "/")
	)
	if !parameterRegexp.MatchString(route.Path) {
		return request
	}

	for i, paramName := range strings.Split(route.Path, "/") {
		if !parameterRegexp.MatchString(paramName) {
			continue
		}

		paramName = strings.Trim(paramName, "{}")

		_, ok := route.Params[placing.InPath][strings.Trim(paramName, "{}")]
		if !ok {
			panic(
				fmt.Sprintf("in-path parameter not found: %s", paramName),
			)
		}

		request.AddInPathParameter(paramName, pathPieces[i])
	}

	return request
}
