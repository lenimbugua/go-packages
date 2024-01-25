package http_request

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type requestHeaders map[string]string

type Request struct {
	Ctx             context.Context
	Client          *http.Client
	Method          string
	Url             string
	Body            io.Reader
	Headers         requestHeaders
	TimeoutDuration time.Duration
}


// MakeHTTPRequest makes an HTTP request using the provided request object.
// It validates the HTTP method, creates a new request, sets the headers, uses a context for the request,
// and makes the request with a timeout.
//
// Inputs:
//   req (Request): The request object containing the method, URL, body, headers, context, client, and timeout duration.
//
// Outputs:
//   response (*http.Response): A pointer to the HTTP response object.
//   err (error): An error, if any occurred during the request.
	

// Custom error types
type InvalidHTTPMethodError struct {
	Method string
}

func (e InvalidHTTPMethodError) Error() string {
	return fmt.Sprintf("invalid HTTP method: %s", e.Method)
}

type RequestCreationError struct {
	Err error
}

func (e RequestCreationError) Error() string {
	return fmt.Sprintf("failed to create request: %v", e.Err)
}

type InvalidHeaderError struct {
	Err error
}

func (e InvalidHeaderError) Error() string {
	return fmt.Sprintf("invalid request header: %v", e.Err)
}

type RequestError struct {
	Err error
}

func (e RequestError) Error() string {
	return fmt.Sprintf("failed to make request: %v", e.Err)
}

// MakeHTTPRequest makes an HTTP request
func MakeHTTPRequest(req Request) (*http.Response, error) {
	// Validate HTTP method
	if !isSupportedHTTPMethod(req.Method) {
		return nil, InvalidHTTPMethodError{Method: req.Method}
	}

	// Create a new request.
	request, err := http.NewRequest(req.Method, req.Url, req.Body)
	if err != nil {
		return nil, RequestCreationError{Err: err}
	}

	// Validate and set request headers
	if err := validateAndSetHeaders(request, req.Headers); err != nil {
		return nil, InvalidHeaderError{Err: err}
	}

	// Use context for the request
	request = request.WithContext(req.Ctx)

	// Make request with timeout using http.Client and context.WithTimeout
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
		Timeout:   req.TimeoutDuration,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, RequestError{Err: err}
	}
	
	return response, nil
}


// isValidHTTPMethod checks if the given HTTP method is valid.
//
// method string
// bool
var validMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodPatch:   true,
	http.MethodTrace:   true,
}

// isSupportedHTTPMethod checks if the method exists in the validMethods map.
//
// It takes a method string as a parameter and returns a boolean.
func isSupportedHTTPMethod(method string) bool {
	// Check if the method exists in the validMethods map
	_, ok := validMethods[method]
	return ok
}

// validateAndSetHeaders validates and sets the headers for the given HTTP request.
//
// request *http.Request - the HTTP request to set headers for
// headers requestHeaders - the headers to validate and set
// error - an error if the headers are invalid, otherwise nil
func validateAndSetHeaders(request *http.Request, headers requestHeaders) error {
	for key, value := range headers {
		if key == "" || value == "" || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return fmt.Errorf("invalid header: %s=%s", key, value)
		}
		escapedValue := url.QueryEscape(value)
		request.Header.Add(http.CanonicalHeaderKey(key), escapedValue)
	}
	return nil
}
