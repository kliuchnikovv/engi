package webapi

import (
	"github.com/KlyuchnikovV/webapi/api"
	"github.com/labstack/echo"
)

type Engine struct {
	*echo.Echo

	services []api.API
}

func New() *Engine {
	return &Engine{
		Echo: echo.New(),
	}
}

func (e *Engine) RegisterServices(services ...api.API) error {
	e.services = services

	var r = e.Echo.Group("/api")

	for i := range e.services {
		e.services[i].Bind(e.services[i].Routers())
		if err := e.services[i].RegisterHandlers(r); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Start(address string) error {
	return e.Echo.Start(address)
}

// TODO: docs
// TODO: readme
// TODO: context rework
