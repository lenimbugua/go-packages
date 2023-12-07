package http_request

import (
	"io"
	"net/http"
)

type requestHeaders map[string]string

/*
  MakeHTTPRequest creates a new request,
  Sets request headers if any,
  Sends the request, and
  Finaly returns the response
*/

func MakeHTTPRequest(client *http.Client, httpMethod string, url string, body io.Reader, headers requestHeaders) ([]byte, error) {
	//Create a new request.
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}

	//Set request headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	//Make request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Close response body
	defer func() error {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				return err
			}
		}
		return nil
	}()

	// Read the response body
	res, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}
