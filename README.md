# Weather App Backend (Go)

## What it is
A lightweight, high-performance REST API backend for the Weather App, written in Go.

## What it does
* **Open-Meteo Integration**: Fetches real-time geocoding and weather data from the public [Open-Meteo API](https://open-meteo.com/).
* **City Coordinates Lookup**: Resolves city names (e.g., "Maputo", "New York") into geographical coordinates (latitude and longitude).
* **Current Weather & Forecasts**: Provides current temperature and multi-day forecasts (minimum and maximum daily temperatures).
* **Caching Layer**: Implements a thread-safe, in-memory caching system. Weather and forecast queries are cached for **5 minutes** to minimize external API calls and speed up responses.
* **CORS Support**: Provides middleware to support Cross-Origin Resource Sharing (CORS) configured via an environment variable.

## How to execute it
### Prerequisites
* Go 1.21+ installed on your system.

### Steps
1. **Configure Environment Variables** (Optional):
   Create a `.env` file in the root directory:
   ```env
   ALLOWED_ORIGIN=http://localhost:5173
   ```
   *If `.env` is omitted, `ALLOWED_ORIGIN` defaults to `*`.*

2. **Run the Server**:
   ```bash
   go run .
   ```
   The backend server will start and listen on port `8080` (`http://localhost:8080`).

### API Endpoints
* **Current Weather**:
  `GET /api/v1/weather?city=<city_name>`
* **Weather Forecast**:
  `GET /api/v1/forecast?city=<city_name>&days=<number_of_days>`
