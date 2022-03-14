package main

import (
	"log"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/service"
)

func main() {
	w := webapi.New()

	err := w.RegisterServices(
		service.NewRequestAPI(),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Start("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
