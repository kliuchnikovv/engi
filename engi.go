package engi

import (
	"log/slog"
	"net/http"
	"os"
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
	services  []*Service

	responseMarshaler types.Marshaler
	responseObject    types.Responser

	server *http.Server
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
