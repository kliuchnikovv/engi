package cors

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
)

type CORS struct {
	name string

	handle func(r *http.Request, w http.ResponseWriter) error
}

func (cors *CORS) Bind(route *routes.Route) error {
	// route.SetAuth(auth.handle)

	return nil
}

func (cors *CORS) Handle(*request.Request, *response.Response) error {
	return nil
}

func (cors *CORS) Docs(route *routes.Route) {
	panic("not implemented")
}

// func AllowedHeaders(headers ...string) engi.Middleware {
// 	return func(route *routes.Route) {
// 		routes.AllowedHeaders(headers...)
// 	}
// }

// func AllowedMethods(methods ...string) engi.Middleware {
// 	return func(route *routes.Route) {
// 		routes.AllowedMethods(methods...)
// 	}
// }

// func AllowedOrigins(origins ...string) engi.Middleware {
// 	return func(route *routes.Route) {
// 		routes.AllowedOrigins(origins...)
// 	}
// }
