package client

import (
	"fmt"
	"net/http"
	"os"
)

type WeatherAPIClient struct{}

func NewWeatherAPIClient() *WeatherAPIClient {
	return &WeatherAPIClient{}
}

func (c *WeatherAPIClient) GetWeather(lat, lon, date string) (*http.Response, error) {
	apiKey := os.Getenv("WEATHERAPI")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s,%s&dt=%s", apiKey, lat, lon, date)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
