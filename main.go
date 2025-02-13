package main

import (
	"Assigment-1/config"
	"Assigment-1/handlers"
	"Assigment-1/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	// Start uptime timer›
	utils.StartUptime()
	log.Println("Started uptime timer:", utils.GetUptime())

	// Define a port and check whether the OS have a designated port or not
	PORT := "8080"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// Create a new router
	router := http.NewServeMux()

	// Routes
	router.HandleFunc(http.MethodGet+config.START_URL+"/info/", handlers.InfoHandler) // Root path
	router.HandleFunc(http.MethodGet+config.START_URL+"/info/{two_letter_country_code}", handlers.InfoHandler)
	router.HandleFunc(http.MethodGet+config.START_URL+"/population/", handlers.PopulationHandler) // Root path
	router.HandleFunc(http.MethodGet+config.START_URL+"/population/{two_letter_country_code}", handlers.PopulationHandler)
	router.HandleFunc(http.MethodGet+config.START_URL+"/status/", handlers.StatusHandler) // Root path

	// Listen on the designated port for traffic
	log.Println("Starting server on port " + PORT)
	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
