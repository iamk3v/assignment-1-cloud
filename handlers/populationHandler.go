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

const ValidCountryCode = `^[a-zA-Z]{2}$`

func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	validYearPattern := `^\d{4}-\d{4}$` // YYYY-YYYY

	use := "Welcome to the population endpoint!\n\nHere is a quick guide to use it:\n" +
		"A valid two letter country code is required as a parameter, with an optional limit query.\nExamples:" +
		"\n	/population/no                 - All registered years" +
		"\n	/population/no?limit=2020-2025 - Limit to last 5 years" +
		"\n\nA list of valid country codes can be found here: https://en.wikipedia.org/wiki/ISO_3166-2"

	yearPattern := regexp.MustCompile(validYearPattern)        // Compile year format to a regex
	countryCodePattern := regexp.MustCompile(ValidCountryCode) // Compile country code
	countryCode := r.PathValue("two_letter_country_code")
	limitQuery := r.URL.Query().Get("limit")

	if countryCode == "" {
		_, err := fmt.Fprintf(w, "%v", use)
		if err != nil {
			log.Print("An error occurred: " + err.Error())
			http.Error(w, "Error when returning output", http.StatusInternalServerError)
			return
		}
		return
	}

	if !countryCodePattern.MatchString(countryCode) {
		http.Error(w, "Invalid country code format:\n"+
			"Expected: /population/{two_letter_country_code}\n"+
			"Example: /population/no\n"+
			"Got: '"+countryCode+"'", http.StatusBadRequest)
		return
	}

	if r.URL.RawQuery != "" { // If we have a query, validate it
		if limitQuery == "" || !yearPattern.MatchString(limitQuery) {
			http.Error(w, "Invalid limit query:\n"+
				"Expected format: 'YYYY-YYYY'\n"+
				"Example: /population/?limit=2020-2025\n"+
				"Got: '"+limitQuery+"'", http.StatusBadRequest)
			return
		}
	}

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

	populationResponse := utils.PopulationData{}
	postData := map[string]string{
		"country": countryName.Name.Common,
	}

	statusCode, err = utils.PostURL(constants.CountriesNowAPI+"countries/population", postData, &populationResponse)
	if err != nil {
		log.Print("Error fetching population data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
	}

	count, sumYears := 0, 0
	for _, v := range populationResponse.Data.PopulationCounts {
		sumYears += v.Value
		count++
	}

	populationInfo := utils.PopulationJson{
		Mean:   sumYears / count,
		Values: populationResponse.Data.PopulationCounts,
	}

	// Convert back to Json
	jsonData, err := json.Marshal(populationInfo)

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send data
	_, err = fmt.Fprintf(w, "%v", string(jsonData))
	if err != nil {
		log.Print("Error when returning output: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
