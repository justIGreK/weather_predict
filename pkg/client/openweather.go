package client

import (
	"fmt"
	"net/http"
	"os"
)

type OpenWeatherClient struct{}

func NewOpenWeatherClient() *OpenWeatherClient {
	return &OpenWeatherClient{}
}

func (c *OpenWeatherClient) GetCoordinates(city string) (*http.Response, error) {
	apiKey := os.Getenv("OPENWEATHER")
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
