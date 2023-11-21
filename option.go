package engi

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/engi/internal/types"
	"github.com/KlyuchnikovV/engi/response"
)

type Option func(*Engine)

// WithResponse - tells server to use object as wrapper for all responses.
func WithResponse(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
	}
}

// AsIsResponse - tells server to response objects without wrapping.
func AsIsResponse(engine *Engine) {
	engine.responseObject = new(response.AsIs)
}

// Use - sets custom configuration function for http.Server.
func Use(f func(*http.Server)) Option {
	return func(engine *Engine) {
		f(engine.server)
	}
}

// WithPrefix - sets api's prefix.
func WithPrefix(prefix string) Option {
	return func(engine *Engine) {
		if prefix == "" {
			engine.apiPrefix = ""
		} else {
			engine.apiPrefix = fmt.Sprintf("/%s", strings.Trim(prefix, "/"))
		}
	}
}

// WithLogger - sets custom logger.
func WithLogger(handler slog.Handler) Option {
	return func(engine *Engine) {
		engine.logger = slog.New(handler)
	}
}

// ResponseAsJSON - tells server to serialize responses as JSON using object as wrapper.
func ResponseAsJSON(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
		engine.responseMarshaler = *types.NewJSONMarshaler()
	}
}

// ResponseAsXML - tells server to serialize responses as XML using object as wrapper.
func ResponseAsXML(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
		engine.responseMarshaler = *types.NewXMLMarshaler()
	}
}
