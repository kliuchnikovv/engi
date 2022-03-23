package main

import (
	"log"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/service"
)

func main() {
	w := webapi.New(
		":8080",
		// (*webapi.Engine).ResponseAsXML,
		(*webapi.Engine).ResponseAsJSON,
	)

	// w.ObjectResponse(new(webapi.ResponseObject))

	err := w.RegisterServices(
		service.NewRequestAPI(w),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		log.Fatal(err)
	}
}
