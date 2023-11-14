package engi

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KlyuchnikovV/engi/internal/types"
	"github.com/KlyuchnikovV/engi/response"
)

// TODO: add checking length of request from comments about field length
// TODO: authorization
// TODO: string builder
// TODO: benchmarks
// TODO: tests
// TODO: rename
// TODO: work with context
// TODO: logging (log url usages)

const (
	defaultPrefix  = "/api"
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
		apiPrefix:         defaultPrefix,
		responseObject:    new(response.AsIs),
		responseMarshaler: *types.NewJSONMarshaler(),
		logger:            slog.New(slog.NewTextHandler(os.Stdout, nil)),
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
		var servicePath = fmt.Sprintf("%s/%s/", e.apiPrefix, srv.Prefix())

		e.services[i] = NewService(e, srv, servicePath)

		for _, middleware := range e.services[i].Middlewares() {
			middleware(e.services[i])
		}

		for path, register := range srv.Routers() {
			if err := register(e.services[i], strings.Trim(path, "/")); err != nil {
				return fmt.Errorf("%w, engine: %s", err, strings.Trim(e.apiPrefix, "/"))
			}
		}

		mux.HandleFunc(servicePath, e.handle(i))
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

func (e *Engine) handle(i int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uri, _ := strings.CutPrefix(r.URL.Path, e.apiPrefix)
		e.services[i].Serve(uri, r, w)
	}
}
