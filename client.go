// Package listennotes provides an API client to access the listennotes API found here: https://listen-api.listennotes.com/.
// API documentation can be found at https://www.listennotes.com/api/docs/.
package listennotes

import (
	"fmt"
	"net/http"
	"net/url"
)

// HTTPClient is the client interface
type HTTPClient interface {
	Search(args map[string]string) (*Response, error)
	Typeahead(args map[string]string) (*Response, error)
	SpellCheck(args map[string]string) (*Response, error)
	FetchRelatedSearches(args map[string]string) (*Response, error)
	FetchTrendingSearches(args map[string]string) (*Response, error)
	FetchBestPodcasts(args map[string]string) (*Response, error)
	FetchPodcastByID(id string, args map[string]string) (*Response, error)
	FetchEpisodeByID(id string, args map[string]string) (*Response, error)
	BatchFetchEpisodes(args map[string]string) (*Response, error)
	BatchFetchPodcasts(args map[string]string) (*Response, error)
	FetchCuratedPodcastsListByID(id string, args map[string]string) (*Response, error)
	FetchPodcastGenres(args map[string]string) (*Response, error)
	FetchPodcastRegions(args map[string]string) (*Response, error)
	FetchPodcastLanguages(args map[string]string) (*Response, error)
	JustListen(args map[string]string) (*Response, error)
	FetchCuratedPodcastsLists(args map[string]string) (*Response, error)
	FetchRecommendationsForPodcast(id string, args map[string]string) (*Response, error)
	FetchRecommendationsForEpisode(id string, args map[string]string) (*Response, error)
	FetchMyPlaylists(args map[string]string) (*Response, error)
	FetchPlaylistByID(id string, args map[string]string) (*Response, error)
	SubmitPodcast(args map[string]string) (*Response, error)
	DeletePodcast(id string, args map[string]string) (*Response, error)
	FetchAudienceForPodcast(id string, args map[string]string) (*Response, error)
	FetchPodcastsByDomain(domainName string, args map[string]string) (*Response, error)	
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

func (c *standardHTTPClient) SearchEpisodeTitles(args map[string]string) (*Response, error) {
	return c.get("search_episode_titles", args)
}

func (c *standardHTTPClient) Typeahead(args map[string]string) (*Response, error) {
	return c.get("typeahead", args)
}

func (c *standardHTTPClient) SpellCheck(args map[string]string) (*Response, error) {
	return c.get("spellcheck", args)
}

func (c *standardHTTPClient) FetchRelatedSearches(args map[string]string) (*Response, error) {
	return c.get("related_searches", args)
}

func (c *standardHTTPClient) FetchTrendingSearches(args map[string]string) (*Response, error) {
	return c.get("trending_searches", args)
}

func (c *standardHTTPClient) FetchBestPodcasts(args map[string]string) (*Response, error) {
	return c.get("best_podcasts", args)
}

func (c *standardHTTPClient) FetchPodcastByID(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("podcasts/%s", id), args)
}

func (c *standardHTTPClient) FetchEpisodeByID(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("episodes/%s", id), args)
}

func (c *standardHTTPClient) BatchFetchEpisodes(args map[string]string) (*Response, error) {
	values := url.Values{}
	for k, v := range args {
		values.Set(k, v)
	}
	return c.post("episodes", args, values)		
}

func (c *standardHTTPClient) BatchFetchPodcasts(args map[string]string) (*Response, error) {
	values := url.Values{}
	for k, v := range args {
		values.Set(k, v)
	}
	return c.post("podcasts", args, values)	
}

func (c *standardHTTPClient) FetchCuratedPodcastsListByID(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("curated_podcasts/%s", id), args)
}

func (c *standardHTTPClient) FetchPodcastGenres(args map[string]string) (*Response, error) {
	return c.get("genres", args)
}

func (c *standardHTTPClient) FetchPodcastRegions(args map[string]string) (*Response, error) {
	return c.get("regions", args)
}

func (c *standardHTTPClient) FetchPodcastLanguages(args map[string]string) (*Response, error) {
	return c.get("languages", args)
}

func (c *standardHTTPClient) JustListen(args map[string]string) (*Response, error) {
	return c.get("just_listen", args)
}

func (c *standardHTTPClient) FetchCuratedPodcastsLists(args map[string]string) (*Response, error) {
	return c.get("curated_podcasts", args)
}

func (c *standardHTTPClient) FetchRecommendationsForPodcast(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("podcasts/%s/recommendations", id), args)
}

func (c *standardHTTPClient) FetchRecommendationsForEpisode(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("episodes/%s/recommendations", id), args)
}

func (c *standardHTTPClient) FetchMyPlaylists(args map[string]string) (*Response, error) {
	return c.get("playlists", args)
}

func (c *standardHTTPClient) FetchPlaylistByID(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("playlists/%s", id), args)
}

func (c *standardHTTPClient) SubmitPodcast(args map[string]string) (*Response, error) {
	values := url.Values{}
	for k, v := range args {
		values.Set(k, v)
	}
	return c.post("podcasts/submit", args, values)
}

func (c *standardHTTPClient) DeletePodcast(id string, args map[string]string) (*Response, error) {
	return c.delete(fmt.Sprintf("podcasts/%s", id), args)
}

func (c *standardHTTPClient) FetchAudienceForPodcast(id string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("podcasts/%s/audience", id), args)
}

func (c *standardHTTPClient) FetchPodcastsByDomain(domainName string, args map[string]string) (*Response, error) {
	return c.get(fmt.Sprintf("podcasts/domains/%s", domainName), args)
}