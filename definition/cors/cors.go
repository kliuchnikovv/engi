package cors

import (
	"net/http"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
)

type CORS struct {
	allowedHeaders []string
	allowedMethods []string
	allowedOrigins []string

	handle func(r *http.Request, w http.ResponseWriter) error
}

func (cors *CORS) Bind(route *routes.Route) error {
	route.AllowedHeaders = append(route.AllowedHeaders, cors.allowedHeaders...)
	route.AllowedMethods = append(route.AllowedMethods, cors.allowedMethods...)
	route.AllowedOrigins = append(route.AllowedOrigins, cors.allowedOrigins...)

	return nil
}

func (cors *CORS) Handle(*request.Request, *response.Response) error {
	return nil
}

func (cors *CORS) Docs(route *routes.Route) {
	panic("not implemented")
}

func AllowedHeaders(headers ...string) engi.Middleware {
	return &CORS{
		allowedHeaders: headers,
	}
}

func AllowedMethods(methods ...string) engi.Middleware {
	return &CORS{
		allowedMethods: methods,
	}
}

func AllowedOrigins(origins ...string) engi.Middleware {
	return &CORS{
		allowedOrigins: origins,
	}
}
