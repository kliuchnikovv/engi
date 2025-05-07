package engi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/KlyuchnikovV/engi/internal/types"
	"github.com/KlyuchnikovV/engi/response"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultAddress = ":8080"
	defaultTimeout = 5 * time.Second
)

// Engine - server provider.
type Engine struct {
	apiPrefix         string
	server            *http.Server
	responseMarshaler types.Marshaler
	responseObject    types.Responser
	services          []*Service
	logger            *slog.Logger
	tracerProvider    trace.TracerProvider
}

// New initializes a new Engine with the given address and options.
func New(address string, configs ...Option) *Engine {
	if address == "" {
		address = defaultAddress
	}

	engine := &Engine{
		responseObject:    new(response.AsIs),
		responseMarshaler: *types.NewJSONMarshaler(),
		server: &http.Server{
			Addr:              address,
			ReadTimeout:       defaultTimeout,
			WriteTimeout:      defaultTimeout,
			IdleTimeout:       defaultTimeout,
			ReadHeaderTimeout: defaultTimeout,
		},
		logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
		tracerProvider: otel.GetTracerProvider(),
	}

	for _, config := range configs {
		config(engine)
	}

	return engine
}

// RegisterServices registers ServiceAPI implementations into the HTTP mux.
func (e *Engine) RegisterServices(services ...ServiceAPI) error {
	e.services = make([]*Service, len(services))
	mux := http.NewServeMux()

	for i, service := range services {
		path := fmt.Sprintf("%s/%s/", e.apiPrefix, service.Prefix())
		srv := NewService(e, service, path)
		e.services[i] = srv

		for route, register := range service.Routers() {
			if err := register(srv, strings.Trim(route, "/")); err != nil {
				return fmt.Errorf("%w, engine: %s", err, strings.Trim(e.apiPrefix, "/"))
			}
			srv.logger.Debug("route registered", slog.String("route", route), slog.String("full", path+route))
		}

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			uri, _ := strings.CutPrefix(r.URL.Path, fmt.Sprintf("%s/%s", e.apiPrefix, srv.api.Prefix()))
			if err := srv.Serve(w, r, uri); err != nil {
				srv.logger.Error(err.Error())
			} else {
				srv.logger.Debug("request handled")
			}
		})
		e.logger.Debug("service registered", slog.String("service", service.Prefix()))
	}

	// Wrap mux with OpenTelemetry instrumentation
	e.server.Handler = otelhttp.NewHandler(mux, fmt.Sprintf("engi-server:%s", e.apiPrefix))
	return nil
}

// Start runs the HTTP server and handles graceful shutdown on SIGINT/SIGTERM.
func (e *Engine) Start() error {
	e.logger.Info("Starting engi", slog.String("address", e.server.Addr))

	// Run server in background
	go func() {
		if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.logger.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	// Listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.logger.Info("Shutting down engi")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.server.Shutdown(ctx); err != nil {
		e.logger.Error("Graceful shutdown failed", slog.String("error", err.Error()))
		return err
	}
	e.logger.Info("engi stopped gracefully")
	return nil
}

// Shutdown allows manual graceful shutdown with a context.
func (e *Engine) Shutdown(ctx context.Context) error {
	e.logger.Info("Manual shutdown initiated")
	return e.server.Shutdown(ctx)
}
