# Weather API Service

This is a Go-based Weather API service that provides weather data for a given city using the [Visual Crossing Weather API](https://www.visualcrossing.com/). The service caches the results in Redis for faster responses and efficient resource usage.

---

## Features

- **Get Current Weather**: Fetch weather information for a specified city and unit (e.g., metric or imperial).
- **Redis Caching**: Weather data is cached for 12 hours to improve performance and reduce API calls.
- **Gin Framework**: Utilizes the Gin framework for HTTP routing and JSON responses.
- **Error Handling**: Handles common errors such as missing API keys, invalid parameters, and Redis issues.

---

## API Endpoints

### `/weather`
Fetch weather information for a specified city.

#### Query Parameters:
- `city` (required): The name of the city.
- `unit` (required): The unit for temperature (`metric` or `imperial`).

#### Example Request:
GET http://localhost:8080/weather?city=London&unit=metric


#### Example Response:
```json
{
  "resolvedAddress": "London, England, United Kingdom",
  "timezone": "Europe/London",
  "days": [
    {
      "datetime": "2025-01-09",
      "temp": 7.2,
      "conditions": "Partially Cloudy",
      "description": "A mix of clouds and sun."
    }
  ]
}
```

#### How to run:

Prerequisites

Go installed on your system.
Redis running locally.
Visual Crossing API Key for accessing weather data.
.env file with the following content:

1.Clone the repository:
```bash
git clone https://github.com/r3iwan/weather-api-service.git
cd weather-api-service
```

2.Install dependencies:
```bash
go mod tidy
```

3.Run the application:
```bash
go run main.go
```

4.Access the API at http://localhost:8080/weather.
