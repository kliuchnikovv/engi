package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/types"
)

// TODO: trie

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
		root *Trie[*Route]
	}
)

func New() Routes {
	return Routes{
		root: NewTrie[*Route](),
	}
}

func (routes Routes) Add(
	method,
	path string,
	handler Handler,
	marshaler types.Marshaler,
	responser types.Responser,
	options ...Middleware,
) error {
	route, err := NewRoute(path, handler, marshaler, responser, options...)
	if err != nil {
		return err
	}

	routes.root.Add(method, path, route)

	return nil
}

func (routes Routes) Handle(
	ctx context.Context,
	req *http.Request,
	writer http.ResponseWriter,
	method string,
	path string,
) error {
	var request = request.New(req)

	route, err := routes.root.Get(request, method, path)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(err.Error()))

		return nil
	}

	if err := (*route).Handle(ctx, request, writer); err != nil {
		fmt.Print(err) // TODO: logger

		return nil
	}

	return nil
}

// func (routes Routes) matchEndpoint(method, path string) (*Route, error) {
// 	var (
// 		exactRoutes, exactHandlerFound   = routes.exactHandlers[method]
// 		regexpRoutes, regexpHandlerFound = routes.regexpHandlers[method]
// 	)

// 	if !exactHandlerFound && !regexpHandlerFound {
// 		return nil, ErrPathNotFound
// 	}

// 	if exactHandlerFound {
// 		route, ok := exactRoutes[path]
// 		if ok {
// 			return route, nil
// 		}
// 	}

// 	var (
// 		index   = 0
// 		minArgs = math.MaxInt
// 		result  = make([]regexpRoute, 0)
// 	)
// 	for _, route := range regexpRoutes {
// 		if !route.regexp.MatchString(path) {
// 			continue
// 		}

// 		result = append(result, route)

// 		// TODO: fix description
// 		// Trying to find route with less number of in path parameters
// 		// to ensure that maximum of non-parametrised path pieces is used
// 		//
// 		//   get/{id} ==> with 'get/2' will go there
// 		// {obj}/{id} ==> with 'abc/2' will go there
// 		if minArgs > len(route.route.Params[placing.InPath]) {
// 			minArgs = len(route.route.Params[placing.InPath])
// 			index = len(result) - 1
// 		}
// 	}

// 	switch len(result) {
// 	case 0:
// 		return nil, ErrPathNotFound
// 	case 1:
// 		return result[0].route, nil
// 	default:
// 		return result[index].route, nil
// 	}
// }
