package middlewares

import (
	"context"

	"github.com/kliuchnikovv/engi/internal/request"
	"github.com/kliuchnikovv/engi/internal/response"
	"github.com/kliuchnikovv/engi/internal/routes"
)

type description string

func Description(desc string) routes.Middleware {
	return description(desc)
}

func (description) Handle(_ context.Context, _ *request.Request, _ *response.Response) error {
	return nil
}

func (description) Docs(*routes.Route) {
	panic("unimplemented")
}

func (description) Priority() int {
	return 100 // TODO: make external priority map
}
