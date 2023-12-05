package main

import (
	"log"

	"github.com/ayo-ajayi/edutech/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file ", err.Error())
	}
	app.NewApp(":8000", app.Router()).Start()
}
