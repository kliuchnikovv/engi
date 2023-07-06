package webapi

import (
	"encoding/json"
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
	if address == "" {
		address = defaultAddress
	}

	var engine = &Engine{
		apiPrefix:         defaultPrefix,
		log:               types.NewLog(nil),
		responseObject:    new(types.AsIsResponse),
		responseMarshaler: json.Marshal,
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
	e.log.Infof("Starting on %s", e.server.Addr)
	e.log.Infof("WebApi started...")

	return e.server.ListenAndServe()
}
