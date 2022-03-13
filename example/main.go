package main

import (
	"github.com/KlyuchnikovV/webapi"
	"github.com/KlyuchnikovV/webapi/example/service"
)

func main() {
	w := webapi.New()

	w.RegisterServices(
		service.NewRequestAPI(),
	)

	w.Start("localhost:8080")
}
