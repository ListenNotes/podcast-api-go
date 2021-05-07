package listennotes_test

import (
	"fmt"
	"net/http"
	"testing"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func TestSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	rt := &integrationTestRoundTripper{}
	httpClient := &http.Client{
		Transport: rt,
	}

	client := listennotes.NewClient("", listennotes.WithHTTPClient(httpClient))

	var expectations = []integrationExpectation{
		{
			Method: "GET",
			Path:   "search",
			Args: map[string]string{
				"q": "term",
			},
			ExecuteFunc: client.Search,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "took", 0.234)
				expectCollection(resp, "results", 10)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "typeahead",
			Args: map[string]string{
				"q": "term",
			},
			ExecuteFunc: client.Typeahead,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "terms", 2)
				expectCollection(resp, "genres", 1)
				expectCollection(resp, "podcasts", 5)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "best_podcasts",
			Args: map[string]string{
				"genre_id": "93",
			},
			ExecuteFunc: client.FetchBestPodcasts,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "name", "Business")
				expectCollection(resp, "podcasts", 20)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "podcasts/4d3fe717742d4963a85562e9f84d8c79",
			Args: map[string]string{
				"sort": "recent_first",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchPodcastByID("4d3fe717742d4963a85562e9f84d8c79", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "id", "4d3fe717742d4963a85562e9f84d8c79")
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "episodes/6b6d65930c5a4f71b254465871fed370",
			Args: map[string]string{
				"show_transcript": "1",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchEpisodeByID("6b6d65930c5a4f71b254465871fed370", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "id", "6b6d65930c5a4f71b254465871fed370")
				return nil
			},
		},
		{
			Method: "POST",
			Path:   "episodes",
			Args:   map[string]string{},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.BatchFetchEpisodes([]string{"123", "456"}, args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "episodes", 2)
				return nil
			},
		},
		{
			Method: "POST",
			Path:   "podcasts",
			Args: map[string]string{
				"show_latest_episodes": "1",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.BatchFetchPodcasts([]string{"123", "456"}, args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "podcasts", 9)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "curated_podcasts/SDFKduyJ47r",
			Args:   map[string]string{},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchCuratedPodcastsListByID("SDFKduyJ47r", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "title", "16 Brilliant Indian Podcasts That'll Make You A Funner, Smarter, Better Informed Person")
				expectCollection(resp, "podcasts", 16)
				return nil
			},
		},

		{
			Method: "GET",
			Path:   "genres",
			Args: map[string]string{
				"top_level_only": "1",
			},
			ExecuteFunc: client.FetchPodcastGenres,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectNoError(respErr)
				expectCollection(resp, "genres", 8)
				return nil
			},
		},
		{
			Method:      "GET",
			Path:        "regions",
			Args:        map[string]string{},
			ExecuteFunc: client.FetchPodcastRegions,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectMap(resp, "regions", 67)
				return nil
			},
		},
		{
			Method:      "GET",
			Path:        "languages",
			Args:        map[string]string{},
			ExecuteFunc: client.FetchPodcastLanguages,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "languages", 20)
				return nil
			},
		},
		{
			Method:      "GET",
			Path:        "just_listen",
			Args:        map[string]string{},
			ExecuteFunc: client.JustListen,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "title", "Miami Heat: Howard Beck on James Harden (Nets, 76ers too)")
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "curated_podcasts",
			Args: map[string]string{
				"page": "2",
			},
			ExecuteFunc: client.FetchCuratedPodcastsLists,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "has_previous", true)
				expectCollection(resp, "curated_lists", 20)
				return nil
			},
		},

		{
			Method: "GET",
			Path:   "podcasts/25212ac3c53240a880dd5032e547047b/recommendations",
			Args: map[string]string{
				"safe_mode": "0",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchRecommendationsForPodcast("25212ac3c53240a880dd5032e547047b", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "recommendations", 8)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "episodes/254444fa6cf64a43a95292a70eb6869b/recommendations",
			Args: map[string]string{
				"page": "2",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchRecommendationsForEpisode("254444fa6cf64a43a95292a70eb6869b", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "recommendations", 8)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "playlists",
			Args: map[string]string{
				"sort": "recent_added_first",
			},
			ExecuteFunc: client.FetchMyPlaylists,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "has_previous", false)
				expectCollection(resp, "playlists", 3)
				return nil
			},
		},
		{
			Method: "GET",
			Path:   "playlists/m1pe7z60bsw",
			Args: map[string]string{
				"type": "episode_list",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchPlaylistByID("m1pe7z60bsw", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "type", "episode_list")
				expectCollection(resp, "items", 20)
				return nil
			},
		},
		{
			Method: "POST",
			Path:   "podcasts/submit",
			Args: map[string]string{
				"rss":   "https://feeds.megaphone.fm/committed",
				"email": "hello@example.com",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.SubmitPodcast(args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "status", "found")
				return nil
			},
		},
		{
			Method: "DELETE",
			Path:   "podcasts/abc",
			Args: map[string]string{
				"reason": "the reason",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.DeletePodcast("abc", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "status", "in review")
				return nil
			},
		},
	}

	for idx, exp := range expectations {
		resp, err := exp.ExecuteFunc(exp.Args)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("[%d] %s failed expectation: %v", idx, exp.Path, r)
				}
			}()
			if err := exp.ValidateFunc(resp, err); err != nil {
				t.Errorf("[%d] %s validation failed: %s", idx, exp.Path, err)
			}
			expectedPath := fmt.Sprintf("/api/v2/%s", exp.Path)
			if rt.LastRequest.URL.Path != expectedPath {
				t.Errorf("[%d] %s had unexpected Path: %s", idx, exp.Path, rt.LastRequest.URL.Path)
			}
			if rt.LastRequest.Method != exp.Method {
				t.Errorf("[%d] %s had unexpected Method: %s", idx, exp.Path, rt.LastRequest.Method)
			}
			for qk, qv := range exp.Args {
				if rt.LastRequest.URL.Query().Get(qk) != qv {
					t.Errorf("[%d] %s did not have expected query arg: %s=%s in %s", idx, exp.Path, qk, qv, rt.LastRequest.URL.RawQuery)
				}
			}
		}()
	}
}
