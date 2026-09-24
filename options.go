package listennotes

import (
	"net/http"
)

// ClientOption allows for options to be passed to the client constructor function
type ClientOption func(c *standardHTTPClient)

// WithHTTPClient uses the caller's HTTP client, including its timeout, redirect,
// proxy, and transport policies. A nil client leaves the SDK defaults intact.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *standardHTTPClient) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBaseURL selects a custom API base URL, including its /api/v2 prefix.
// The client's API key is sent to that server. Invalid URLs fail before sending.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *standardHTTPClient) {
		c.baseURL = baseURL
	}
}
