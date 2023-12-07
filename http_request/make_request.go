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

/*
  MakeHTTPRequest creates a new request,
  Sets request headers if any,
  Sends the request, and
  Finaly returns the response.
*/

func MakeHTTPRequest(ctx context.Context, client *http.Client, httpMethod string, url string, body io.Reader, headers requestHeaders, timeoutDuration time.Duration) ([]byte, error) {
	// Validate HTTP method
	if !isValidHTTPMethod(httpMethod) {
		return nil, fmt.Errorf("invalid HTTP method: %s", httpMethod)
	}

	// Create a new request.
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Validate and set request headers
	if err := validateAndSetHeaders(request, headers); err != nil {
		return nil, err
	}

	// Use context for the request
	request = request.WithContext(ctx)

	// Make request with timeout
	client.Timeout = timeoutDuration 
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

	// Read the response body
	res, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return res, nil
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
