package main

import (
	"log"
	"spy_cat/internal/pkg/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Could not initialize application: %v", err)
	}

	application.Run()
}
