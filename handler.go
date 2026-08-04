package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func weatherHandler(client *http.Client, cache *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cityName := r.URL.Query().Get("city")
		if cityName == "" {
			http.Error(w, "city name is required", http.StatusBadRequest)
			return
		}

		cache.mu.RLock()
		entry, found := cache.items[cityName]
		cache.mu.RUnlock()
		if found && time.Since(entry.CreatedAt) < 5*time.Minute {
			log.Println("Cache hit! Returning cached data")
			w.Header().Set("Content-Type", "application/json")
			w.Write(entry.Data)
			return
		}

		log.Println("Cache miss! Fetching from Open-Meteo")

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

		response := APIResponse{
			City:        location.Name,
			Country:     location.Country,
			Temperature: fmt.Sprintf("%.1f%s", weather.Current.Temperature2m, weather.CurrentUnits.Temperature2m),
		}

		jsonBytes, _ := json.Marshal(response)

		cache.mu.Lock()
		cache.items[cityName] = CacheEntry{
			Data:      jsonBytes,
			CreatedAt: time.Now(),
		}
		cache.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}

}

func weatherForecastHandler(client *http.Client, cache *Cache) http.HandlerFunc {
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

		cacheKey := fmt.Sprintf("%s:%s", url.QueryEscape(cityName), days)

		cache.mu.RLock()
		entry, found := cache.items[cacheKey]
		cache.mu.RUnlock()

		if found && time.Since(entry.CreatedAt) < 5*time.Minute {
			log.Println("Forecast Cache HIT for key:", cacheKey)
			w.Header().Set("Content-Type", "application/json")
			w.Write(entry.Data)
			return
		}

		log.Println("Forecast Cache MISS for key:", cacheKey)

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

		response := ForecastAPIResponse{
			City:      location.Name,
			Country:   location.Country,
			Forecasts: forecasts,
		}

		jsonBytes, _ := json.Marshal(response)

		cache.mu.Lock()
		cache.items[cacheKey] = CacheEntry{
			Data:      jsonBytes,
			CreatedAt: time.Now(),
		}
		cache.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
