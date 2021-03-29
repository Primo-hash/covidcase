package utils

/*
Utilities that might be of use to multiple packages
*/

import (
	"encoding/json"
	"errors"
	"net/http"
)

/*
DecodeResponseToMap returns a decoded map from a decoded JSON
* Optional removal of a key in decoded map
*/
func DecodeResponseToMap(data *http.Response, filter string) (map[string]interface{}, error) {
	var result = make(map[string]interface{}) // Body object

	defer data.Body.Close()     // Closing body after finishing read
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
DecodeRequestToMap returns a decoded map from a decoded JSON
* Optional removal of a key in decoded map
*/
func DecodeRequestToMap(data *http.Request, filter string) (map[string]interface{}, error) {
	var result = make(map[string]interface{}) // Body object
	defer data.Body.Close()                   // Closing body after finishing read
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
