package webapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// TODO: errors and errors chan
// TODO: logging
// TODO: ServiceAPI methods refactor

type ServiceAPI interface {
	// PathPrefix - prefix of all paths for this service.
	PathPrefix() string
	// Routers returns the handlers and their relative paths (relative to the service) for registration.
	Routers() map[string]RouterFunc

	// GET - implements GET api method call.
	GET(HandlerFunc, ...HandlerFunc) RouterFunc
	// PUT - implements PUT api method call.
	PUT(HandlerFunc, ...HandlerFunc) RouterFunc
	// HEAD - implements HEAD api method call.
	HEAD(HandlerFunc, ...HandlerFunc) RouterFunc
	// POST - implements POST api method call.
	POST(HandlerFunc, ...HandlerFunc) RouterFunc
	// PATCH - implements PATCH api method call.
	PATCH(HandlerFunc, ...HandlerFunc) RouterFunc
	// TRACE - implements TRACE api method call.
	TRACE(HandlerFunc, ...HandlerFunc) RouterFunc
	// DELETE - implements DELETE api method call.
	DELETE(HandlerFunc, ...HandlerFunc) RouterFunc
	// CONNECT - implements CONNECT api method call.
	CONNECT(HandlerFunc, ...HandlerFunc) RouterFunc
	// OPTIONS - implements OPTIONS api method call.
	OPTIONS(HandlerFunc, ...HandlerFunc) RouterFunc

	// WithBody - takes pointer to structure and saves casted request body into context.
	// Result can be retrieved from context using 'context.QueryParams.Body()'.
	WithBody(interface{}) HandlerFunc
	// WithBool - queries mandatory boolean parameter from request by 'key'.
	// Result can be retrieved from context using 'context.QueryParams.Bool(key)'.
	WithBool(key string) HandlerFunc
	// WithInt - queries mandatory integer parameter from request by 'key'.
	// Result can be retrieved from context using 'context.QueryParams.Integer(key)'.
	WithInt(key string) HandlerFunc
	// WithFloat - queries mandatory floating point number parameter from request by 'key'.
	// Result can be retrieved from context using 'context.QueryParams.Float(key)'.
	WithFloat(key string) HandlerFunc
	// WithString - queries mandatory string parameter from request by 'key'.
	// Result can be retrieved from context using 'context.QueryParams.String(key)'.
	WithString(key string) HandlerFunc
	// WithTime - queries mandatory time parameter from request by 'key' using 'layout'.
	// Result can be retrieved from context using 'context.QueryParams.Time(key, layout)'.
	WithTime(key, layout string) HandlerFunc

	ServeHTTP(w http.ResponseWriter, r *http.Request)

	bind(ResponseMarshaler, Responser)
}

type ResponseMarshaler func(interface{}) ([]byte, error)

type Engine struct {
	// *echo.Echo
	mux *http.ServeMux

	responseBinder ResponseMarshaler
	responseObject Responser

	services []ServiceAPI
}

func New(configs ...func(*Engine)) *Engine {
	e := &Engine{
		responseBinder: json.Marshal,
		responseObject: new(AsIsObject),
	}

	for _, config := range configs {
		config(e)
	}

	return e
}

// RegisterServices - registering service routes.
func (e *Engine) RegisterServices(services ...ServiceAPI) error {
	e.services = services

	e.mux = http.NewServeMux()

	var prefix = "api"

	for i := range e.services {
		e.services[i].bind(e.responseBinder, e.responseObject)

		var servicePrefix = fmt.Sprintf("/%s/%s/",
			strings.Trim(prefix, "/"),
			strings.Trim(e.services[i].PathPrefix(), "/"),
		)

		for path, register := range e.services[i].Routers() {
			register(fmt.Sprintf("%s%s",
				servicePrefix,
				strings.Trim(path, "/"),
			))
		}

		e.mux.Handle(
			servicePrefix,
			e.services[i],
		)
	}

	return nil
}

func (e *Engine) Start(address string) error {
	log.Print("WebApi started...")

	return http.ListenAndServe(address, e.mux)
}

func (e *Engine) ResponseAsJSON() {
	e.responseBinder = json.Marshal
}

func (e *Engine) ResponseAsXML() {
	e.responseBinder = xml.Marshal
}

func (e *Engine) AsIsResponse() {
	e.responseObject = new(AsIsObject)
}

func (e *Engine) ObjectResponse(object Responser) {
	e.responseObject = object
}
