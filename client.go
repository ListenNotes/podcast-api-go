// Package listennotes provides a client for the Listen Notes Podcast API.
// See https://www.listennotes.com/api/docs/ for the API reference.
package listennotes

import "net/http"

type standardHTTPClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ HTTPClient = &standardHTTPClient{}

// NewClient creates a client with a 30-second timeout and redirects disabled.
// An empty API key selects the public mock API; otherwise it selects production.
// Reuse a client to reuse HTTP connections. Configure options before sharing it.
func NewClient(apiKey string, opts ...ClientOption) HTTPClient {
	baseURL := BaseURLTest
	if apiKey != "" {
		baseURL = BaseURLProduction
	}
	client := &standardHTTPClient{apiKey: apiKey, httpClient: newHTTPClient(), baseURL: baseURL}
	for _, opt := range opts {
		if opt != nil {
			opt(client)
		}
	}
	return client
}
