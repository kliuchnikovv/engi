package engi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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
		path, err := url.JoinPath(e.apiPrefix, service.Prefix())
		if err != nil {
			return fmt.Errorf("failed to join path: %w", err)
		}

		var srv = NewService(e, service, path)
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

		e.services[i] = srv

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

// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
// Start listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Start always returns a non-nil error. After Shutdown or Close, the returned error is ErrServerClosed.
func (e *Engine) Start() error {
	e.logger.Info("Starting engi", slog.String("address", e.server.Addr))

	// Run shutdown in background
	go e.shutdownEngine()

	// Listen for OS signals
	signal.Notify(e.signalChan, os.Interrupt, syscall.SIGTERM)

	if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		e.logger.Error("Server error", slog.String("error", err.Error()))

		return err
	}

	return nil
}

// Shutdown allows manual graceful shutdown with a context.
func (e *Engine) Shutdown(ctx context.Context) {
	e.logger.Info("manual shutdown initiated")

	e.signalChan <- syscall.SIGTERM
}

func (e *Engine) shutdownEngine() {
	<-e.signalChan

	e.logger.Info("shutting down engi")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.server.Shutdown(ctx); err != nil {
		e.logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
	}

	e.logger.Info("engi stopped gracefully")
}
