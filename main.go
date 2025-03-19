package main

import (
	"context"
	"log"

	"github.com/MaksimovDenis/loadinator2000/internal/app"
)

func main() {
	ctx := context.Background()

	loadinator2000, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = loadinator2000.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
