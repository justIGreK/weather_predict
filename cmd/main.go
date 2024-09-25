package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"weather/cmd/handler"
	"weather/internal/service"
	"weather/internal/storage"
	"weather/pkg/client"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loadEnv()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	weatherAPIClient := client.NewWeatherAPIClient()
	tomorrowIOClient := client.NewTomorrowIOClient()
	openWeatherClient := client.NewOpenWeatherClient()
	mongoDB := storage.NewMongoDBClient(mongoClient)

	weatherService := service.NewWeatherService(weatherAPIClient, tomorrowIOClient, openWeatherClient, mongoDB)
	weatherRouter := handler.NewHandler(weatherService)

	err = http.ListenAndServe(":7777", weatherRouter.InitRouters())
	if err != nil {
		log.Fatal("error during running server: ", err)
	}

}
