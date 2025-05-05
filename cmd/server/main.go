package main

import (
	"github.com/ZaneH/keep-talking/internal/application/services"
)

func main() {
	g := services.NewGameService()

	g.Start()
}
