package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system OS environment variables")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*" // Fallback to allow all if not set
	}

	client := &http.Client{}

	cache := &Cache{
		items: make(map[string]CacheEntry),
	}

	http.HandleFunc("/api/v1/weather", corsMiddleware(weatherHandler(client, cache), allowedOrigin))
	http.HandleFunc("/api/v1/forecast", corsMiddleware(weatherForecastHandler(client, cache), allowedOrigin))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
