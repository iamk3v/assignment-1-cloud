package handlers

import (
	"Assigment1/constants"
	"Assigment1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	infoLimitRegex := `^\d+$`

	use := "Welcome to the info endpoint!\n\nHere is a quick guide to use it:\n" +
		"A valid two letter country code is required as a parameter, with an optional limit query.\nExamples: " +
		"\n	/info/no          - General info, including ALL cities" +
		"\n	/info/no?limit=10 - General info, LIMIT to 10 cities\n" +
		"\nA list of valid country codes can be found here: https://en.wikipedia.org/wiki/ISO_3166-2"

	infoLimitPattern := regexp.MustCompile(infoLimitRegex)     // Compile year format to a regex
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
		return // Return as there is no country code given
	}

	if !countryCodePattern.MatchString(countryCode) {
		http.Error(w, "Invalid country code format:\n"+
			"Expected: /info/{two_letter_country_code}\n"+
			"Example: /info/no\n"+
			"Got: '"+countryCode+"'", http.StatusBadRequest)
		return
	}

	if r.URL.RawQuery != "" { // If we have a query, validate it
		if limitQuery == "" || !infoLimitPattern.MatchString(limitQuery) {
			http.Error(w, "Invalid limit query:\n"+
				"Expected format: 'NUMBER'\n"+
				"Example: /info/no?limit=10\n"+
				"Got: '"+limitQuery+"'", http.StatusBadRequest)
			return
		}
	}

	var infoResponse []utils.RestCountriesJson

	statusCode, err := utils.GetURL(constants.RestCountriesAPI+"alpha/"+countryCode, &infoResponse)
	if err != nil {
		if statusCode == http.StatusNotFound {
			http.Error(w, "No country found with that country code..", http.StatusNotFound)
		} else {
			log.Print("Error fetching country data: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		}
		return
	}

	postData := map[string]string{
		"country": infoResponse[0].Name.Common,
	}

	cityResponse := utils.CitiesJson{}
	statusCode, err = utils.PostURL(constants.CountriesNowAPI+"countries/cities", postData, &cityResponse)
	if err != nil {
		log.Print("Error fetching city data: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}

	// If there is a limit, slice cities
	if limitQuery != "" {
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			log.Print("Error converting limit query to int: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
		cityResponse.Cities = cityResponse.Cities[:limit]
	}

	finalInfo := utils.InfoJson{
		Name:       infoResponse[0].Name.Common,
		Continents: infoResponse[0].Continents,
		Population: infoResponse[0].Population,
		Languages:  infoResponse[0].Languages,
		Borders:    infoResponse[0].Borders,
		Flag:       infoResponse[0].Flags.Png,
		Capital:    infoResponse[0].Capital,
		Cities:     cityResponse.Cities,
	}

	// Convert back to Json
	jsonData, err := json.Marshal(finalInfo)
	if err != nil {
		log.Print("Failed to Marshal finalInfo: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
	}

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
