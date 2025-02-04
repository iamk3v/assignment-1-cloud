package handlers

import (
	"Assigment1/constants"
	"Assigment1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	// Testing status with Norway
	countryCode := "za"
	postData := map[string]string{
		"country": "norway",
	}

	restStatusCode, err := utils.TestGetApi(constants.RestCountriesAPI + "alpha/" + countryCode)
	if err != nil {
		log.Print("Error fetching country data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	countriesNowStatusCode, err := utils.TestPostApi(constants.CountriesNowAPI+"countries/cities", postData)
	if err != nil {
		log.Print("Error fetching city data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	statusData := map[string]string{
		"countriesnowapi":  strconv.Itoa(restStatusCode),
		"restcountriesapi": strconv.Itoa(countriesNowStatusCode),
		"version":          constants.Version,
		"uptime":           "Not implemented yet",
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert back to Json
	jsonData, err := json.Marshal(statusData)

	// Send data
	_, err = fmt.Fprintf(w, "%v", string(jsonData))
	if err != nil {
		log.Print("An error occurred: " + err.Error())
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
		return
	}
}
