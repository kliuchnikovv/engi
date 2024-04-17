package routes

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/api/parameter/placing"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/types"
)

type Route struct {
	handler Handler

	marshaler types.Marshaler
	responser types.Responser

	// cors   request.Middleware
	auth   func(r *http.Request, w http.ResponseWriter) error
	Body   Option
	Params map[placing.Placing][]Option //func(request *request.Request, response *response.Response) error // TODO: change type of params
	other  []Option                     //request.Middleware

	// CORS options
	allowedHeaders []string
	allowedMethods []string
	allowedOrigins []string
}

func NewRoute(
	handler Handler,
	marshaler types.Marshaler,
	responser types.Responser,
	options ...Option,
) (*Route, error) {
	var route = &Route{
		handler:   handler,
		marshaler: marshaler,
		responser: responser,
		auth: func(r *http.Request, w http.ResponseWriter) error {
			return nil
		},
		Params: make(map[placing.Placing][]Option),
	}

	for _, option := range options {
		if err := option.Bind(route); err != nil {
			return nil, err
		}
	}

	return route, nil
}

func (route *Route) Handle(
	r *http.Request,
	w http.ResponseWriter,
) error {
	// TODO: handle middlewares
	var (
		// TODO: handle in path params
		request  = request.New(r)
		response = response.New(
			w,
			route.marshaler,
			route.responser,
		)
	)

	if route.Body != nil {
		if err := route.Body.Handle(request, response); err != nil {
			return err
		}
	}

	for _, params := range route.Params {
		for _, param := range params {
			if err := param.Handle(request, response); err != nil {
				return err
			}
		}
	}

	for _, other := range route.other {
		if err := other.Handle(request, response); err != nil {
			return err
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

// func (route *Route) SetBody(
// 	body func(*request.Request, *response.Response) error,
// ) {
// 	route.body = body
// }

// func (route *Route) AddParam(
// 	param func(request *request.Request, response *response.Response) error,
// ) {
// 	route.params = append(route.params, param)
// }

// func (route *Route) AddParameter(
// 	key string,
// 	place placing.Placing,
// 	handler ParameterHandler,
// 	options ...request.Option,
// ) {
// 	if _, ok := route.params[place]; !ok {
// 		route.params[place] = make([]IParameter, 0, 1)
// 	}

// 	route.params[place] = append(route.params[place]) // IParameter{
// 	//		key:     key,
// 	//		place:   place,
// 	//		handler: handler,
// 	//		options: options,
// 	//	})
// }
