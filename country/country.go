package country

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alediaferia/gocountries"
	"net/http"
)

/*
URL list for 'REST Countries API' to be modified to query needs
*/
const BASEURL = "https://covid-api.mmediagroup.fr/v1/cases" // For healthchecks
const CASEURL = "https://covid-api.mmediagroup.fr/v1/cases?country=%s" // For all covid cases
const SCOPEURL = "https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Confirmed" // Cases within a date scope

// CaseInfo struct for JSON encoding HTTP request data
type CaseInfo struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            float64 `json:"confirmed"`
	Recovered            float64 `json:"recovered"`
	PopulationPercentage string `json:"population_percentage"`
}

/*
GetCountryData returns a map of a decoded json object with
specified total confirmed cases and total of recovered based on a timescope(date) specified
*/
func GetCountryData(startDate, endDate, countryName string) (CaseInfo, error) {
	var result map[string]interface{}
	var caseInfo CaseInfo

	if startDate == "" || endDate == "" { // Format within complete scope
		// Insert parameters into CASEURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(CASEURL, countryName))
		if err != nil { // Error handling data
			return caseInfo, err
		}
		result, err = DecodeToMap(resData, "")	// Decode for data extraction into Caseinfo struct
		if err != nil { // Error handling data
			return caseInfo, err
		}

		// Inserting and processing data into caseInfo struct
		// Remember to use type assertion at the end ".(float64)/.(string)" since the program has to deal with interface{}
		caseInfo.Country = result["All"].(map[string]interface{})["country"].(string)			// Country
		caseInfo.Continent = result["All"].(map[string]interface{})["continent"].(string)		// Continent
		caseInfo.Scope = "total"												                // Scope
		caseInfo.Confirmed = result["All"].(map[string]interface{})["confirmed"].(float64)		    // Confirmed cases
		caseInfo.Recovered = result["All"].(map[string]interface{})["recovered"].(float64)		    // Recovered cases
		// Percentage of population with a confirmed case
		percentage := caseInfo.Confirmed / result["All"].(map[string]interface{})["population"].(float64) * 100
		caseInfo.PopulationPercentage = fmt.Sprintf("%.2f", percentage)

		return caseInfo, nil
	} else {							  // Format within scope of date specified

		// Insert parameters into SCOPEURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(SCOPEURL, countryName))			// Confirmed cases
		if err != nil { // Error handling data
			return caseInfo, err
		}
		// Processing scope
		result, err = DecodeToMap(resData, "")	// Decode for data extraction into Caseinfo struct
		if err != nil { // Error handling data
			return caseInfo, err
		}

		// Extracting confirmed cases at start date and end date for scope calculation
		startDateCases := result["All"].(map[string]interface{})["dates"].(map[string]interface{})[startDate].(float64)
		endDateCases := result["All"].(map[string]interface{})["dates"].(map[string]interface{})[endDate].(float64)

		// Inserting data into caseInfo struct
		// Remember to use type assertion at the end ".(float)/.(string)" since the program has to deal with interface{}
		caseInfo.Country = result["All"].(map[string]interface{})["country"].(string)	    // Country
		caseInfo.Continent = result["All"].(map[string]interface{})["continent"].(string)	// Continent
		caseInfo.Scope = startDate + "-" + endDate											// Scope
		caseInfo.Confirmed = endDateCases - startDateCases	// Confirmed cases
		caseInfo.Recovered = 0		    // No recovery cases response
		// Percentage of population with a confirmed case
		percentage := caseInfo.Confirmed / result["All"].(map[string]interface{})["population"].(float64) * 100
		caseInfo.PopulationPercentage = fmt.Sprintf("%.2f", percentage)

		return caseInfo, nil
	}
}

/*
GetCurrency returns a string of specified Country's currency code e.g.(NOK, USD, EUR...)
*/
func GetCurrency(countryName string) (string, error) {
	// Query for structs of possible countries
	countries, err := gocountries.CountriesByName(countryName)
	// Extract first country
	c := (countries)[0]
	// Extract policy code
	currencyCode := c.Currencies[0]
	return currencyCode, err
}

/*
HealthCheck returns an http status code after checking for a response from REST Countries API servers
*/
func HealthCheck() (string, error) {
	// Using http API for restcountriesAPI because gocountries pckg does not support searching by country code
	// Send HTTP GET request
	resData, err := http.Get(BASEURL)
	if err != nil { // Error handling HTTP request
		return "", err
	}
	return resData.Status, nil
}

/*
Decode returns a decoded map from a decoded JSON
* Optional removal of a key in decoded map
*/
func DecodeToMap(data *http.Response, filter string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})		// Body object

	defer data.Body.Close() // Closing body after finishing read
	if data.StatusCode != 200 { // Error handling HTTP request
		e := errors.New(data.Status)
		return nil, e
	}
	// Decoding body
	err := json.NewDecoder(data.Body).Decode(&result)
	if err != nil { // Error handling decoding
		return nil, err
	}
	// Optional filtering of certain key in map
	if filter != "" {
		// Check for filter word existence
		_, ok := result[filter]
		if ok {
			delete(result, filter)
		}
	}
	// Return map with requested data
	return result, err
}
