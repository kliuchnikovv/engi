package main

import (
	"log"

	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/services"
	"github.com/KlyuchnikovV/webapi/types"
)

func main() {
	w := webapi.New(
		":8080",

		// Equals to 'w.ResponseAsJSON()'
		// This form can be used with any engine methods having no parameters.
		(*webapi.Engine).ResponseAsJSON,
	)

	w.WithPrefix("api")
	w.ObjectResponse(new(types.AsIsResponse))

	err := w.RegisterServices(
		services.NewRequestAPI(w),
		services.NewNotesAPI(w),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Start()
	if err != nil {
		log.Fatal(err)
	}
}
