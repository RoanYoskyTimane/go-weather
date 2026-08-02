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

// Custom response for our API
type APIResponse struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	Temperature string `json:"temperature"`
}

func main() {
	client := &http.Client{}

	http.HandleFunc("/api/v1/weather", weatherHandler(client))
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
