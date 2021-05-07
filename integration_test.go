package listennotes_test

import (
	"fmt"
	"testing"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func TestAuthErrorIntegration(t *testing.T) {
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

	client := listennotes.NewClient("")

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
					t.Errorf("Execution [%d] %s failed expectation: %v", idx, exp.Path, r)
				}
			}()
			if err := exp.ValidateFunc(resp, err); err != nil {
				t.Errorf("Execution [%d] %s validation failed: %s", idx, exp.Path, err)
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
