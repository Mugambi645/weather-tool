package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv" 
)

const (
	currentWeatherURL = "https://api.openweathermap.org/data/2.5/weather"
	forecastURL       = "https://api.openweathermap.org/data/2.5/forecast"
)

// --- Data Structures (Remain the same) ---
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main describes the main weather parameters (temperature, humidity, pressure)
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

// Wind describes wind speed and direction
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

// Clouds describes cloudiness
type Clouds struct {
	All int `json:"all"`
}

// Sys describes sunrise and sunset times (for current weather)
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

// Coord describes geographical coordinates
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// CurrentWeatherResponse is the top-level struct for current weather API response
type CurrentWeatherResponse struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int64     `json:"dt"` // Time of data calculation, Unix, UTC
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

// City describes the city information in the forecast response
type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Coord      Coord  `json:"coord"`
	Country    string `json:"country"`
	Population int    `json:"population"`
	Timezone   int    `json:"timezone"`
	Sunrise    int64  `json:"sunrise"`
	Sunset     int64  `json:"sunset"`
}

// ForecastListEntry describes a single 3-hour forecast entry
type ForecastListEntry struct {
	Dt         int64     `json:"dt"` // Time of data calculation, Unix, UTC
	Main       Main      `json:"main"`
	Weather    []Weather `json:"json:"weather"`
	Clouds     Clouds    `json:"clouds"`
	Wind       Wind      `json:"wind"`
	Visibility int       `json:"visibility"`
	Pop        float64   `json:"pop"` // Probability of precipitation
	Sys        struct {
		Pod string `json:"pod"` // Part of the day (d = day, n = night)
	} `json:"sys"`
	DtTxt string `json:"dt_txt"` // Date and time in UTC
}

// ForecastResponse is the top-level struct for 5-day / 3-hour forecast API response
type ForecastResponse struct {
	Cod     string              `json:"cod"`
	Message float64             `json:"message"`
	Cnt     int                 `json:"cnt"`
	List    []ForecastListEntry `json:"list"`
	City    City                `json:"city"`
}

// --- API Client Functions (Remain the same) ---
func fetchWeatherData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return nil
}

// GetCurrentWeather fetches current weather data for a given city.
func GetCurrentWeather(city string, apiKey string) (*CurrentWeatherResponse, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", currentWeatherURL, city, apiKey)
	var weatherData CurrentWeatherResponse
	err := fetchWeatherData(url, &weatherData)
	if err != nil {
		return nil, err
	}
	return &weatherData, nil
}

// GetForecast fetches 5-day / 3-hour forecast data for a given city.
func GetForecast(city string, apiKey string) (*ForecastResponse, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", forecastURL, city, apiKey)
	var forecastData ForecastResponse
	err := fetchWeatherData(url, &forecastData)
	if err != nil {
		return nil, err
	}
	return &forecastData, nil
}

// --- Display Functions (Remain the same) ---
func displayCurrentWeather(data *CurrentWeatherResponse) {
	fmt.Printf("Current Weather for %s, %s:\n", data.Name, data.Sys.Country)
	fmt.Printf("  Temperature: %.1f°C (Feels like: %.1f°C)\n", data.Main.Temp, data.Main.FeelsLike)
	fmt.Printf("  Conditions: %s (%s)\n", data.Weather[0].Main, data.Weather[0].Description)
	fmt.Printf("  Humidity: %d%%\n", data.Main.Humidity)
	fmt.Printf("  Wind: %.1f m/s\n", data.Wind.Speed)
	fmt.Printf("  Pressure: %d hPa\n", data.Main.Pressure)
	fmt.Printf("  Cloudiness: %d%%\n", data.Clouds.All)
	fmt.Printf("  Sunrise: %s\n", time.Unix(data.Sys.Sunrise, 0).Local().Format("15:04"))
	fmt.Printf("  Sunset: %s\n", time.Unix(data.Sys.Sunset, 0).Local().Format("15:04"))
	fmt.Println("------------------------------------")
}
// displayForecast prints the 5-day / 3-hour forecast details.
func displayForecast(data *ForecastResponse) {
	fmt.Printf("5-Day / 3-Hour Forecast for %s, %s:\n", data.City.Name, data.City.Country)
	fmt.Println("------------------------------------")

	// Group forecast entries by day
	dailyForecasts := make(map[string][]ForecastListEntry)
	for _, entry := range data.List {
		date := time.Unix(entry.Dt, 0).Local().Format("2006-01-02 (Mon)")
		dailyForecasts[date] = append(dailyForecasts[date], entry)
	}

	// Sort dates for consistent output
	var dates []string
	for date := range dailyForecasts {
		dates = append(dates, date)
	}
	// Simple bubble sort for demonstration, for larger sets use sort.Strings
	for i := 0; i < len(dates)-1; i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] > dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	for _, date := range dates {
		fmt.Printf("\nDate: %s\n", date)
		for _, entry := range dailyForecasts[date] {
			forecastTime := time.Unix(entry.Dt, 0).Local().Format("15:04")

			// --- FIX STARTS HERE ---
			var mainWeather, descWeather string
			if len(entry.Weather) > 0 {
				mainWeather = entry.Weather[0].Main
				descWeather = entry.Weather[0].Description
			} else {
				// Provide default values if weather array is empty
				mainWeather = "N/A"
				descWeather = "No specific conditions"
			}
			// --- FIX ENDS HERE ---

			fmt.Printf("  %s: Temp: %.1f°C, Feels: %.1f°C, Cond: %s (%s), Wind: %.1f m/s, Pop: %.0f%%\n",
				forecastTime,
				entry.Main.Temp,
				entry.Main.FeelsLike,
				mainWeather,       // Use the checked variable
				descWeather,       // Use the checked variable
				entry.Wind.Speed,
				entry.Pop*100,
			)
		}
	}
	fmt.Println("------------------------------------")
}

func main() {
	// Load environment variables from .env file
	// godotenv.Load() without arguments looks for .env in the current directory
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Could not load .env file. Falling back to system environment variables.")
		// It's okay if .env doesn't exist, as system env vars might be used in production
	}

	// Define command-line flags
	cityPtr := flag.String("city", "", "City name (e.g., 'London', 'Nairobi')")
	forecastPtr := flag.Bool("forecast", false, "Get 5-day / 3-hour forecast instead of current weather")

	flag.Parse()

	// Read API key from environment variable (will now check loaded .env first, then system env)
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	// Validate API Key
	if apiKey == "" {
		fmt.Println("Error: OpenWeatherMap API key not found.")
		fmt.Println("Please set the OPENWEATHER_API_KEY environment variable in a .env file or directly in your shell.")
		fmt.Println("Example .env entry: OPENWEATHER_API_KEY=\"YOUR_ACTUAL_API_KEY\"")
		os.Exit(1)
	}

	// Validate city input
	if *cityPtr == "" {
		fmt.Println("Error: Please provide a city name using the --city flag.")
		fmt.Println("Usage: go run main.go --city \"YourCity\" [--forecast]")
		os.Exit(1)
	}

	if *forecastPtr {
		forecastData, err := GetForecast(*cityPtr, apiKey)
		if err != nil {
			fmt.Printf("Error fetching forecast for %s: %v\n", *cityPtr, err)
			os.Exit(1)
		}
		displayForecast(forecastData)
	} else {
		weatherData, err := GetCurrentWeather(*cityPtr, apiKey)
		if err != nil {
			fmt.Printf("Error fetching current weather for %s: %v\n", *cityPtr, err)
			os.Exit(1)
		}
		displayCurrentWeather(weatherData)
	}
}