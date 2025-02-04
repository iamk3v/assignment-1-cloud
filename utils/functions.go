package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetURL(url string, targetStruct interface{}) (statusCode int, err error) {
	// Send GET request
	res, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send GET request: %v", err)
	}

	// Ensure that the body closes at the end of function call
	defer res.Body.Close()

	// Decode the response JSON into target struct
	err = json.NewDecoder(res.Body).Decode(targetStruct)
	if err != nil {
		return res.StatusCode, fmt.Errorf("failed to decode response body: %v", err)
	}

	// If no errors, return status ok and nil for error
	return http.StatusOK, nil
}

func PostURL(url string, postData interface{}, targetStruct interface{}) (statusCode int, err error) {
	// Convert into JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to convert data to JSON: %v", err)
	}

	// Send POST request
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send POST request: %v", err)
	}

	// Ensure that the body closes at the end of function call
	defer res.Body.Close()

	// Decode the response JSON into target struct
	err = json.NewDecoder(res.Body).Decode(targetStruct)
	if err != nil {
		return res.StatusCode, fmt.Errorf("failed to decode response body: %v", err)
	}

	// If no errors, return status ok and nil for error
	return http.StatusOK, nil
}

func TestGetApi(url string) (statusCode int, err error) {
	// Send get request
	res, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send GET request: %v", err)
	}

	// Return the status code from the API and nil for error
	return res.StatusCode, nil
}

func TestPostApi(url string, data interface{}) (statusCode int, err error) {
	// Convert into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to convert data to JSON: %v", err)
	}

	// Send POST request
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send POST request: %v", err)
	}

	// Return the status code from the API and nil for error
	return res.StatusCode, nil
}
