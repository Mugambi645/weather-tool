# Go Weather CLI Tool

A simple yet powerful command-line interface (CLI) tool built with Go to fetch current weather and 5-day / 3-hour forecast data for any city worldwide using the OpenWeatherMap API.

## Features

- **Current Weather**: Get up-to-date temperature, conditions, humidity, wind, and more.
- **5-Day / 3-Hour Forecast**: Plan ahead with detailed hourly forecasts for the next five days.
- **City-Based Search**: Easily query weather for any specified city.
- **Secure API Key Handling**: Utilizes environment variables (with `.env` support for development) to keep your API key out of your codebase.

## Getting Started

Follow these steps to get your Go Weather Tool up and running.

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** (1.16 or newer recommended): Download from [golang.org](https://golang.org).
- **An OpenWeatherMap API Key**:
  - Go to [OpenWeatherMap](https://openweathermap.org).
  - Sign up for a free account.
  - Navigate to your API keys (usually found under your profile settings).
  - Copy your unique API key.

### Installation & Setup

#### Clone the repository (or create the project structure):

If you're starting from scratch, create a directory and initialize a Go module:

```bash
mkdir go-weather-tool
cd go-weather-tool
go mod init go-weather-tool
```

Then create a `main.go` file inside.

#### Add the godotenv dependency:

This package helps load environment variables from a `.env` file during development.

```bash
go get github.com/joho/godotenv
```

#### Create your `.env` file:

In the root of your `go-weather-tool` directory (where `main.go` is), create a file named `.env`.  
Add your OpenWeatherMap API key to it:

```env
# .env
OPENWEATHER_API_KEY="YOUR_ACTUAL_OPENWEATHERMAP_API_KEY"
```

> Replace `"YOUR_ACTUAL_OPENWEATHERMAP_API_KEY"` with your actual key!

#### Populate `main.go`:

Paste the entire Go code from our previous discussion into your `main.go` file.  
Ensure the `godotenv` import and API key fetching logic are correctly implemented.

#### Secure your API Key (Important!):

Add `.env` to your `.gitignore` file to prevent accidentally committing your API key:

```gitignore
# .gitignore
.env
```

## Usage

Run the tool from your terminal.

### Fetch Current Weather

To get the current weather for a city:

```bash
go run main.go --city "Meru"
```

### Fetch 5-Day Forecast

To get the 5-day / 3-hour forecast for a city:

```bash
go run main.go --city "Nairobi" --forecast
```

### Cities with Spaces

For cities with spaces in their names, enclose the city name in quotes:

```bash
go run main.go --city "Mombasa" --forecast
```

## Environment Variables

The tool expects your OpenWeatherMap API key to be available as an environment variable named `OPENWEATHER_API_KEY`.

- **During Development**: The tool will automatically load this from your `.env` file.
- **In Production/Deployment**: Set the environment variable directly in your deployment environment.

**Examples:**

- **Linux/macOS (Bash/Zsh):** `export OPENWEATHER_API_KEY="YOUR_KEY"`
- **Windows (CMD):** `set OPENWEATHER_API_KEY=YOUR_KEY`
- **Windows (PowerShell):** `$env:OPENWEATHER_API_KEY="YOUR_KEY"`

## Error Handling

The tool includes basic error handling for network issues, invalid API keys, and unreadable responses.  
If an error occurs, it will print a descriptive message to the console.

## Contributing

Feel free to fork this repository, open issues, or submit pull requests for any improvements or bug fixes.

## License

This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).
