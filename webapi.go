package webapi

import (
	"fmt"
	"strings"

	"github.com/KlyuchnikovV/webapi/api"
	"github.com/labstack/echo"
)

type Engine struct {
	*echo.Echo

	services []api.ServiceAPI
}

func New() *Engine {
	return &Engine{
		Echo: echo.New(),
	}
}

func (e *Engine) RegisterServices(services ...api.ServiceAPI) error {
	e.services = services

	var r = e.Echo.Group("/api")

	for i := range e.services {
		if err := e.registerHandlers(e.services[i], r); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) registerHandlers(service api.ServiceAPI, r *echo.Group) error {
	var group = r.Group(
		fmt.Sprintf(
			"/%s", strings.Trim(service.PathPrefix(), "/"),
		),
	)

	for path, register := range service.Routers() {
		register(group, fmt.Sprintf(
			"/%s", strings.Trim(path, "/"),
		))
	}

	return nil
}

func (e *Engine) Start(address string) error {
	return e.Echo.Start(address)
}

// TODO: docs
// TODO: context rework
