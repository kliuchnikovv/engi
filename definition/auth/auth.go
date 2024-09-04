package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/definition/parameter/placing"
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

func APIKey(
	key, value string, place placing.Placing,
) engi.Middleware {
	if place == placing.InPath {
		panic("placing api key in path not supported")
	}

	return &Authorization{
		name: "api key",
		handle: func(r *http.Request, w http.ResponseWriter) error {
			var parameter string

			switch place {
			case placing.InHeader:
				parameter = r.Header.Get(key)
			case placing.InQuery:
				parameter = strings.Join(r.URL.Query()[key], "")
			case placing.InCookie:
				for _, cookie := range r.Cookies() {
					if cookie.Name == key {
						parameter = cookie.Value
						break
					}
				}
			}

			if value != parameter {
				w.WriteHeader(http.StatusUnauthorized)

				return errUnathorized
			}

			return nil
		},
	}
}
