// Package podcast provides an API client to access the listennotes API found here: https://listen-api.listennotes.com/.
// API documentation can be found at https://www.listennotes.com/api/docs/.
package podcast // import "github.com/ListenNotes/podcast"

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Base urls for access the available api endpoints
const (
	BaseURLProduction = "https://listen-api.listennotes.com/api/v2"
	BaseURLTest       = "https://listen-api-test.listennotes.com/api/v2"
)

// Request header keys
const (
	RequestHeaderKeyAPI = "X-ListenAPI-Key"
)

// Reponse header keys
const (
	ResponseHeaderKeyFreeQuota       = "X-ListenAPI-FreeQuota"
	ResponseHeaderKeyUsage           = "X-ListenAPI-Usage"
	ResponseHeaderKeyLatencySeconds  = "X-listenAPI-Latency-Seconds"
	ResponseHeaderKeyNextBillingDate = "X-Listenapi-NextBillingDate"
)

// TimeFormat is the string format of all response times
const TimeFormat = "2006-01-02T15:04:05.000000+07:00"

var defaultHTTPClient *http.Client = &http.Client{}

type Response struct {
	Stats ResponseStatistics
	Data  map[string]interface{}
}

type HTTPClient interface{}

type StandardHTTPClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ HTTPClient = &StandardHTTPClient{}

// NewClient will create a client with reasonable defaults.
// If an apiKey is not provided, the client will use the mock test API by default.
// You can optionally override some configuration.
func NewClient(apiKey string, opts ...ClientOption) HTTPClient {
	baseURL := BaseURLTest
	if apiKey != "" {
		baseURL = BaseURLProduction
	}

	client := &StandardHTTPClient{
		apiKey:     apiKey,
		httpClient: defaultHTTPClient,
		baseURL:    baseURL,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *StandardHTTPClient) Search(args map[string]string) (*Response, error) {

	// TODO: move all this common stuff to a function
	url := fmt.Sprintf("%s/search", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add(RequestHeaderKeyAPI, c.apiKey)

	q := req.URL.Query()
	for k, v := range args {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to executing request: %w", err)
	}
	defer resp.Body.Close()

	// map any generic status code errors
	if mappedError, ok := errMap[resp.StatusCode]; ok && mappedError != nil {
		return nil, mappedError
	}

	// generic body parsing
	var genericJSON map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&genericJSON); err != nil {
		return nil, fmt.Errorf("failed parsing the response: %w", err)
	}

	// gather the header statistics
	stats := parseStats(resp)

	return &Response{
		Stats: stats,
		Data:  genericJSON,
	}, nil
}
