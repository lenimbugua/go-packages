package http_request

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

/*
MakeHTTPRequest creates a new request,
Sets request headers if any,
Sends the request, and
Finaly returns the response.
*/
func MakeHTTPRequest(req Request) (*http.Response, error) {
	// Validate HTTP method
	if !isValidHTTPMethod(req.Method) {
		return nil, fmt.Errorf("invalid HTTP method: %s", req.Method)
	}

	// Create a new request.
	request, err := http.NewRequest(req.Method, req.Url, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Validate and set request headers
	if err := validateAndSetHeaders(request, req.Headers); err != nil {
		return nil, err
	}

	// Use context for the request
	request = request.WithContext(req.Ctx)

	// Make request with timeout
	client := req.Client
	client.Timeout = req.TimeoutDuration
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}
	}()

	return response, nil
}

func isValidHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		return true
	}
	return false
}

func validateAndSetHeaders(request *http.Request, headers requestHeaders) error {
	for key, value := range headers {
		if key == "" || value == "" {
			return fmt.Errorf("invalid header: %s=%s", key, value)
		}
		request.Header.Set(key, value)
	}
	return nil
}
