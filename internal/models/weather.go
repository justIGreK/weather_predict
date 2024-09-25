package models

import "time"

type WeatherCache struct {
	City        string    `bson:"city"`
	Date        string    `bson:"date"`
	API         string    `bson:"api"`         
	Temperature string    `bson:"temperature"`  
	RetrievedAt time.Time `bson:"retrieved_at"` 
}
