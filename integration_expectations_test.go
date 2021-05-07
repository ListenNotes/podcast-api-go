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
			Path:   "podcasts/123",
			Args: map[string]string{
				"sort": "recent_first",
			},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.FetchPodcastByID("123", args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectData(resp, "id", "4d3fe717742d4963a85562e9f84d8c79")
				return nil
			},
		},
		{
			Method: "POST",
			Path:   "episodes",
			Args:   map[string]string{},
			ExecuteFunc: func(args map[string]string) (*listennotes.Response, error) {
				return client.BatchFetchPodcastByID([]string{"123", "456"}, args)
			},
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectNoError(respErr)
				expectCollection(resp, "episodes", 2)
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
