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
	// AddInPathParameter(key string, value string)
	GetParameter(string, placing.Placing) string
	Headers() map[string][]string
	All() map[string]string

	// Parameters
	GetRequest() *http.Request
	Body() interface{}
	Bool(string, placing.Placing) bool
	Float(string, placing.Placing) float64
	Integer(string, placing.Placing) int64
	String(key string, paramPlacing placing.Placing) string
	Time(key string, layout string, paramPlacing placing.Placing) time.Time

	// Responses
	GetResponse() http.ResponseWriter
	BadRequest(format string, args ...interface{}) error
	Created() error
	Error(code int, format string, args ...interface{}) error
	Forbidden(format string, args ...interface{}) error
	InternalServerError(format string, args ...interface{}) error
	JSON(code int, payload interface{}) error
	MethodNotAllowed(format string, args ...interface{}) error
	NoContent() error
	NotFound(format string, args ...interface{}) error
	OK(payload interface{}) error
	WithoutContent(code int) error
}

type Route func(Context) error
