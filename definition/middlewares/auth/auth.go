package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/kliuchnikovv/engi"
	"github.com/kliuchnikovv/engi/definition/parameter/placing"
	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
)

// TODO: add authorization middleware type

const (
	authHeader = "Authorization"

	bearerPrefix = "Bearer "
)

var errUnathorized = errors.New("Unauthorized.")

type Authorization struct {
	name string

	handle func(context.Context, *http.Request, http.ResponseWriter) error
}

func (auth *Authorization) Bind(route *routes.Route) error {
	// route.SetAuth(auth.handle)

	return nil
}

func (auth *Authorization) Handle(ctx context.Context, req *request.Request, resp *response.Response) error {
	return auth.handle(ctx, req.GetRequest(), resp.ResponseWriter())
}

func (auth *Authorization) Docs(route *routes.Route) {
	panic("not implemented")
}

func (auth *Authorization) Priority() int {
	return 20
}

func NoAuth() engi.Middleware {
	return &Authorization{
		name: "no auth",
		handle: func(_ context.Context, _ *http.Request, _ http.ResponseWriter) error {
			return nil
		},
	}
}

func Basic(username, password string) engi.Middleware {
	return &Authorization{
		name: "basic",
		handle: func(_ context.Context, r *http.Request, w http.ResponseWriter) error {
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
		handle: func(_ context.Context, r *http.Request, w http.ResponseWriter) error {
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
		handle: func(_ context.Context, r *http.Request, w http.ResponseWriter) error {
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
