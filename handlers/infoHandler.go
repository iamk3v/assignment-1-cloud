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

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleInfoGetRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}

}

func handleInfoGetRequest(w http.ResponseWriter, r *http.Request) {
	validLimit := `^\d+/?$` // Must be an integer

	use := "Welcome to the info endpoint!\n\nHere is a quick guide to use it:\n" +
		"A valid two letter country code is required as a parameter, with an optional limit query.\nExamples: " +
		"\n	/info/no          - General info, including ALL cities" +
		"\n	/info/no?limit=10 - General info, LIMIT to 10 cities\n" +
		"\nA list of valid country codes can be found here: https://en.wikipedia.org/wiki/ISO_3166-2"

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
			"Expected: /info/{two_letter_country_code}\n"+
			"Example: /info/no\n"+
			"Got: '"+countryCode+"'", http.StatusBadRequest)
		return
	}

	// If we have a query, validate it
	if r.URL.RawQuery != "" {
		if limitQuery == "" || !limitPattern.MatchString(limitQuery) {
			http.Error(w, "Invalid limit query:\n"+
				"Expected format: 'NUMBER'\n"+
				"Example: /info/no?limit=10\n"+
				"Got: '"+r.URL.RawQuery+"'", http.StatusBadRequest)
			return
		}
	}

	var infoResponse []utils.RestCountriesJson

	// Get the country info from country code
	statusCode, err := clients.GetCountryInfo(w, countryCode, &infoResponse)
	if err != nil {
		log.Print("Error fetching country data with status code '" + strconv.Itoa(statusCode) + "': " + err.Error())
		return
	}

	// Define the post data for the cities request
	postData := map[string]string{
		"country": infoResponse[0].Name.Common,
	}

	// Get the cities for the country
	cityResponse := utils.CitiesJson{}
	statusCode, err = clients.GetCities(w, postData, &cityResponse)
	if err != nil {
		log.Print("Error fetching city data with status code '" + strconv.Itoa(statusCode) + "': " + err.Error())
		return
	}

	// If there is a limit query, slice cities
	if limitQuery != "" {
		// Check if it ends with "/" and remove it
		if strings.HasSuffix(limitQuery, "/") {
			limitQuery = strings.TrimSuffix(limitQuery, "/")
		}
		limit, convertErr := strconv.Atoi(limitQuery)
		if convertErr != nil {
			log.Print("Error converting limit query to int: " + convertErr.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
		cityResponse.Cities = cityResponse.Cities[:limit]
	}

	// Define the response data to the user
	finalInfo := utils.InfoResponseJson{
		Name:       infoResponse[0].Name.Common,
		Continents: infoResponse[0].Continents,
		Population: infoResponse[0].Population,
		Languages:  infoResponse[0].Languages,
		Borders:    infoResponse[0].Borders,
		Flag:       infoResponse[0].Flags.Png,
		Capital:    infoResponse[0].Capital,
		Cities:     cityResponse.Cities,
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response to the user
	err = json.NewEncoder(w).Encode(finalInfo)
	if err != nil {
		log.Print("Error occurred when trying to encode and send response: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
