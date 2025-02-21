package clients

import (
	"Assigment-1/config"
	"Assigment-1/utils"
	"fmt"
	"net/http"
)

func GetCountryName(w http.ResponseWriter, countryCode string, countryName *utils.CountryName) (statusCode int, err error) {
	// Get the country name
	status, getErr := utils.GetURL(config.RESTCOUNTRIES_ROOT+"alpha/"+countryCode+"?fields=name", &countryName)
	// If no country by that code
	if status == http.StatusNotFound {
		http.Error(w, "No country found with that country code..", http.StatusNotFound)
		return status, fmt.Errorf("invalid country code")
	}
	// If an error, send internal server error and return status + error
	if getErr != nil {
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return status, getErr
	}
	// Everything went well
	return http.StatusOK, nil
}

func GetCountryInfo(w http.ResponseWriter, countryCode string, infoRes *[]utils.RestCountriesJson) (statusCode int, err error) {
	// Get the country info
	status, getErr := utils.GetURL(config.RESTCOUNTRIES_ROOT+"alpha/"+countryCode, &infoRes)
	// If no country by that code
	if status == http.StatusNotFound {
		http.Error(w, "No country found with that country code..", http.StatusNotFound)
		return status, fmt.Errorf("invalid country code")
	}
	// If an error, send internal server error and return status + error
	if getErr != nil {
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return status, getErr
	}
	// Everything went well
	return http.StatusOK, nil
}
