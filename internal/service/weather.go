package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
	"weather/internal/storage"
	"weather/pkg/client"
)

type WeatherService struct {
	weatherAPIClient  *client.WeatherAPIClient
	tomorrowIOClient  *client.TomorrowIOClient
	openWeatherClient *client.OpenWeatherClient
	mongoDB           *storage.MongoDB
}

func NewWeatherService(weatherAPIClient *client.WeatherAPIClient, tomorrowIOClient *client.TomorrowIOClient, openWeatherClient *client.OpenWeatherClient, mongoDB *storage.MongoDB) *WeatherService {
	return &WeatherService{
		weatherAPIClient:  weatherAPIClient,
		tomorrowIOClient:  tomorrowIOClient,
		openWeatherClient: openWeatherClient,
		mongoDB:           mongoDB,
	}
}

func (s *WeatherService) GetWeatherSummary(city, date string) (map[string]string, error) {
	summary := map[string]string{}
	err := s.CheckForAvailibleDate(date, 15)
	if err != nil {
		return nil, err
	}
	coordResp, err := s.openWeatherClient.GetCoordinates(city)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates: %v", err)
	}
	defer coordResp.Body.Close()

	var coords []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	if err := json.NewDecoder(coordResp.Body).Decode(&coords); err != nil || len(coords) == 0 {
		return nil, fmt.Errorf("failed to decode coordinates response: %v", err)
	}

	lat, lon := fmt.Sprintf("%.4f", coords[0].Lat), fmt.Sprintf("%.4f", coords[0].Lon)

	cachedWeatherAPI, err := s.mongoDB.GetCachedWeather(city, date, "WeatherAPI")
	if err == nil && cachedWeatherAPI != nil {
		log.Println("data has taken from weatherAPidb")
		summary["WeatherAPI"] = cachedWeatherAPI.Temperature
	} else {
		weatherAPIResp, err := s.weatherAPIClient.GetWeather(lat, lon, date)
		if err != nil {
			return nil, fmt.Errorf("failed to get weather from WeatherAPI: %v", err)
		}
		defer weatherAPIResp.Body.Close()
		weatherAPIData, _ := io.ReadAll(weatherAPIResp.Body)
		weatherAPIAvgTemp, err := extractWeatherAPIAvgTemp(weatherAPIData, date)
		if err != nil {
			return nil, err
		}
		summary["WeatherApi"] = weatherAPIAvgTemp
		err = s.mongoDB.CacheWeather(city, date, "WeatherAPI", weatherAPIAvgTemp)
		if err != nil {
			return nil, fmt.Errorf("failed to cache WeatherAPI data: %v", err)
		}
	}

	endTime := parseDate(date)
	cachedTomorrowIO, err := s.mongoDB.GetCachedWeather(city, date, "TomorrowIO")
	if err == nil && cachedTomorrowIO != nil {
		log.Println("data has taken from TomorroIOdb")
		summary["TomorrowIO"] = cachedTomorrowIO.Temperature
	} else {
		tomorrowIOResp, err := s.tomorrowIOClient.GetWeather(lat, lon, date, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get weather from TomorrowIO: %v", err)
		}
		defer tomorrowIOResp.Body.Close()
		tomorrowIOData, _ := io.ReadAll(tomorrowIOResp.Body)
		TomorrowIOAvgTemp, err := extractTomorrowIOAvgTemp(tomorrowIOData)
		if err != nil {
			return nil, err
		}
		summary["TomorrowIo"] = TomorrowIOAvgTemp
		err = s.mongoDB.CacheWeather(city, date, "TomorrowIO", TomorrowIOAvgTemp)
		if err != nil {
			return nil, fmt.Errorf("failed to cache TomorrowIO data: %v", err)
		}
	}
	return summary, nil
}

func parseDate(datestr string) string {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, datestr)
	date = date.AddDate(0, 0, 1)
	newDateStr := date.Format(layout)
	log.Println("New date:", newDateStr)
	return newDateStr
}

func (s *WeatherService) CheckForAvailibleDate(dateStr string, days int) error {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Parsing error:", err)
		return err
	}
	now := time.Now()
	futureDate := now.AddDate(0, 0, days)
	if date.After(futureDate) {
		return fmt.Errorf("the specified date is %v days from the current date", days)
	} else {
		return nil
	}
}

func extractWeatherAPIAvgTemp(responseBody []byte, date string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", fmt.Errorf("failed to decode WeatherAPI response: %v", err)
	}

	if forecast, ok := data["forecast"].(map[string]interface{}); ok {
		if forecastday, ok := forecast["forecastday"].([]interface{}); ok {
			for _, day := range forecastday {
				if dayData, ok := day.(map[string]interface{}); ok {
					if dateStr, ok := dayData["date"].(string); ok && dateStr == date {
						if dayDetail, ok := dayData["day"].(map[string]interface{}); ok {
							if avgTemp, ok := dayDetail["avgtemp_c"].(float64); ok {
								return fmt.Sprintf("%.1f", avgTemp), nil
							}
						}
					}
				}
			}
		}
	}
	return "", nil
}

func extractTomorrowIOAvgTemp(responseBody []byte) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", fmt.Errorf("failed to decode TomorrowIO response: %v", err)
	}

	if timelines, ok := data["data"].(map[string]interface{})["timelines"].([]interface{}); ok && len(timelines) > 0 {
		var totalTemp float64
		var count int
		for _, timeline := range timelines {
			if intervals, ok := timeline.(map[string]interface{})["intervals"].([]interface{}); ok {
				for _, interval := range intervals {
					if values, ok := interval.(map[string]interface{})["values"].(map[string]interface{}); ok {
						if temp, ok := values["temperature"].(float64); ok {
							totalTemp += temp
							count++
						}
					}
				}
			}
		}
		if count > 0 {
			averageTemp := totalTemp / float64(count)
			return fmt.Sprintf("%.1f", averageTemp), nil
		}
	}
	return "", nil
}
