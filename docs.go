//go:build !docs
// +build !docs

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

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// RegisterServices registers ServiceAPI implementations into the HTTP mux.
func (e *Engine) RegisterServices(services ...ServiceDefinition) error {
	e.services = make([]*Service, len(services))
	mux := http.NewServeMux()

	for i, service := range services {
		path := fmt.Sprintf("%s/%s/", e.apiPrefix, service.Prefix())
		srv := NewService(e, service, path)
		e.services[i] = srv

		for route, register := range service.Routers() {
			if err := register(srv, route.method, strings.Trim(route.path, "/")); err != nil {
				return fmt.Errorf("%w, engine: %s", err, strings.Trim(e.apiPrefix, "/"))
			}
			srv.logger.Debug("route registered",
				slog.String("method", route.method),
				slog.String("route", route.path),
				slog.String("full", path+route.path),
			)
		}

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if err := srv.Serve(w, r); err != nil {
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

// TODO: refactor
// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
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
