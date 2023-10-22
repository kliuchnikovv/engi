package main

import (
	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/services"
	"github.com/KlyuchnikovV/webapi/response"
	"log"
)

func main() {
	w := webapi.New(
		":8080",
		webapi.WithPrefix("api"),
		webapi.ResponseAsJSON(new(response.AsIs)),
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
