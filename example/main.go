package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/api/response"
	"github.com/KlyuchnikovV/engi/example/services"
)

func main() {
	w := engi.New(
		":8080",
		engi.WithPrefix("api"),
		engi.ResponseAsJSON(new(response.AsIs)),
		engi.WithLogger(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)),
	)

	if err := w.RegisterServices(
		new(services.NotesAPI),
		new(services.RequestAPI),
	); err != nil {
		log.Fatal(err)
	}

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}
