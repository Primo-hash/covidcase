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
const CASEURL = "https://covid-api.mmediagroup.fr/v1/cases?country=%s" // For all cases
const CONFSCOPEURL = "https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Confirmed" // Confirmed history
const RECSCOPEURL = "https://covid-api.mmediagroup.fr/v1/history?country=%s&status=Recovered" // Recovered history

// Diagnose struct for JSON encoding
type CaseInfo struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage string `json:"population_percentage"`
}

// Country struct for JSON encoding
type Country []struct {
	Name           string    `json:"name"`
	TopLevelDomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	CallingCodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	AltSpellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	NativeName     string    `json:"nativeName"`
	NumericCode    string    `json:"numericCode"`
	Currencies     []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Iso6391    string `json:"iso639_1"`
		Iso6392    string `json:"iso639_2"`
		Name       string `json:"name"`
		NativeName string `json:"nativeName"`
	} `json:"languages"`
	Translations struct {
		De string `json:"de"`
		Es string `json:"es"`
		Fr string `json:"fr"`
		Ja string `json:"ja"`
		It string `json:"it"`
		Br string `json:"br"`
		Pt string `json:"pt"`
		Nl string `json:"nl"`
		Hr string `json:"hr"`
		Fa string `json:"fa"`
	} `json:"translations"`
	Flag          string `json:"flag"`
	RegionalBlocs []struct {
		Acronym       string        `json:"acronym"`
		Name          string        `json:"name"`
		OtherAcronyms []interface{} `json:"otherAcronyms"`
		OtherNames    []interface{} `json:"otherNames"`
	} `json:"regionalBlocs"`
	Cioc string `json:"cioc"`
}

/*
GetCountryData returns a map of a decoded json object with
specified total confirmed cases and total of recovered based on a timescope(date) specified
*/
func GetCountryData(startDate, endDate, countryName string) (map[string]interface{}, error) {
	var result map[string]interface{}
	var caseInfo CaseInfo

	if startDate == "" || endDate == "" { // Format within complete scope
		// Insert parameters into CASEURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(CASEURL, countryName))
		if err != nil { // Error handling data
			return nil, err
		}
		result, err = Decode(resData, "")	// Decode for data extraction into Caseinfo struct

		// Inserting and processing data into caseInfo struct
		caseInfo.Country = result["All"].(map[string]string)["country"]			// Country
		caseInfo.Continent = result["All"].(map[string]string)["continent"]		// Continent
		caseInfo.Scope = "total"												// Scope
		caseInfo.Confirmed = result["All"].(map[string]int)["confirmed"]		// Confirmed cases
		caseInfo.Recovered = result["All"].(map[string]int)["recovered"]		// Recovered cases
		// Percentage of population with a confirmed case
		percentage := float64(caseInfo.Confirmed) / float64(result["All"].(map[string]int)["population"])
		caseInfo.PopulationPercentage = fmt.Sprintf("%.2f", percentage)	//

		return Decode(resData, "date")
	} else {							  // Format within scope of date specified
		return nil, nil
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
	// Extract currency code
	currencyCode := c.Currencies[0]
	return currencyCode, err
}

/*
GetNeighbour returns a string of specified Country's Neighbours' currency codes
* limit parameter is for restricting the amount of currencies returned
*/
func GetNeighbour(countryName string, limit int) (string, error) {
	var borderURL = BASEURL // URL string for modification
	var countries Country // Holds JSON object values

	// Query for structs of possible countries
	country, err := gocountries.CountriesByName(countryName)
	// Extract first country
	c := country[0]
	// Extract border alpha codes
	neighbourAlpha := c.Borders[:]
	// To avoid indexing out of neighbourAlpha's range
	if limit > len(neighbourAlpha) {limit = len(neighbourAlpha)}
	// parse neighbour alpha codes and append to API call URL
	for i:= 0; i < limit; i++ {
		if i >= limit {
			// Nothing happens if index exceeds limit
		} else if i == limit - 1  {		// If last element in array
			borderURL += neighbourAlpha[i] // Avoid appending with ';' at the end
		} else if i < limit { 				// If not last element in array
			borderURL += neighbourAlpha[i] + ";"
		}
	}
	// Using http API for restcountriesAPI because gocountries pckg does not support searching by country code
	// Send HTTP GET request
	resData, err := http.Get(borderURL)
	if err != nil { // Error handling HTTP request
		return "", err
	}
	defer resData.Body.Close() // Closing body after finishing read
	// Decoding body
	err = json.NewDecoder(resData.Body).Decode(&countries)
	if err != nil {
		fmt.Println("Decoding: " + err.Error())
		return "", err
	}
	// Make string value of neighbour country currencies for return
	currencyCodes := ""
	for i, a := range countries {
		if i != len(countries) - 1 {		// If not last element in array
			currencyCodes += a.Currencies[0].Code + ","
		} else {
			currencyCodes += a.Currencies[0].Code // Avoid appending with ',' at the end
		}
	}
	return currencyCodes, nil
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
func Decode(data *http.Response, filter string) (map[string]interface{}, error) {
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