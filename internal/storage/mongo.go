package storage

import (
	"context"
	"time"
	"weather/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDBClient(client *mongo.Client) *MongoDB {
	collection := client.Database("weatherDB").Collection("forecast_cache")
	return &MongoDB{collection: collection}
}

func (c *MongoDB) GetCachedWeather(city, date, api string) (*models.WeatherCache, error) {
	var result models.WeatherCache
	filter := bson.M{"city": city, "date": date, "api": api}

	err := c.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *MongoDB) CacheWeather(city, date, api, temperature string) error {
	weatherCache := models.WeatherCache{
		City:        city,
		Date:        date,
		API:         api,
		Temperature: temperature,
		RetrievedAt: time.Now(),
	}

	_, err := c.collection.InsertOne(context.TODO(), weatherCache)
	return err
}
