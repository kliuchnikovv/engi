package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/api/parameter/placing"
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
	"github.com/KlyuchnikovV/engi/internal/routes"
)

const (
	authHeader = "Authorization"

	bearerPrefix = "Bearer "
)

var errUnathorized = errors.New("Unauthorized.")

type Authorization struct {
	name string

	handle func(r *http.Request, w http.ResponseWriter) error
}

func (auth *Authorization) Bind(route *routes.Route) error {
	route.SetAuth(auth.handle)

	return nil
}

func (auth *Authorization) Handle(*request.Request, *response.Response) error {
	return nil
}

func (auth *Authorization) Docs(route *routes.Route) {
	panic("not implemented")
}

func NoAuth() engi.Middleware {
	return &Authorization{
		name: "no auth",
		handle: func(r *http.Request, w http.ResponseWriter) error {
			return nil
		},
	}
}

func Basic(username, password string) engi.Middleware {
	return &Authorization{
		name: "basic",
		handle: func(r *http.Request, w http.ResponseWriter) error {
			gotUser, gotPassword, ok := r.BasicAuth()
			if !ok {
				return errUnathorized
			}

			if username != gotUser || password != gotPassword {
				return errUnathorized
			}

			return nil
		},
	}
}

func BearerToken(
	token string,
) engi.Middleware {
	return BearerFunc(
		func(s string) bool { return s == token },
	)
}

func BearerFunc(
	isValid func(string) bool,
) engi.Middleware {
	return &Authorization{
		name: "bearer",
		handle: func(r *http.Request, w http.ResponseWriter) error {
			var header = r.Header.Get(authHeader)
			if len(header) == 0 {
				w.WriteHeader(http.StatusUnauthorized)

				return errUnathorized
			}

			if isValid(strings.TrimPrefix(authHeader, bearerPrefix)) {
				w.WriteHeader(http.StatusUnauthorized)

				return errUnathorized
			}

			return nil
		},
	}
}

// TODO: implement
func APIKey(
	key, value string, place placing.Placing,
) engi.Middleware {
	return &Authorization{
		name: "api key",
		handle: func(r *http.Request, w http.ResponseWriter) error {
			var param string

			switch place {
			case placing.InCookie:
				// TODO:
			case placing.InHeader:
			case placing.InPath:
			case placing.InQuery:
			default:
				return nil // TODO:
			}

			// var param = r.GetParameter(key, place)
			if len(param) == 0 {
				w.WriteHeader(http.StatusUnauthorized)

				return errUnathorized
			}

			if value != param {
				w.WriteHeader(http.StatusUnauthorized)

				return errUnathorized
			}

			return nil
		},
	}
}
