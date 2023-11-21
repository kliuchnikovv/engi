package engi

import (
	"github.com/KlyuchnikovV/engi/internal/middlewares"
	"github.com/KlyuchnikovV/engi/internal/middlewares/auth"
	"github.com/KlyuchnikovV/engi/internal/middlewares/cors"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
)

type (
	Middleware  func(*Service)
	Middlewares []Middleware

	Register middlewares.Register
)

var (
	AllowedHeaders = cors.AllowedHeaders
	AllowedMethods = cors.AllowedMethods
	AllowedOrigins = cors.AllowedOrigins
)

func UseCORS(opts ...cors.CORSOption) Register {
	var cors = new(cors.CORS)
	for _, opt := range opts {
		opt(cors)
	}

	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddCORS(cors.Handle)
	}
}

var (
	NoAuth     = auth.NoAuth
	BasicAuth  = auth.Basic
	BearerAuth = auth.Bearer // TODO: remake it
	APIKeyAuth = func(key, value string, place AuthKeyPlacing) request.Middleware {
		return auth.APIKey(key, value, placing.Placing(place))
	}
)

type AuthKeyPlacing placing.Placing

const (
	InQuery  AuthKeyPlacing = AuthKeyPlacing(placing.InQuery)
	InCookie AuthKeyPlacing = AuthKeyPlacing(placing.InCookie)
	InHeader AuthKeyPlacing = AuthKeyPlacing(placing.InHeader)
)

func UseAuthorization(option request.Middleware) Register {
	return func(middlewares *middlewares.Middlewares) {
		middlewares.AddAuth(option)
	}
}
