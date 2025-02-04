package handlers

import (
	"Assigment1/constants"
	"Assigment1/utils"
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

	populationResponse := utils.PopulationJson{}
	postData := map[string]string{
		"country": "norway",
	}

	statusCode, err := utils.PostURL(constants.CountriesNowAPI+"countries/population", postData, &populationResponse)
	if err != nil {
		if statusCode == http.StatusNotFound {
			http.Error(w, "No country found with that country code..", http.StatusNotFound)
		} else {
			fmt.Println("Error fetching population data:", err)
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		}
		return
	}

	fmt.Println("Population Response:", populationResponse)

	if err != nil {
		log.Fatal(err)
	}

	/*
		_, err = fmt.Fprintf(w, "%v", output)
		if err != nil {
			log.Print("An error occurred: " + err.Error())
			http.Error(w, "Error when returning output", http.StatusInternalServerError)
			return
		}
	*/

}
