package webapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KlyuchnikovV/webapi/types"
)

const (
	defaultPrefix  = "api"
	defaultAddress = ":8080"
	defaultTimeout = 5 * time.Second
)

// Engine - server provider.
type Engine struct {
	apiPrefix string

	server *http.Server

	responseMarshaler types.Marshaler
	responseObject    types.Responser

	services []*Service

	log *types.Log
}

func New(address string, configs ...Option) *Engine {
	var engine = &Engine{
		apiPrefix:         defaultPrefix,
		log:               types.NewLog(nil),
		responseObject:    new(types.AsIsResponse),
		responseMarshaler: json.Marshal,
		server: &http.Server{
			Addr:              defaultAddress,
			ReadTimeout:       defaultTimeout,
			WriteTimeout:      defaultTimeout,
			IdleTimeout:       defaultTimeout,
			ReadHeaderTimeout: defaultTimeout,
		},
	}

	for _, config := range configs {
		config(engine)
	}

	return engine
}

// RegisterServices - registering service routes.
func (e *Engine) RegisterServices(services ...ServiceAPI) error {
	e.services = make([]*Service, len(services))

	var mux = http.NewServeMux()

	for i, srv := range services {
		e.services[i] = NewService(e, srv)

		var servicePath = fmt.Sprintf("/%s/%s/", e.apiPrefix, srv.Prefix())

		if middlewares, ok := srv.(MiddlewaresAPI); ok {
			for _, middleware := range middlewares.Middlewares() {
				middleware(e.services[i])
			}
		}

		for path, register := range srv.Routers() {
			register(e.services[i], fmt.Sprintf("%s%s", servicePath, strings.Trim(path, "/")))
		}

		mux.Handle(servicePath, e.services[i])
	}

	e.server.Handler = mux

	return nil
}

// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
func (e *Engine) Start() error {
	e.log.Info("Starting on %s", e.server.Addr)
	e.log.Info("WebApi started...")

	return e.server.ListenAndServe()
}

type Option func(*Engine)

// WithResponse - tells server to use object as wrapper for all responses.
func WithResponse(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
	}
}

// AsIsResponse - tells server to response objects without wrapping.
func AsIsResponse(engine *Engine) {
	engine.responseObject = new(types.AsIsResponse)
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
func WithLogger(log types.Logger) Option {
	return func(engine *Engine) {
		engine.log = types.NewLog(log)
	}
}

// WithSendingErrors - sets errors channel capacity.
func WithSendingErrors(capacity int) Option {
	return func(engine *Engine) {
		if engine.log == nil {
			engine.log = types.NewLog(nil)
		} else {
			engine.log.SetChannelCapacity(capacity)
		}
	}
}

// ResponseAsJSON - tells server to serialize responses as JSON using object as wrapper.
func ResponseAsJSON(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
		engine.responseMarshaler = json.Marshal
	}
}

// ResponseAsXML - tells server to serialize responses as XML using object as wrapper.
func ResponseAsXML(object types.Responser) Option {
	return func(engine *Engine) {
		engine.responseObject = object
		engine.responseMarshaler = func(i interface{}) ([]byte, error) {
			bytes, err := xml.Marshal(i)
			if err != nil {
				return nil, err
			}

			// Should append header for proper visualization.
			return append([]byte(xml.Header), bytes...), nil
		}
	}
}
