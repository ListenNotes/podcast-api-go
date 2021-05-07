package listennotes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is the standard response for all client functions
type Response struct {
	Stats ResponseStatistics
	Data  map[string]interface{}
}

func (c *standardHTTPClient) execute(args map[string]string, path string) (*Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", path, err)
	}
	req.Header.Add(RequestHeaderKeyAPI, c.apiKey)

	q := req.URL.Query()
	for k, v := range args {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to executing request to %s: %w", path, err)
	}
	defer resp.Body.Close()

	// map any generic status code errors
	if mappedError, ok := errMap[resp.StatusCode]; ok && mappedError != nil {
		return nil, mappedError
	}

	// generic body parsing
	var genericJSON map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&genericJSON); err != nil {
		return nil, fmt.Errorf("failed parsing the response from %s: %w", path, err)
	}

	// gather the header statistics
	stats := parseStats(resp)

	return &Response{
		Stats: stats,
		Data:  genericJSON,
	}, nil
}
