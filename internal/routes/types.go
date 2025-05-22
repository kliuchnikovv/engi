package routes

import (
	"context"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
)

type NamedParameter interface {
	Name() string
	Regexp() string
}

type (
	Middleware interface {
		// Bind(*Route) error
		Handle(context.Context, *request.Request, *response.Response) error
		Docs(*Route)
		Priority() int
	}
)

func contains(slice []string, item string) bool {
	if len(slice) == 0 {
		return true
	}

	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}
