package engi

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KlyuchnikovV/engi/internal/types"
)

// TODO: add checking length of request from comments about field length
// TODO: benchmarks
// TODO: tests
// TODO: logging (log url usages)
// TODO: documentation

const (
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

	logger *slog.Logger
}

func New(address string, configs ...Option) *Engine {
	if address == "" {
		address = defaultAddress
	}

	var engine = &Engine{
		responseObject:    new(types.ResponseAsIs),
		responseMarshaler: types.NewJSONMarshaler(),
		server: &http.Server{
			Addr:              address,
			ReadTimeout:       defaultTimeout,
			WriteTimeout:      defaultTimeout,
			IdleTimeout:       defaultTimeout,
			ReadHeaderTimeout: defaultTimeout,
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	for _, config := range configs {
		config(engine)
	}

	return engine
}

// RegisterServices - registering service routes.
func (e *Engine) RegisterServices(services ...ServiceDefinition) error {
	e.services = make([]*Service, len(services))

	var mux = http.NewServeMux()

	for i, service := range services {
		var (
			servicePath = fmt.Sprintf("%s/%s/", e.apiPrefix, service.Prefix())
			srv         = NewService(e, service, servicePath)
		)

		e.services[i] = srv

		for path, register := range service.Routers() {
			register(
				e.services[i],
				strings.Trim(path, "/"),
			)

			e.services[i].logger.Debug("route registered",
				slog.String("path", path),
				slog.String("full_path", fmt.Sprintf("%s%s", servicePath, path)),
			)
		}

		mux.Handle(servicePath, e.services[i])

		e.logger.Debug("service registered", slog.String("service", service.Prefix()))
	}

	e.server.Handler = mux

	return nil
}

// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
func (e *Engine) Start() error {
	e.logger.Info("Starting engi", slog.String("address", e.server.Addr))
	e.logger.Info("engi started...")

	return e.server.ListenAndServe()
}
