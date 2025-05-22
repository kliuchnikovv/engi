package cors

import (
	"errors"

	"github.com/kliuchnikovv/engi"
)

var (
	ErrOriginNotAllowed           = errors.New("origin is not allowed")
	ErrMissingConanicalHeader     = errors.New("missing canonical header")
	ErrCORSMethodHeaderNotFound   = errors.New("CORS-Method header not found")
	ErrCORSMethodHeaderNotAllowed = errors.New("CORS-Method header not allowed")
)

func AllowedHeaders(headers ...string) engi.Middleware {
	return corsAllowedHeaders(headers)
}

func AllowedMethods(methods ...string) engi.Middleware {
	return corsAllowedMethods(methods)
}

func AllowedOrigins(origins ...string) engi.Middleware {
	return corsAllowedOrigins(origins)
}
