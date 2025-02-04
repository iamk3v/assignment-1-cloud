package handlers

import (
	"Assigment-1/constants"
	"Assigment-1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Define test data to use
	countryCode := "no"
	postData := map[string]string{
		"country": "norway",
	}

	// Test and get the status code from RestCountriesAPI
	restStatusCode, err := utils.TestGetApi(constants.RestCountriesAPI + "alpha/" + countryCode + "?fields=name")
	if err != nil {
		log.Print("Error fetching country data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	// Test and get the status code from CountriesNowAPI
	countriesNowStatusCode, err := utils.TestPostApi(constants.CountriesNowAPI+"countries/cities", postData)
	if err != nil {
		log.Print("Error fetching city data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	// Define the response data to the user
	statusData := map[string]string{
		"countriesnowapi":  strconv.Itoa(restStatusCode),
		"restcountriesapi": strconv.Itoa(countriesNowStatusCode),
		"version":          constants.Version,
		"uptime":           "Not implemented yet",
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert response data to JSON
	jsonData, err := json.Marshal(statusData)
	if err != nil {
		log.Print("Failed to Marshal statusInfo: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
	}

	// Send the response to the user
	_, err = fmt.Fprintf(w, "%v", string(jsonData))
	if err != nil {
		log.Print("An error occurred: " + err.Error())
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
		return
	}
}
