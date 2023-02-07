package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mewejo/go-watering/pkg/app"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	app.NewApp().Run()
}
