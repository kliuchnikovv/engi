package auth

import (
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/parameter/placing"
	"github.com/KlyuchnikovV/engi/response"
)

const (
	authHeader = "Authorization"

	bearerPrefix = "Bearer "

	unathorizedResponse = "Unauthorized."
)

func NoAuth(*request.Request, http.ResponseWriter) *response.AsObject { return nil }

func Basic(username, password string) request.Middleware {
	return func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
		gotUser, gotPassword, ok := r.GetRequest().BasicAuth()
		if !ok {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		if username != gotUser || password != gotPassword {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		return nil
	}
}

func APIKey(key, value string, place placing.Placing) request.Middleware {
	return func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
		var param = r.GetParameter(key, place)
		if len(param) == 0 {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		if value != param {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		return nil
	}
}

func Bearer(isValid func(string) bool) request.Middleware {
	return func(r *request.Request, _ http.ResponseWriter) *response.AsObject {
		var header = r.GetRequest().Header.Get(authHeader)
		if len(header) == 0 {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		if isValid(strings.TrimPrefix(authHeader, bearerPrefix)) {
			return response.AsError(http.StatusUnauthorized, unathorizedResponse)
		}

		return nil
	}
}
