package main

import (
	"Assigment1/constants"
	"Assigment1/handlers"
	"log"
	"net/http"
	"os"
)

func main() {

	PORT := "8080"
	if os.Getenv("PORT") != "" { // Check whether the os have a designated port or not
		PORT = os.Getenv("PORT")
	}

	router := http.NewServeMux()

	router.HandleFunc(http.MethodGet+constants.StartURL+"/info/", handlers.InfoHandler) // Root path
	router.HandleFunc(http.MethodGet+constants.StartURL+"/info/{two_letter_country_code}", handlers.InfoHandler)
	router.HandleFunc(http.MethodGet+constants.StartURL+"/population/", handlers.PopulationHandler) // Root path
	router.HandleFunc(http.MethodGet+constants.StartURL+"/population/{two_letter_country_code}", handlers.PopulationHandler)
	router.HandleFunc(http.MethodGet+constants.StartURL+"/status/", handlers.StatusHandler) // Root path

	log.Println("Starting server on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))

}
