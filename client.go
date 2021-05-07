// Package listennotes provides an API client to access the listennotes API found here: https://listen-api.listennotes.com/.
// API documentation can be found at https://www.listennotes.com/api/docs/.
package listennotes

import (
	"net/http"
)

// HTTPClient is the client interface
type HTTPClient interface {
	Search(args map[string]string) (*Response, error)
}

type standardHTTPClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ HTTPClient = &standardHTTPClient{}

// NewClient will create a client with reasonable defaults.
// If an apiKey is not provided, the client will use the mock test API by default.
// You can optionally override some configuration.
func NewClient(apiKey string, opts ...ClientOption) HTTPClient {
	baseURL := BaseURLTest
	if apiKey != "" {
		baseURL = BaseURLProduction
	}

	client := &standardHTTPClient{
		apiKey:     apiKey,
		httpClient: defaultHTTPClient,
		baseURL:    baseURL,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *standardHTTPClient) Search(args map[string]string) (*Response, error) {
	return c.get(args, "search")
}
