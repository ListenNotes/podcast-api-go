package listennotes_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	listennotes "github.com/ListenNotes/podcast-api-go"
)

func TestMockURL(t *testing.T) {
	noDial := fmt.Errorf("no-dial")

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				if addr != "listen-api-test.listennotes.com:443" {
					t.Errorf("custom baseURL was not used: %s", addr)
				}
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

func TestProductionURL(t *testing.T) {
	noDial := fmt.Errorf("no-dial")

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				if addr != "listen-api.listennotes.com:443" {
					t.Errorf("custom baseURL was not used: %s", addr)
				}
				return nil, noDial
			},
		},
	}

	client := listennotes.NewClient("anapikey", listennotes.WithHTTPClient(httpClient))
	_, err := client.Search(map[string]string{"q": "a"})

	if !errors.Is(err, noDial) {
		t.Errorf("mock http client failure was not as expected: %s", err)
	}
}
