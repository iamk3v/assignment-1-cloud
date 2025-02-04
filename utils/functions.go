package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetURL(url string, targetStruct interface{}) (statusCode int, err error) {
	// Send get request
	res, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer res.Body.Close() // Ensure that the body is close

	err = json.NewDecoder(res.Body).Decode(targetStruct)
	if err != nil {
		return res.StatusCode, fmt.Errorf("failed to decode response body: %v", err)
	}

	return http.StatusOK, nil
}

func PostURL(url string, data interface{}, targetStruct interface{}) (statusCode int, err error) {
	// Convert into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to convert data to JSON: %v", err)
	}

	// Send post request
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer res.Body.Close() // Ensure that the body is close

	err = json.NewDecoder(res.Body).Decode(targetStruct)
	if err != nil {
		return res.StatusCode, fmt.Errorf("failed to decode response body: %v", err)
	}

	return http.StatusOK, nil
}

func TestGetApi(url string) (statusCode int, err error) {
	// Send get request
	res, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer res.Body.Close() // Ensure that the body is close

	return res.StatusCode, nil
}

func TestPostApi(url string, data interface{}) (statusCode int, err error) {
	// Convert into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to convert data to JSON: %v", err)
	}

	// Send post request
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer res.Body.Close() // Ensure that the body is close

	return res.StatusCode, nil
}
