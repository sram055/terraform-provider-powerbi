package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPUnsuccessfulError represents an error thrown when a non 2xx response is received
type HTTPUnsuccessfulError struct {
	Request   *http.Request
	Response  *http.Response
	ErrorBody *ErrorBody
}

// ErrorBody represents the error returend in the body of Power BI API requests
type ErrorBody struct {
	Code    string
	Message string
}

type roundTripperErrorOnUnsuccessful struct {
	rt http.RoundTripper
}

func (err HTTPUnsuccessfulError) Error() string {

	message := fmt.Sprintf("status code '%s'", err.Response.Status)
	if err.ErrorBody != nil {
		message += fmt.Sprintf(" with code '%s' and message '%s'", err.ErrorBody.Code, err.ErrorBody.Message)
	}
	return message
}

func (h roundTripperErrorOnUnsuccessful) RoundTrip(req *http.Request) (*http.Response, error) {

	resp, err := h.rt.RoundTrip(req)

	if err != nil || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return resp, err
	}

	// try and read the body to get the formatted error
	var errorBody ErrorBody
	if resp.Body != http.NoBody {
		respBody, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(respBody, &errorBody)
	}

	return resp, HTTPUnsuccessfulError{
		Request:   req,
		Response:  resp,
		ErrorBody: &errorBody,
	}
}