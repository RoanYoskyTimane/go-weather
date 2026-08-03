package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type GeocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type WeatherResponse struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
	CurrentUnits struct {
		Temperature2m string `json:"temperature_2m"`
	} `json:"current_units"`
}

type OpenMeteoForecastResponse struct {
	DailyUnits struct {
		Temperature2mMax string `json:"temperature_2m_max"`
		Temperature2mMin string `json:"temperature_2m_min"`
	} `json:"daily_units"`
	Daily struct {
		Time             []string  `json:"time"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
		Temperature2mMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type DailyForecast struct {
	Date    string `json:"date"`
	MaxTemp string `json:"max_temp"`
	MinTemp string `json:"min_temp"`
}

type ForecastAPIResponse struct {
	City      string          `json:"city"`
	Country   string          `json:"country"`
	Forecasts []DailyForecast `json:"forecasts"`
}

// Custom response for our API
type APIResponse struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	Temperature string `json:"temperature"`
}

func main() {
	client := &http.Client{}

	http.HandleFunc("/api/v1/weather", weatherHandler(client))
	http.HandleFunc("/api/v1/forecast", weatherForDays(client))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func weatherHandler(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cityName := r.URL.Query().Get("city")
		if cityName == "" {
			http.Error(w, "city name is required", http.StatusBadRequest)
			return
		}

		location, err := getCoordinates(client, cityName)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{Error: "%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		weather, err := getWeather(client, location.Latitude, location.Longitude)
		if err != nil {
			http.Error(w, `{"error":"failed to get weather"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{
			City:        location.Name,
			Country:     location.Country,
			Temperature: fmt.Sprintf("%.1f%s", weather.Current.Temperature2m, weather.CurrentUnits.Temperature2m),
		})
	}

}

func weatherForDays(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cityName := r.URL.Query().Get("city")
		if cityName == "" {
			http.Error(w, "city name is required", http.StatusBadRequest)
			return
		}

		days := r.URL.Query().Get("days")
		if days == "" {
			http.Error(w, "days is required", http.StatusBadRequest)
			return
		}

		location, err := getCoordinates(client, cityName)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{Error: "%s"}`, err.Error()), http.StatusNotFound)
			return
		}

		rawForecast, err := getWeatherForDays(client, location.Latitude, location.Longitude, days)
		if err != nil {
			http.Error(w, `{"error":"failed to get weather"}`, http.StatusInternalServerError)
			return
		}

		forecasts := make([]DailyForecast, 0, len(rawForecast.Daily.Time))
		for i := 0; i < len(rawForecast.Daily.Time); i++ {
			forecasts = append(forecasts, DailyForecast{
				Date:    rawForecast.Daily.Time[i],
				MaxTemp: fmt.Sprintf("%.1f%s", rawForecast.Daily.Temperature2mMax[i], rawForecast.DailyUnits.Temperature2mMax),
				MinTemp: fmt.Sprintf("%.1f%s", rawForecast.Daily.Temperature2mMin[i], rawForecast.DailyUnits.Temperature2mMin),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ForecastAPIResponse{
			City:      location.Name,
			Country:   location.Country,
			Forecasts: forecasts,
		})
	}
}

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
