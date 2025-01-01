package hyperliquid

import (
	"encoding/json"
	"fmt"
)

// API implementation general error
type APIError struct {
	Message string
}

func (e APIError) Error() string {
	return e.Message
}

// IAPIService is an interface for making requests to the API Service.
//
// It has a Request method that takes a path and a payload and returns a byte array and an error.
// It has a debug method that takes a format string and args and returns nothing.
// It has an Endpoint method that returns a string.
type IAPIService interface {
	debug(format string, args ...interface{})
	Request(path string, payload any) ([]byte, error)
	Endpoint() string
	KeyManager() *PKeyManager
}

// MakeUniversalRequest is a generic function that takes an
// IAPIService and a request and returns a pointer to the result and an error.
// It makes a request to the API Service and unmarshals the result into the result type T
func MakeUniversalRequest[T any](api IAPIService, request any) (*T, error) {
	if api.Endpoint() == "" {
		return nil, APIError{Message: "Endpoint not set"}
	}
	if api == nil {
		return nil, APIError{Message: "API not set"}
	}
	if api.Endpoint() == "/exchange" && api.KeyManager() == nil {
		return nil, APIError{Message: "API key not set"}
	}

	response, err := api.Request(api.Endpoint(), request)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal(response, &result)
	if err == nil {
		return &result, nil
	}

	var errResult map[string]interface{}
	err = json.Unmarshal(response, &errResult)
	if err != nil {
		api.debug("Error second json.Unmarshal: %s", err)
		return nil, APIError{Message: "Unexpected response"}
	}

	if errResult["status"] == "err" {
		return nil, APIError{Message: errResult["response"].(string)}
	}

	return nil, APIError{Message: fmt.Sprintf("Unexpected response: %v", errResult)}
}
