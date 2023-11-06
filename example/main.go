package main

import (
	"log"

	"github.com/KlyuchnikovV/engi"
	"github.com/KlyuchnikovV/engi/example/services"
	"github.com/KlyuchnikovV/engi/response"
)

func main() {
	w := engi.New(
		":8080",
		engi.WithPrefix("api"),
		engi.ResponseAsJSON(new(response.AsIs)),
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
