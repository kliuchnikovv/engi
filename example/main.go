package main

import (
	"log"

	webapi "github.com/KlyuchnikovV/webapi/api"
	"github.com/KlyuchnikovV/webapi/api/example/service"
)

func main() {
	w := webapi.New(
		(*webapi.Engine).ResponseAsXML,
	)

	w.ObjectResponse(new(webapi.ResponseObject))

	err := w.RegisterServices(
		service.NewRequestAPI(),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
