package main

import (
	"log"

	"github.com/yloveya1/metricsalert/internal/app"
)

func main() {
	err := app.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
