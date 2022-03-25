package webapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/KlyuchnikovV/webapi/types"
)

// Engine - server provider.
type Engine struct {
	apiPrefix string

	server http.Server

	responseMarshaler types.Marshaler
	responseObject    types.Responser

	services []ServiceAPI

	log *types.Log
}

func New(address string, configs ...func(*Engine)) *Engine {
	e := &Engine{
		apiPrefix: "api",
		server: http.Server{
			Addr: address,
		},
		log:               types.NewLog(nil),
		responseMarshaler: json.Marshal,
		responseObject:    new(types.AsIsResponse),
	}

	for _, config := range configs {
		config(e)
	}

	return e
}

// RegisterServices - registering service routes.
func (e *Engine) RegisterServices(services ...ServiceAPI) error {
	e.services = services

	var mux = http.NewServeMux()

	for i := range e.services {
		var servicePrefix = fmt.Sprintf("/%s/%s/", e.apiPrefix, e.services[i].PathPrefix())

		for path, register := range e.services[i].Routers() {
			register(fmt.Sprintf("%s%s", servicePrefix, strings.Trim(path, "/")))
		}

		mux.Handle(
			servicePrefix,
			e.services[i],
		)
	}

	e.server.Handler = mux

	return nil
}

// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// If srv.Addr is blank, ":http" is used.
//
// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
func (e *Engine) Start() error {
	e.log.Infof("WebApi started...")

	return e.server.ListenAndServe()
}

// Use - sets custom configuration function.
func (e *Engine) Use(f func(*http.Server)) {
	f(&e.server)
}

// WithPrefix - sets api's prefix.
func (e *Engine) WithPrefix(prefix string) {
	e.apiPrefix = strings.Trim(prefix, "/")
}

// ResponseAsJSON - tells server to serialize all responses as JSON.
func (e *Engine) ResponseAsJSON() {
	e.responseMarshaler = json.Marshal
}

// ResponseAsXML - tells server to serialize all responses as XML.
func (e *Engine) ResponseAsXML() {
	e.responseMarshaler = func(i interface{}) ([]byte, error) {
		bytes, err := xml.Marshal(i)
		if err != nil {
			return nil, err
		}

		// Should append header for proper visualization.
		return append([]byte(xml.Header), bytes...), nil
	}
}

// AsIsResponse - tells server to response objects without wrapping.
func (e *Engine) AsIsResponse() {
	e.responseObject = new(types.AsIsResponse)
}

// ObjectResponse - tells server to use object as wrapper for all responses.
func (e *Engine) ObjectResponse(object types.Responser) {
	e.responseObject = object
}

// WithLogger - sets custom logger.
func (e *Engine) WithLogger(log types.Logger) {
	e.log = types.NewLog(log)
}

// WithSendingErrors - sets errors channel capacity.
func (e *Engine) WithSendingErrors(capacity int) {
	if e.log == nil {
		e.log = types.NewLog(nil)
	} else {
		e.log.SetChannelCapacity(capacity)
	}
}
