package webapi

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KlyuchnikovV/webapi/internal/types"
	"github.com/KlyuchnikovV/webapi/logger"
	"github.com/KlyuchnikovV/webapi/placing"
	"github.com/KlyuchnikovV/webapi/response"
)

// TODO: add checking length of request from comments about field length

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

	log *logger.Log
}

func New(address string, configs ...Option) *Engine {
	if address == "" {
		address = defaultAddress
	}

	var engine = &Engine{
		apiPrefix:         defaultPrefix,
		log:               logger.New(nil),
		responseObject:    new(response.AsIs),
		responseMarshaler: *types.NewJSONMarshaler(),
		server: &http.Server{
			Addr:              address,
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
		var servicePath = fmt.Sprintf("/%s/%s/", e.apiPrefix, srv.Prefix())

		e.services[i] = NewService(e, srv, servicePath)

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
	e.log.Infof("Starting on %s", e.server.Addr)
	e.log.Infof("WebApi started...")

	return e.server.ListenAndServe()
}

type Context interface {
	// Request.

	// Headers - returns request headers.
	Headers() map[string][]string
	// All - returns all parsed parameters.
	All() map[placing.Placing]map[string]string
	// GetParameter - returns parameter value from defined place.
	GetParameter(string, placing.Placing) string
	// GetRequest - return http.Request object associated with request.
	GetRequest() *http.Request
	// Body - returns request body.
	// Body must be requested by 'api.Body(pointer)' or 'api.CustomBody(unmarshaler, pointer)'.
	Body() interface{}
	// Bool - returns boolean parameter.
	// Mandatory parameter should be requested by 'api.Bool'.
	// Otherwise, parameter will be obtained by key and its value will be checked for truth.
	Bool(string, placing.Placing) bool
	// Integer - returns integer parameter.
	// Mandatory parameter should be requested by 'api.Integer'.
	// Otherwise, parameter will be obtained by key and its value will be converted. to int64.
	Integer(string, placing.Placing) int64
	// Float - returns floating point number parameter.
	// Mandatory parameter should be requested by 'api.Float'.
	// Otherwise, parameter will be obtained by key and its value will be converted to float64.
	Float(string, placing.Placing) float64
	// String - returns String parameter.
	// Mandatory parameter should be requested by 'api.String'.
	// Otherwise, parameter will be obtained by key.
	String(string, placing.Placing) string
	// Time - returns date-time parameter.
	// Mandatory parameter should be requested by 'api.Time'.
	// Otherwise, parameter will be obtained by key and its value will be converted to time using 'layout'.
	Time(key string, layout string, paramPlacing placing.Placing) time.Time

	// Responses.

	// ResponseWriter - returns http.ResponseWriter associated with request.
	ResponseWriter() http.ResponseWriter
	// Object - responses with provided custom code and body.
	// Body will be marshaled using service-defined object and marshaler.
	Object(code int, payload interface{}) error
	// WithourContent - responses with provided custom code and no body.
	WithoutContent(code int) error
	// Error - responses custom error with provided code and formatted string message.
	Error(code int, format string, args ...interface{}) error
	// OK - writes payload into json's 'result' field with 200 http code.
	OK(payload interface{}) error
	// Created - responses with 201 http code and no content.
	Created() error
	// NoContent - responses with 204 http code and no content.
	NoContent() error
	// BadRequest - responses with 400 code and provided formatted string message.
	BadRequest(format string, args ...interface{}) error
	// Forbidden - responses with 403 error code and provided formatted string message.
	Forbidden(format string, args ...interface{}) error
	// NotFound - responses with 404 error code and provided formatted string message.
	NotFound(format string, args ...interface{}) error
	// MethodNotAllowed - responses with 405 error code and provided formatted string message.
	MethodNotAllowed(format string, args ...interface{}) error
	// InternalServerError - responses with 500 error code and provided formatted string message.
	InternalServerError(format string, args ...interface{}) error
}

type Route func(Context) error
