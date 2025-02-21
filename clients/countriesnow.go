package clients

import (
	"Assigment-1/config"
	"Assigment-1/utils"
	"net/http"
)

func GetPopulation(w http.ResponseWriter, data interface{}, popRes *utils.PopulationData) (statusCode int, err error) {
	// Get population with post request
	status, getErr := utils.PostURL(config.COUNTRIESNOW_ROOT+"countries/population", data, &popRes)
	// If an error, send internal server error and return status + error
	if getErr != nil {
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return status, getErr
	}
	// Everything went well
	return http.StatusOK, nil
}

func GetCities(w http.ResponseWriter, data interface{}, cityRes *utils.CitiesJson) (statusCode int, err error) {
	// Get cities with post request
	status, getErr := utils.PostURL(config.COUNTRIESNOW_ROOT+"countries/cities", data, &cityRes)
	// If an error, send internal server error and return status + error
	if getErr != nil {
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return status, getErr
	}
	// Everything went well
	return http.StatusOK, nil
}
