package handlers

import (
	"Assigment-1/clients"
	"Assigment-1/config"
	"Assigment-1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlePopulationGetRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}

}

func handlePopulationGetRequest(w http.ResponseWriter, r *http.Request) {
	validLimit := `^\d{4}-\d{4}/?$` // YYYY-YYYY

	use := "Welcome to the population endpoint!\n\nHere is a quick guide to use it:\n" +
		"A valid two letter country code is required as a parameter, with an optional limit query.\nExamples:" +
		"\n	/population/no                 - All registered years" +
		"\n	/population/no?limit=2020-2025 - Limit to last 5 years" +
		"\n\nA list of valid country codes can be found here: https://en.wikipedia.org/wiki/ISO_3166-2"

	// Regexes to validate country code and limit format
	limitPattern := regexp.MustCompile(validLimit)
	countryCodePattern := regexp.MustCompile(config.VALID_COUNTRY_CODE)

	// Extract country code and potential limit query from request
	countryCode := r.PathValue("two_letter_country_code")
	limitQuery := r.URL.Query().Get("limit")

	// If only the root path without country code, send use message
	if countryCode == "" {
		_, err := fmt.Fprintf(w, "%v", use)
		if err != nil {
			log.Print("Error occurred when trying to send response: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle invalid country code
	if !countryCodePattern.MatchString(countryCode) {
		http.Error(w, "Invalid country code format:\n"+
			"Expected: /population/{two_letter_country_code}\n"+
			"Example: /population/no\n"+
			"Got: '"+countryCode+"'", http.StatusBadRequest)
		return
	}

	// If we have a query, validate it
	if r.URL.RawQuery != "" {
		if limitQuery == "" || !limitPattern.MatchString(limitQuery) {
			http.Error(w, "Invalid limit query:\n"+
				"Expected format: 'YYYY-YYYY'\n"+
				"Example: /population/?limit=2020-2025\n"+
				"Got: '"+r.URL.RawQuery+"'", http.StatusBadRequest)
			return
		}
	}

	// Get the country name from the country code
	countryName := utils.CountryName{}
	statusCode, err := clients.GetCountryName(w, countryCode, &countryName)
	if err != nil {
		log.Print("Error fetching country name with status code '" + strconv.Itoa(statusCode) + "': " + err.Error())
		return
	}

	// Define the response struct and post data for population request
	postData := map[string]string{
		"country": countryName.Name.Common,
	}

	// Get the population
	populationResponse := utils.PopulationData{}
	statusCode, err = clients.GetPopulation(w, postData, &populationResponse)
	if err != nil {
		log.Print("Error fetching population with status code '" + strconv.Itoa(statusCode) + "': " + err.Error())
		return
	}

	// If there is a limit query
	if limitQuery != "" {
		if strings.HasSuffix(limitQuery, "/") {
			limitQuery = strings.TrimSuffix(limitQuery, "/")
		}

		// Split the query into two years
		years := strings.Split(limitQuery, "-")

		// Extract the individual years
		startYear, convertErr := strconv.Atoi(years[0])
		if convertErr != nil {
			log.Print("Error converting startYear to int" + convertErr.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
		endYear, convertErr := strconv.Atoi(years[1])
		if convertErr != nil {
			log.Print("Error converting endYear to int" + convertErr.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}

		// Filter for limit
		utils.FilterYears(&populationResponse, startYear, endYear)
	}

	// Get count and sum of years
	count, sumYears := utils.CalculateYears(&populationResponse)

	// Define the response data to the user
	populationInfo := utils.PopulationResponseJson{
		Mean:   sumYears / count,
		Values: populationResponse.Data.PopulationCounts,
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response to the user
	err = json.NewEncoder(w).Encode(populationInfo)
	if err != nil {
		log.Print("Error occurred when trying to encode and send response: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
