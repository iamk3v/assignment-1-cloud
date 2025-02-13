package handlers

import (
	"Assigment-1/config"
	"Assigment-1/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}

}

func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	// Define test data to use
	countryCode := "no"
	postData := map[string]string{
		"country": "norway",
	}

	// Test and get the status code from RestCountriesAPI
	restStatusCode, err := utils.TestGetApi(config.RESTCOUNTRIES_ROOT + "alpha/" + countryCode + "?fields=name")
	if err != nil {
		log.Print("Error fetching country name with status code '" + strconv.Itoa(restStatusCode) + "': " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	// Test and get the status code from CountriesNowAPI
	countriesNowStatusCode, err := utils.TestPostApi(config.COUNTRIESNOW_ROOT+"countries/cities", postData)
	if err != nil {
		log.Print("Error fetching population data with status code '" +
			strconv.Itoa(countriesNowStatusCode) + "': " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	// Define the response data to the user
	statusData := map[string]string{
		"countriesnowapi":  strconv.Itoa(restStatusCode) + " - " + config.HTTP_CAT + strconv.Itoa(restStatusCode),
		"restcountriesapi": strconv.Itoa(countriesNowStatusCode) + " - " + config.HTTP_CAT + strconv.Itoa(countriesNowStatusCode),
		"version":          config.VERSION,
		"uptime":           utils.GetUptime(),
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response to the user
	err = json.NewEncoder(w).Encode(statusData)
	if err != nil {
		log.Print("Error occurred when trying to encode and send response: " + err.Error())
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
		return
	}
}
