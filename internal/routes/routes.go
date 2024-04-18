package routes

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/KlyuchnikovV/engi/definition/parameter/placing"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

var (
	parameterRegexp = regexp.MustCompile("{[a-zA-Z]*}")

	ErrPathNotFound = errors.New("path not found")
)

type (
	Handler func(ctx context.Context, request *request.Request, response *response.Response) error

	regexpRoute struct {
		regexp regexp.Regexp

		route *Route
	}

	Routes struct {
		exactHandlers map[string]map[string]*Route

		regexpHandlers map[string][]regexpRoute
	}
)

func New() Routes {
	return Routes{
		exactHandlers:  make(map[string]map[string]*Route),
		regexpHandlers: make(map[string][]regexpRoute),
	}
}

func (routes Routes) Add(
	method,
	path string,
	handler Handler,
	marshaler types.Marshaler,
	responser types.Responser,
	options ...Option,
) error {
	route, err := NewRoute(path, handler, marshaler, responser, options...)
	if err != nil {
		return err
	}

	if !parameterRegexp.MatchString(path) {
		if _, ok := routes.exactHandlers[method]; !ok {
			routes.exactHandlers[method] = make(map[string]*Route)
		}

		routes.exactHandlers[method][path] = route

		return nil
	}

	if _, ok := routes.regexpHandlers[method]; !ok {
		routes.regexpHandlers[method] = make([]regexpRoute, 0, 1)
	}

	for _, param := range route.Params[placing.InPath] {
		named, ok := param.(NamedParameter)
		if !ok {
			panic("parameter can't be unnamed")
		}

		var (
			name  = named.Name()
			reg   = named.Regexp()
			param = fmt.Sprintf("{%s}", name)
		)
		if !strings.Contains(path, param) {
			return fmt.Errorf("in-path parameter not found in path: %s", name)
		}

		path = strings.Replace(path, param, reg, 1)
	}

	path = strings.Replace(path, "/", `\/`, 1)

	routes.regexpHandlers[method] = append(routes.regexpHandlers[method], regexpRoute{
		regexp: *regexp.MustCompile(path),
		route:  route,
	})

	return nil
}

func (routes Routes) Handle(
	r *http.Request,
	w http.ResponseWriter,
	path string,
) error {
	route, err := routes.matchEndpoint(r.Method, path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))

		return nil
	}

	if code, err := route.cors(r, w); err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))

		return nil
	}

	if err := route.auth(r, w); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return nil
	}

	if err := route.Handle(r, w, path); err != nil {
		// w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte(err.Error()))
		fmt.Print(err) // TODO: logger

		return nil
	}

	return nil
}

func (routes Routes) matchEndpoint(method, path string) (*Route, error) {
	var (
		exactRoutes, exactHandlerFound   = routes.exactHandlers[method]
		regexpRoutes, regexpHandlerFound = routes.regexpHandlers[method]
	)

	if !exactHandlerFound && !regexpHandlerFound {
		return nil, ErrPathNotFound
	}

	if exactHandlerFound {
		route, ok := exactRoutes[path]
		if ok {
			return route, nil
		}
	}

	var (
		index   = 0
		minArgs = math.MaxInt
		result  = make([]regexpRoute, 0)
	)
	for _, route := range regexpRoutes {
		if !route.regexp.MatchString(path) {
			continue
		}

		result = append(result, route)

		// TODO: fix description
		// Trying to find route with less number of in path parameters
		// to ensure that maximum of non-parametrised path pieces is used
		//
		//   get/{id} ==> with 'get/2' will go there
		// {obj}/{id} ==> with 'abc/2' will go there
		if minArgs > len(route.route.Params[placing.InPath]) {
			minArgs = len(route.route.Params[placing.InPath])
			index = len(result) - 1
		}
	}

	switch len(result) {
	case 0:
		return nil, ErrPathNotFound
	case 1:
		return result[0].route, nil
	default:
		return result[index].route, nil
	}
}
