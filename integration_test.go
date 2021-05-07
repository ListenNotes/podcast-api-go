package listennotes_test

import (
	"fmt"
	"net/http"
	"testing"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func TestAuthErrorIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := listennotes.NewClient("not-valid")
	_, err := client.Search(nil)
	if err != listennotes.ErrUnauthorized {
		t.Errorf("Expected bad token to result in unauthorized")
	}
}

func TestSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	rt := &integrationTestRoundTripper{}
	httpClient := &http.Client{
		Transport: rt,
	}

	client := listennotes.NewClient("", listennotes.WithHTTPClient(httpClient))

	expectations := []integrationExpectation{
		{
			Method: "GET",
			Path:   "search",
			Args: map[string]string{
				"q": "term",
			},
			ExecuteFunc: client.Search,
			ValidateFunc: func(resp *listennotes.Response, respErr error) error {
				expectData(resp, "took", 0.234)
				expectResults(resp, 10)
				expectNoError(respErr)
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

type integrationExpectation struct {
	Method       string
	Path         string
	Args         map[string]string
	ExecuteFunc  func(args map[string]string) (*listennotes.Response, error)
	ValidateFunc func(resp *listennotes.Response, respErr error) error
}

func expectResults(resp *listennotes.Response, expectedCount int) {
	if resp == nil {
		panic(fmt.Errorf("no response"))
	}

	results, ok := resp.Data["results"].([]interface{})
	if !ok {
		panic(fmt.Errorf("no results"))
	}

	if len(results) != expectedCount {
		panic(fmt.Errorf("expected %d results, but got %d", expectedCount, len(results)))
	}
}

func expectData(resp *listennotes.Response, key string, expectedValue interface{}) {
	actual, ok := resp.Data["took"]
	if !ok {
		panic(fmt.Errorf("no data at %s", key))
	}
	if actual != expectedValue {
		panic(fmt.Errorf("expected %v got %v", expectedValue, actual))
	}
}

func expectNoError(err error) {
	if err != nil {
		panic(fmt.Errorf("Expected no error but got: %s", err))
	}
}

type integrationTestRoundTripper struct {
	LastRequest *http.Request
}

func (rt *integrationTestRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	rt.LastRequest = req
	return http.DefaultTransport.RoundTrip(req)
}
