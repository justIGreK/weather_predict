package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type TomorrowIOClient struct{}

func NewTomorrowIOClient() *TomorrowIOClient {
	return &TomorrowIOClient{}
}

func (c *TomorrowIOClient) GetWeather(lat, lon, startTime, endTime string) (*http.Response, error) {
	log.Println(startTime, endTime)
	startTime = fmt.Sprintf("%sT00:00:00Z", startTime)
	endTime = fmt.Sprintf("%sT01:00:00Z", endTime)
	apiKey := os.Getenv("TOMORROWIO")
	url := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%s,%s&fields=temperature&units=metric&timesteps=1d&startTime=%s&endTime=%s&apikey=%s",
		lat, lon, startTime, endTime, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
