package main

import (
	"covidcase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	COUNTRY = "{country_name:[A-Za-z]+}" // Country name
	WEBID   = "{id}"                     // Webhook id
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

	// CORS
	// source: https://github.com/go-chi/cors
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		Debug:            true,
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Logs the start and end of each request with the elapsed processing time
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes GET
	r.Get("/corona/v1/notifications/", covidcase.HandlerNotifications())
	r.Get("/corona/v1/notifications/"+WEBID, covidcase.HandlerNotification())
	r.Get("/corona/v1/country/"+COUNTRY, covidcase.HandlerCountry()) // optional query parameter "scope" as start/end date
	r.Get("/corona/v1/policy/"+COUNTRY, covidcase.HandlerPolicy())   // optional query parameter "scope" as start/end date
	r.Get("/diag", covidcase.HandlerDiag(appStart))                  // Pass appStart time value for use in this route
	r.Get("/*", covidcase.HandlerLostUser)                           // Route for any other query not handled by API

	// Routes POST
	r.Post("/corona/v1/notifications/", covidcase.HandlerNotifications())

	// Routes DELETE
	r.Delete("/corona/v1/notifications/"+WEBID, covidcase.HandlerNotification())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
