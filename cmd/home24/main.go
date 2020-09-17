package main

import (
	"github.com/donutloop/home24/internal/api"
)

func main() {
	a := api.NewAPI(false)
	a.Bootstrap()
	a.Start()
}
