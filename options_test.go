package listennotes_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func TestWithHTTPOptions(t *testing.T) {

	noDial := fmt.Errorf("no-dial")

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return nil, noDial
			},
		},
	}

	client := listennotes.NewClient("", listennotes.WithHTTPClient(httpClient))
	_, err := client.Search(map[string]string{"q": "a"})

	if !errors.Is(err, noDial) {
		t.Errorf("mock http client failure was not as expected: %s", err)
	}

}

func TestWithBaseURL(t *testing.T) {
	noDial := fmt.Errorf("no-dial")

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				if addr != "localhost:80" {
					t.Errorf("custom baseURL was not used: %s", addr)
				}
				return nil, noDial
			},
		},
	}

	client := listennotes.NewClient("", listennotes.WithHTTPClient(httpClient), listennotes.WithBaseURL("http://localhost/test-url"))
	_, err := client.Search(map[string]string{"q": "a"})

	if !errors.Is(err, noDial) {
		t.Errorf("mock http client failure was not as expected: %s", err)
	}
}
