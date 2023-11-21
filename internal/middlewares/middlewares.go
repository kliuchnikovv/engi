package middlewares

import (
	"net/http"

	"github.com/KlyuchnikovV/engi/internal/request"
	"github.com/KlyuchnikovV/engi/response"
)

func noOpMiddleware(*request.Request, http.ResponseWriter) *response.AsObject {
	return nil
}

type Register func(middlewares *Middlewares)

type Middlewares struct {
	cors   request.Middleware
	auth   request.Middleware
	params []request.Middleware
	other  []request.Middleware
}

func New(registrators ...Register) *Middlewares {
	var middlewares = &Middlewares{
		cors:   noOpMiddleware,
		auth:   noOpMiddleware,
		params: make([]request.Middleware, 0),
		other:  make([]request.Middleware, 0),
	}

	for _, register := range registrators {
		register(middlewares)
	}

	return middlewares
}

func (m *Middlewares) Add(registrators ...Register) {
	for _, register := range registrators {
		register(m)
	}
}

func (m *Middlewares) AddParams(middlewares ...request.Middleware) {
	if m.params == nil {
		m.params = make([]request.Middleware, 0, len(middlewares))
	}

	m.params = append(m.params, middlewares...)
}

func (m *Middlewares) AddAuth(middleware request.Middleware) {
	m.auth = middleware
}

func (m *Middlewares) AddCORS(middleware request.Middleware) {
	m.cors = middleware
}

func (m *Middlewares) AddOther(middlewares ...request.Middleware) {
	if m.other == nil {
		m.other = make([]request.Middleware, 0, len(middlewares))
	}

	m.other = append(m.other, middlewares...)
}

func (m *Middlewares) Handle(r *request.Request, w http.ResponseWriter) *response.AsObject {
	if err := m.cors(r, w); err != nil {
		return err
	}

	if err := m.auth(r, w); err != nil {
		return err
	}

	for _, param := range m.params {
		if err := param(r, w); err != nil {
			return err
		}
	}

	for _, other := range m.other {
		if err := other(r, w); err != nil {
			return err
		}
	}

	return nil
}
