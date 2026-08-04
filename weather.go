package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func getCoordinates(client *http.Client, cityName string) (GeocodingResult, error) {
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", url.QueryEscape(cityName))

	resp, err := client.Get(geoURL)
	if err != nil {
		return GeocodingResult{}, err
	}
	defer resp.Body.Close()

	var geoData GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return GeocodingResult{}, err
	}

	if len(geoData.Results) == 0 {
		return GeocodingResult{}, fmt.Errorf("city %s not found", cityName)
	}

	return geoData.Results[0], nil
}

func getWeather(client *http.Client, lat, lon float64) (WeatherResponse, error) {
	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m", lat, lon)

	resp, err := client.Get(weatherURL)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return WeatherResponse{}, err
	}

	return weather, nil
}

func getWeatherForDays(client *http.Client, lat, lon float64, days string) (OpenMeteoForecastResponse, error) {
	weatherUrl := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&daily=temperature_2m_max,temperature_2m_min&forecast_days=%s", lat, lon, days)
	resp, err := client.Get(weatherUrl)
	if err != nil {
		return OpenMeteoForecastResponse{}, err
	}
	defer resp.Body.Close()

	var weather OpenMeteoForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return OpenMeteoForecastResponse{}, err
	}

	return weather, nil
}
