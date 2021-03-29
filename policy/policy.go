package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

/*
URL list for 'Policy API' to be modified to query needs
 */
const BASEURL = "https://covidtrackerapi.bsg.ox.ac.uk/api/"
const LATESTURL = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/%s/%s" // URL for latest policy
const SCOPEURL = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/%s/%s" // URL for policy in scope
const ALPHA3URL = "https://restcountries.eu/rest/v2/name/%s" // Retrieves general info about a country

// StringencyInfo struct for JSON encoding HTTP request data
type StringencyInfo struct {
	Country              string  `json:"country"`
	Scope                string  `json:"scope"`
	Stringency           float64 `json:"stringency"`
	Trend            	 float64 `json:"trend"`
}

// Country struct for data extraction of ALPHA-3 code
type Country []struct {
	Alpha3Code     string    `json:"alpha3Code"`
	Region         string    `json:"region"`
}

/*
GetPolicyData returns a map of a decoded json object with
specified trend of a countries' stringency policy based on date (scope) specified.
*/
func GetPolicyData(startDate, endDate, countryName string) (StringencyInfo, error) {
	var result map[string]interface{}
	var stringencyInfo StringencyInfo

	// Get ALPHA3 code of requested country for API request
	alpha3, _, err := GetAlpha3(countryName)
	if err != nil { // Error handling data
		return stringencyInfo, err
	}

	if startDate == "" || endDate == "" { // Format within complete scope
		now := time.Now()
		now = now.AddDate(0, 0, -10)	// Latest values are from 10 days ago
		latestDate := now.Format("2006-01-02") 	// YYYY-MM-DD string

		// Insert parameters into POLICYURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(LATESTURL, alpha3, latestDate))
		if err != nil { // Error handling data
			return stringencyInfo, err
		}
		result, err = DecodeToMap(resData, "")	// Decode for data extraction into StringencyInfo struct
		if err != nil { // Error handling data
			return stringencyInfo, err
		}

		// Inserting and processing data into stringencyInfo struct
		// Remember to use type assertion at the end ".(float64)/.(string)" since the program has to deal with interface{}
		stringencyInfo.Country = countryName				// Country
		stringencyInfo.Scope = "total"						// Scope
		stringencyInfo.Stringency = getStringency(result, "stringency_actual").(float64)	// Stringency value
		stringencyInfo.Trend = 0

		return stringencyInfo, nil
	} else {							  // Format within scope of date specified
		// Insert parameters into POLICYURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(SCOPEURL, startDate, endDate))
		if err != nil { // Error handling data
			return stringencyInfo, err
		}
		result, err = DecodeToMap(resData, "")	// Decode for data extraction into StringencyInfo struct
		if err != nil { // Error handling data
			return stringencyInfo, err
		}

		// extract from data key
		stringencyData := result["data"].(map[string]interface{})
		// Get stringency values from start date
		startDateStringency, _ := getStringencyScope(stringencyData, startDate, alpha3, "stringency_actual")
		// Get stringency values from end date
		endDateStringency, _ := getStringencyScope(stringencyData, endDate, alpha3, "stringency_actual")

		// Inserting and processing data into stringencyInfo struct
		// Remember to use type assertion at the end ".(float64)/.(string)" since the program has to deal with interface{}
		stringencyInfo.Country = countryName							    					// Country
		stringencyInfo.Scope = "total"									    					// Scope
		stringencyInfo.Stringency = endDateStringency.(float64)									// Stringency
		if stringencyInfo.Stringency == -1 {	// Check if missing information, set to 0 if so
			stringencyInfo.Trend = 0
		} else {
			stringencyInfo.Trend = endDateStringency.(float64) - startDateStringency.(float64)	// Trend
		}

		return stringencyInfo, nil
	}
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

/*
GetAlpha3 returns a string of specified Country's ALPHA-3 code and Continent
*/
func GetAlpha3(countryName string) (string, string, error) {
	var countries Country // Holds JSON object values, query might return multiple countries

	// Insert parameters into CASEURL for HTTP GET request
	resData, err := http.Get(fmt.Sprintf(ALPHA3URL, countryName))
	if err != nil { // Error handling data
		return "", "", err
	}

	// Decode into countries object
	err = json.NewDecoder(resData.Body).Decode(&countries)
	if err != nil {
		return "", "", err
	}
	// Extract ALPHA-3 from first country
	alpha := countries[0].Alpha3Code
	continent := countries[0].Region
	return alpha, continent, nil
}

/*
HealthCheck returns an http status code after checking for a response from exchangeratesAPI servers
*/
func HealthCheck() (string, error) {
	// Send HTTP GET request
	resData, err := http.Get(BASEURL)
	if err != nil { // Error handling HTTP request
		return "", err
	}
	return resData.Status, nil
}

/*
getStringency returns a value 'stringency_actual' or 'stringency' if it exists
*/
func getStringency(data map[string]interface{}, key string) interface{} {
	res := data["stringencyData"].(map[string]interface{})[key]
	return res
}

/*
getStringencyScope recursive func that returns a value 'stringency_actual' or 'stringency'
if it exists for the scope parameter
*/
func getStringencyScope(data map[string]interface{}, date, alpha3, key string) (interface{}, bool) {
	res, ok := data[date].(map[string]interface{})[alpha3].(map[string]interface{})[key]	// First key
	if !ok {	// incase first key doesn't work try second
		res, ok = getStringencyScope(data, date, alpha3, "stringency")
		if !ok {
			return -1, true		// Return -1 if no keys found
		} else {
			return res, true   // Return second key if found
		}
	}
	return res, true
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
