package main

import (
	"covidcase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

/*
	const definitions for chi routing regex and Query string parameters
 */
const (
	// Chi regex parameters
	COUNTRY = "{country_name:[A-Za-z]+}"			  // Country name
	//BY = "{b_year:\\d\\d\\d\\d}"		   		  // Begin year
	//BM = "{b_month:\\d\\d}"		   	  			  // Begin month
	//BD = "{b_day:\\d\\d}"		   		          // Begin day
	//EY = "{e_year:\\d\\d\\d\\d}"		   		  // End year
	//EM = "{e_month:\\d\\d}"		   	              // End month
	//ED = "{e_day:\\d\\d}"		   		          // End day
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Define application startup time value
	appStart := time.Now()

	// Define new router
	r := chi.NewRouter()

	// Logs the start and end of each request with the elapsed processing time
	r.Use(middleware.Logger)

	// Routes
	r.Get("/corona/v1/country/"+COUNTRY, covidcase.HandlerCountry()) // optional query parameter "scope" as start/end date
	r.Get("/exchange/v1/exchangeborder/"+COUNTRY, covidcase.HandlerBorder())
	r.Get("/exchange/v1/diag/", covidcase.HandlerDiag(appStart))	// Pass appStart time value for use in this route
	r.Get("/*", covidcase.HandlerLostUser) // Route for any other query not handled by API

	log.Fatal(http.ListenAndServe(":"+port, r))
}