package main

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
