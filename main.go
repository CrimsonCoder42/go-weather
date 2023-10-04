package main

import (
	"encoding/json"     // Importing the encoding/json package for JSON processing
	"net/http"          // Importing the net/http package for HTTP server and client implementation
	"os"                // Importing the os package for interacting with the OS
	"strings"           // Importing the strings package for string manipulation
)

// Defining a struct to hold API configuration data
type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"openWeatherMapApiKey"` // JSON field tag to specify the format of the data in JSON
}

// Defining a struct to hold the weather data
type weatherData struct {
	Name string `json:"name"` // Name of the city
	Main struct {             // Nested struct to hold the temperature data
		Kelvin float64 `json:"temp"` // Temperature in Kelvin
	} `json:"main"`                  // JSON field tag
}

// loadApiConfig function reads and unmarshals the API configuration from a file
func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename) // Reading the file content

	if err != nil {
		return apiConfigData{}, err // Returning an error if file reading fails
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c) // Unmarshaling the JSON content into the apiConfigData struct
	if err != nil {
		return apiConfigData{}, err // Returning an error if unmarshaling fails
	}
	return c, nil // Returning the API configuration data
}

// hello function responds with a greeting message
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!")) // Writing a greeting message to the response
}

// query function requests weather data for a specified city from the OpenWeatherMap API
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig") // Loading the API configuration data from the .apiConfig file
	if err != nil {
		return weatherData{}, err // Returning an error if loading fails
	}

	// Making an HTTP GET request to the OpenWeatherMap API with the specified city and API key
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err // Returning an error if the request fails
	}
	defer resp.Body.Close() // Closing the response body when the function exits

	var d weatherData
	// Decoding the JSON response into the weatherData struct
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err // Returning an error if decoding fails
	}

	return d, nil // Returning the weather data
}

// The main function where the execution of the program begins
func main() {
	http.HandleFunc("/", hello) // Setting up the route for the greeting message

	// Setting up the route for the weather data, with a function to handle the requests
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			// Extracting the city from the URL path
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			// Querying the weather data for the specified city
			data, err := query(city)
			if err != nil {
				// Writing an error to the response if the query fails
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Setting the content type of the response to JSON
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			// Encoding and writing the weather data to the response
			json.NewEncoder(w).Encode(data)
		})
	// Starting the HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}
