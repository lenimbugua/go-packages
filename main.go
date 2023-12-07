package util

import (
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type Headers map[string]string

func MakeHTTPRequest(client *http.Client, logger *zerolog.Logger, httpMethod string, url string, body io.Reader, headers Headers) ([]byte, error) {
	//Create a new request.
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		logger.Error().Err(err).Msg("Could not create new request")
		return nil, err
	}

	//Set request headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	//Make request
	response, err := client.Do(request)
	if err != nil {
		logger.Error().Err(err).Msg("Could not send request")
		return nil, err
	}

	// Close response body
	defer func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				logger.Error().Err(err).Msg("Could not close body stream")
			}
		}
	}()

	// Read the response body
	res, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Could read response body")
		return nil, err
	}
	return res, nil
}
