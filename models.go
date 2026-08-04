package main

import (
	"sync"
	"time"
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

type CacheEntry struct {
	Data      []byte
	CreatedAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheEntry
}
