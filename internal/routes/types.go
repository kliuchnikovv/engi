package routes

import (
	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/internal/response"
)

type NamedParameter interface {
	Name() string
	Regexp() string
}

type Option interface {
	Bind(*Route) error
	Handle(*request.Request, *response.Response) error
	Docs(*Route)
}

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
