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

type integrationExpectation struct {
	Method       string
	Path         string
	Args         map[string]string
	ExecuteFunc  func(args map[string]string) (*listennotes.Response, error)
	ValidateFunc func(resp *listennotes.Response, respErr error) error
}

func expectCollection(resp *listennotes.Response, collectionKey string, expectedCount int) {
	if resp == nil {
		panic(fmt.Errorf("no response"))
	}

	results, ok := resp.Data[collectionKey].([]interface{})
	if !ok {
		panic(fmt.Errorf("collection Key %s not found", collectionKey))
	}

	if len(results) != expectedCount {
		panic(fmt.Errorf("expected %d results, but got %d", expectedCount, len(results)))
	}
}

func expectMap(resp *listennotes.Response, mapKey string, expectedCount int) {
	if resp == nil {
		panic(fmt.Errorf("no response"))
	}

	results, ok := resp.Data[mapKey].(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("map Key %s not found", mapKey))
	}

	if len(results) != expectedCount {
		panic(fmt.Errorf("expected %d results, but got %d", expectedCount, len(results)))
	}
}

func expectData(resp *listennotes.Response, key string, expectedValue interface{}) {
	actual, ok := resp.Data[key]
	if !ok {
		panic(fmt.Errorf("no data at %s", key))
	}
	if actual != expectedValue {
		panic(fmt.Errorf("expected (%T)%v got (%T)%v", expectedValue, expectedValue, actual, actual))
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
