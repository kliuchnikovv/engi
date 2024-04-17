package routes

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/api/response"
	"github.com/KlyuchnikovV/engi/internal/request"
)

func noOpMiddleware(
	*request.Request,
	http.ResponseWriter,
) *response.AsObject {
	return nil
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
