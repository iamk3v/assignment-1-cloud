package main

import (
	"Assigment-1/constants"
	"Assigment-1/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Define a port and check whether the os have a designated port or not
	PORT := "8080"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// Create a new router
	router := http.NewServeMux()

	// Routes
	router.HandleFunc(http.MethodGet+constants.StartURL+"/info/", handlers.InfoHandler) // Root path
	router.HandleFunc(http.MethodGet+constants.StartURL+"/info/{two_letter_country_code}", handlers.InfoHandler)
	router.HandleFunc(http.MethodGet+constants.StartURL+"/population/", handlers.PopulationHandler) // Root path
	router.HandleFunc(http.MethodGet+constants.StartURL+"/population/{two_letter_country_code}", handlers.PopulationHandler)
	router.HandleFunc(http.MethodGet+constants.StartURL+"/status/", handlers.StatusHandler) // Root path

	// Listen on the designated port for traffic
	log.Println("Starting server on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))

}
