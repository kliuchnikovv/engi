package main

import (
	"github.com/KlyuchnikovV/webapi"
)

func main() {
	w := webapi.New()

	w.RegisterServices(
		NewRequestAPI(),
	)

	w.Start("localhost:8080")
}
