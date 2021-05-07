// Package listennotes provides an API client to access the listennotes API found here: https://listen-api.listennotes.com/.
// API documentation can be found at https://www.listennotes.com/api/docs/.
package listennotes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient is the client interface
type HTTPClient interface {
	Search(args map[string]string) (*Response, error)
	FetchPodcastByID(id string, args map[string]string) (*Response, error)
	BatchFetchPodcastByID(ids []string, args map[string]string) (*Response, error)
	DeletePodcast(id string, args map[string]string) (*Response, error)
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
	return c.get("search", args)
}

// GET /typeahead

// GET /best_podcasts

func (c *standardHTTPClient) FetchPodcastByID(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("podcasts/%s", id), args)
}

// GET /episodes/{id}

// POST /episodes
func (c *standardHTTPClient) BatchFetchPodcastByID(ids []string, args map[string]string) (*Response, error) {
	return c.post("episodes", args, url.Values{
		"ids": []string{strings.Join(ids, ",")},
	})
}

// POST /podcasts

// GET /curated_podcasts/{id}

// GET /genres

// GET /regions

// GET /languages

// GET /just_listen

// GET /curated_podcasts

// GET /podcats/{id}/recommendations

// GET /episodes/{id}/recommendations

// GET /playlists

// GET /playlists/{id}

// POST /podcasts/submit

// -H 'Content-Type: application/x-www-form-urlencoded'\
// -d 'rss=https://feeds.megaphone.fm/committed'\
// -d 'email=hello@example.com'

func (c *standardHTTPClient) DeletePodcast(id string, args map[string]string) (*Response, error) {
	return c.delete(fmt.Sprintf("podcasts/%s", id), args)
}
