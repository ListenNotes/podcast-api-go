package listennotes

import (
	"net/http"
)

// ClientOption allows for options to be passed to the client constructor function
type ClientOption func(c *standardHTTPClient)

// WithHTTPClient allows providing an underlying http client.  It is good practice to _not_ use the default http client
// that Go provides as it has no timeouts.  If you do not provide your own default client, a reasonable one will be created for you.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *standardHTTPClient) {
		c.httpClient = httpClient
	}
}

// WithBaseURL allows for providing a custom base URL.  If not provided a reasonable url will be selected based on your apiKey.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *standardHTTPClient) {
		c.baseURL = baseURL
	}
}
