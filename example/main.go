package main

import (
	"log"

	"github.com/KlyuchnikovV/webapi"
	service "github.com/KlyuchnikovV/webapi/example/services"
	"github.com/KlyuchnikovV/webapi/types"
)

func main() {
	w := webapi.New(
		":8080",
		// (*webapi.Engine).ResponseAsXML,
		(*webapi.Engine).ResponseAsJSON,
	)

	// w.WithPrefix("server")

	w.ObjectResponse(new(types.AsIsResponse))

	err := w.RegisterServices(
		service.NewRequestAPI(w),
		service.NewNotesAPI(w),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		log.Fatal(err)
	}
}
