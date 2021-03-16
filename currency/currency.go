package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/*
URL list for 'Foreign exchange rates API' to be modified to query needs
order: (startDate, endDate, currencyCode, currencyBase)
NOTE: 'symbols' and 'base' parameters can be omitted as query will return default values for these
 */
const BASEURL = "https://api.exchangeratesapi.io/"
const HISTORYURL = "https://api.exchangeratesapi.io/history?start_at=%s&end_at=%s&symbols=%s&base=%s"
const LATESTURL = "https://api.exchangeratesapi.io/latest?&symbols=%s&base=%s"	// URL for latest currency

/*
GetExchangeData returns a map of a decoded json object with
specified history of exchange rates based on date specified
or specify multiple countries' currencies and a base currency for comparison.
* Date parameters are mandatory when querying currency history.
* CurrencyCode is optional for history.
* CurrencyBase is optional.
*/
func GetExchangeData(startDate, endDate, currencyCode, currencyBase string) (map[string]interface{}, error) {
	if startDate == "" || endDate == "" { // Request for currency
		// Insert parameters into CURRENCYURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(LATESTURL, currencyCode, currencyBase))
		if err != nil { // Error handling data
			return nil, err
		}
		return Decode(resData, "date")
	} else {							  // Request for history
		// Insert parameters into HISTORYURL for HTTP GET request
		resData, err := http.Get(fmt.Sprintf(HISTORYURL, startDate, endDate, currencyCode, currencyBase))
		if err != nil { // Error handling data
			return nil, err
		}
		return Decode(resData, "")
	}
}

// TODO make this function return a struct instead of a map, no random placement
// of JSON arrays and keys like map handles it

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
