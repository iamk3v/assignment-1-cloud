package handlers

import (
	"Assigment1/constants"
	"Assigment1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

const ValidCountryCode = `^[a-zA-Z]{2}$` // aa to ZZ

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	validLimit := `^\d{4}-\d{4}$` // YYYY-YYYY

	use := "Welcome to the population endpoint!\n\nHere is a quick guide to use it:\n" +
		"A valid two letter country code is required as a parameter, with an optional limit query.\nExamples:" +
		"\n	/population/no                 - All registered years" +
		"\n	/population/no?limit=2020-2025 - Limit to last 5 years" +
		"\n\nA list of valid country codes can be found here: https://en.wikipedia.org/wiki/ISO_3166-2"

	// Regexes to validate country code and limit format
	limitPattern := regexp.MustCompile(validLimit)
	countryCodePattern := regexp.MustCompile(ValidCountryCode)

	// Extract country code and potential limit query from request
	countryCode := r.PathValue("two_letter_country_code")
	limitQuery := r.URL.Query().Get("limit")

	// If only the root path without country code, send use message
	if countryCode == "" {
		_, err := fmt.Fprintf(w, "%v", use)
		if err != nil {
			log.Print("An error occurred: " + err.Error())
			http.Error(w, "Error when returning output", http.StatusInternalServerError)
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
				"Got: '"+limitQuery+"'", http.StatusBadRequest)
			return
		}
	}

	// Get the country name from the country code
	countryName := utils.CountryName{}
	statusCode, err := utils.GetURL(constants.RestCountriesAPI+"alpha/"+countryCode+"?fields=name", &countryName)
	if err != nil {
		if statusCode == http.StatusNotFound {
			log.Print("Invalid country code: " + err.Error())
			http.Error(w, "No country found with that country code..", http.StatusNotFound)
		} else {
			log.Print("Error fetching country name: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		}
		return
	}

	// Define the response struct and post data for population request
	populationResponse := utils.PopulationData{}
	postData := map[string]string{
		"country": countryName.Name.Common,
	}

	// Send the population post request
	statusCode, err = utils.PostURL(constants.CountriesNowAPI+"countries/population", postData, &populationResponse)
	if err != nil {
		log.Print("Error fetching population data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
	}

	// Loop and calculate the sum of population and number of years gotten
	count, sumYears := 0, 0
	for _, v := range populationResponse.Data.PopulationCounts {
		sumYears += v.Value
		count++
	}

	// Define the response data to the user
	populationInfo := utils.PopulationJson{
		Mean:   sumYears / count,
		Values: populationResponse.Data.PopulationCounts,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(populationInfo)
	if err != nil {
		log.Print("Failed to Marshal populationInfo: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send the response to the user
	_, err = fmt.Fprintf(w, "%v", string(jsonData))
	if err != nil {
		log.Print("Error when returning output: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
