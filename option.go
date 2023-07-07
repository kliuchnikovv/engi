package webapi

import (
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/webapi/internal/types"
	"github.com/KlyuchnikovV/webapi/logger"
	"github.com/KlyuchnikovV/webapi/response"
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
		engine.apiPrefix = strings.Trim(prefix, "/")
	}
}

// WithLogger - sets custom logger.
func WithLogger(log logger.Logger) Option {
	return func(engine *Engine) {
		engine.log = logger.New(log)
	}
}

// WithSendingErrors - sets errors channel capacity.
func WithSendingErrors(capacity int) Option {
	return func(engine *Engine) {
		if engine.log == nil {
			engine.log = logger.New(nil)
		} else {
			engine.log.SetChannelCapacity(capacity)
		}
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
