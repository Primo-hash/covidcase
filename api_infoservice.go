package covidcase

import (
	"covidcase/country"
	"covidcase/policy"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DATELEN = 21	// Length of startDate and endDate optional param characters
var appStart time.Time // Uptime variable

// Diagnose struct for JSON encoding
type Diagnose struct {
	Mmediagroupapi   string `json:"mmediagroupapi"`
	Covidtrackerapi  string `json:"covidtrackerapi"`
	Registered		 float64`json:"registered"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

/*
HandlerLostUser is a function for guiding lost souls back to the right 'relative' path
*/
func HandlerLostUser(w http.ResponseWriter, r *http.Request) {
	protocol := "http://"
	host := r.Host 					// Host URL
	pathAPI := "/exchange/v1/" 		// API path

	// API request queries
	pathHistory := "exchangehistory/norway/2020-03-01-2020-03-03"
	pathBorder := "exchangeborder/norway?limit=2"
	pathDiag := "diag/"

	line1 := "Hello! You seem lost! Let me help you!"
	line2 := "These are some examples:"
	history := protocol + host + pathAPI + pathHistory
	border := protocol + host + pathAPI + pathBorder
	diagnose := protocol + host + pathAPI + pathDiag

	// HTML form for response such that URLs are hyperlinks
	var form = `<p>`+line1+`</p>
				<p>`+line2+`</p>
			    <p><a href="`+history+`">`+history+`</a></p>
				<p><a href="`+border+`">`+border+`</a></p>
				<p><a href="`+diagnose+`">`+diagnose+`</a></p>`

	// Generate HTML template from Form
	res := template.New("table")
	res, err := res.Parse(form)
	if err != nil {
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		fmt.Println("Could not parse HTML form: " + err.Error())
	}

	// Write template
	err = res.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		fmt.Println("Could not write back response: " + err.Error())
	}
}

// HandlerCountry main handler for route related to `/country` requests
func HandlerCountry() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleCountryGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}

// HandlerPolicy main handler for route related to `/policy` requests
func HandlerPolicy() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlePolicyGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}

// HandlerDiag main handler for route related to `/diag` requests
func HandlerDiag(t time.Time) func(http.ResponseWriter, *http.Request) {
	appStart = t	// Pass application start time for multiple function access
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleDiagGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}

// handleCountryGet utility function, package level, to handle GET request to country route
func handleCountryGet(w http.ResponseWriter, r *http.Request) {
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 5 || parts[3] != "country" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// extract URL parameters
	countryName := p(r, "country_name")
	// Handling case sensitivity of country name (lowercase all letters then capitalize first letter)
	countryName = strings.ToLower(countryName) // All letters lower case
	countryName = strings.Title(countryName) // First letter capitalized

	// Extract optional 'scope' parameter
	scope := r.URL.Query().Get("scope")
	// Extract start and end date from scope
	sDate, eDate := split(scope, "-", 3)
	// Request covid info for queried country

	result, err := country.GetCountryData(sDate, eDate, countryName)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}

// handlePolicyGet utility function, package level, to handle GET request to policy route
func handlePolicyGet(w http.ResponseWriter, r *http.Request) {
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 5 || parts[3] != "policy" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// extract URL parameters
	countryName := p(r, "country_name")
	// Handling case sensitivity of country name (lowercase all letters then capitalize first letter)
	countryName = strings.ToLower(countryName) // All letters lower case
	countryName = strings.Title(countryName) // First letter capitalized

	// Extract optional 'scope' parameter
	scope := r.URL.Query().Get("scope")
	// Extract start and end date from scope
	sDate, eDate := split(scope, "-", 3)
	// Request covid info for queried country

	result, err := policy.GetPolicyData(sDate, eDate, countryName)
	if err != nil { // Error handling bad request parameter for params
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}

// handleDiagGet utility function, package level, to handle GET request to diag route
func handleDiagGet(w http.ResponseWriter, r *http.Request) {
	var diag Diagnose
	var err error
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 2 || parts[1] != "diag" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// Insert covidetrackerapi status code
	diag.Covidtrackerapi, err = policy.HealthCheck()
	if err != nil {
		// Error could be a 400, print internally as well
		fmt.Println("HTTP status: " + err.Error())
	}
	// Insert mmediagroupapi status code
	diag.Mmediagroupapi, err = country.HealthCheck()
	if err != nil {
		// Error could be a 400, print internally as well
		fmt.Println("HTTP status: " + err.Error())
	}
	// Insert number of registered webhooks
	diag.Registered = 0 // TODO update to fetch number of webhooks!!!
	// Insert API version
	diag.Version = "v1"
	// Insert API uptime in hr min sec
	diag.Uptime = time.Since(appStart).String()
	// Encode diagnostic report
	report, _ := json.Marshal(diag)
	if err != nil {
		// In case of no server response, reply with 500
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("Encode: " + err.Error())
	}
	// Send status and diagnostic report
	w.WriteHeader(http.StatusOK)
	w.Write(report)		 		 // Send result for processing
}

// p is a shortened function for extracting URL parameters
func p(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// split uses strings.Split and strings.Join to separate a string on the Nth occurrence of a character and returns
// both parts of the string
func split(s, sep string, pos int) (string, string) {
	// Check if date entered has a valid length of 21 chars
	if len(s) != DATELEN {
		return "", "" // Empty return
	} else {
		str := strings.Split(s, sep)
		// Join seperated parts using sep character and return both split ends of the string
		return strings.Join(str[:pos], sep), strings.Join(str[pos:], sep)
	}
}

// getLimit converts string number into int and returns int, for limiting option
func getLimit(s string) int {
	// convert string number to an int and handle error for non digit characters
	if n, err := strconv.Atoi(s); err == nil {
		return n
	} else {
		// return '20' to fixed limit
		return 20
	}
}

// resWithData write objects/types encoded as a JSON to http response
func resWithData(w io.Writer, response interface{}) {
	// handle JSON objects
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}