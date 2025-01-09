package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type jsonResponse struct {
	ResolvedAddress string `json:"resolvedAddress"`
	Timezone        string `json:"timezone"`
	Days            []Days `json:"days"`
}

type Days struct {
	Datetime    string  `json:"datetime"`
	Temperature float64 `json:"temp"`
	Conditions  string  `json:"conditions"`
	Description string  `json:"description"`
}

var rdb *redis.Client

func GetWeatherHandler() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("API key not set")
	}

	router := gin.Default()
	router.GET("/weather", GetWeatherAPI)
	router.Run(":8080")
}

func GetWeatherAPI(c *gin.Context) {
	startTime := time.Now()
	
	city := c.Query("city")
	unit := c.Query("unit")

	if city == "" || unit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City and Unit parameters are required!"})
	}

	var weather jsonResponse
	cachedWeather, err := rdb.Get(city).Result()
	if err == redis.Nil {
		apiKey := os.Getenv("WEATHER_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API key not set"})
			return
		}

		link := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?unitGroup=%s&key=%s&contentType=json", city, unit, apiKey)
		resp, err := http.Get(link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making HTTP request"})
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&weather)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding JSON response"})
			return
		}

		weatherData, err := json.Marshal(weather)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling weather data"})
			return
		}

		err = rdb.Set(city, weatherData, 12*time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving to Redis"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error accessing Redis"})
		return
	} else {
		err = json.Unmarshal([]byte(cachedWeather), &weather)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling cached data"})
			return
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("Request for city: %s took %v\n", city, duration)

	c.JSON(http.StatusOK, weather)
}
