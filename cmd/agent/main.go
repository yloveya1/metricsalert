package main

import (
	"context"
	"log"

	"github.com/yloveya1/metricsalert/internal/app"
)

func main() {
	ctx := context.Background()
	err := app.RunAgent(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
