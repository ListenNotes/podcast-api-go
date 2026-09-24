//go:build integration

package listennotes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type mockOnlyTransport struct {
	transport *http.Transport
	last      *http.Request
	body      string
}

func (m *mockOnlyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != "listen-api-test.listennotes.com" || r.URL.User != nil || !strings.HasPrefix(r.URL.Path, "/api/v2/") || len(r.Header.Values(RequestHeaderKeyAPI)) != 0 {
		return nil, fmt.Errorf("integration request rejected: destination or credentials are not allowed")
	}
	m.last, m.body = r.Clone(r.Context()), ""
	if r.GetBody != nil {
		body, err := r.GetBody()
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(body)
		body.Close()
		if err != nil {
			return nil, err
		}
		m.body = string(raw)
	}
	return m.transport.RoundTrip(r)
}

func mockIntegrationClient(t *testing.T) (HTTPClient, *mockOnlyTransport) {
	t.Helper()
	httpClient := newHTTPClient()
	transport := httpClient.Transport.(*http.Transport)
	transport.Proxy = nil
	t.Cleanup(transport.CloseIdleConnections)
	guard := &mockOnlyTransport{transport: transport}
	httpClient.Transport = guard
	return NewClient("", WithHTTPClient(httpClient)), guard
}

func TestMockIntegrationAllMethods(t *testing.T) {
	client, guard := mockIntegrationClient(t)
	for _, op := range loadSDKContract(t) {
		t.Run(op.OperationID, func(t *testing.T) {
			response, err := callSDKMethod(client, op.OperationID, exampleArgs(op))
			if err != nil {
				t.Fatal(err)
			}
			status := 200
			if op.OperationID == "createPlaylist" || op.OperationID == "addPlaylistItem" {
				status = 201
			}
			if guard.last.Method != op.Method || response.StatusCode != status || len(response.Data) == 0 {
				t.Fatalf("unexpected response: %+v", response)
			}
			for _, header := range []string{ResponseHeaderKeyUsage, ResponseHeaderKeyFreeQuota, ResponseHeaderKeyLatencySeconds, ResponseHeaderKeyNextBillingDate} {
				if response.Headers.Get(header) == "" {
					t.Errorf("missing response header %s", header)
				}
			}
			switch op.OperationID {
			case "createPlaylist", "updatePlaylist", "getPlaylistById":
				if _, ok := response.Data["id"].(string); !ok {
					t.Error("missing playlist ID")
				}
				if kind := response.Data["type"]; kind != "episode_list" && kind != "podcast_list" {
					t.Error("incorrect playlist type")
				}
			case "addPlaylistItem", "updatePlaylistItemNotes":
				if _, ok := response.Data["id"].(float64); !ok {
					t.Error("missing item ID")
				}
				if _, ok := response.Data["notes"].(string); !ok {
					t.Error("missing notes")
				}
			case "deletePlaylistItem":
				if response.Data["deleted"] != true {
					t.Error("missing deletion confirmation")
				}
			}
		})
		time.Sleep(100 * time.Millisecond)
	}
}

func TestMockIntegrationEncodedAndEmptyFields(t *testing.T) {
	client, guard := mockIntegrationClient(t)
	if _, err := client.Search(map[string]string{"q": "c++ & café / podcasts"}); err != nil {
		t.Fatal(err)
	}
	if guard.last.URL.Query().Get("q") != "c++ & café / podcasts" {
		t.Fatal("incorrect search encoding")
	}
	if _, err := client.UpdatePlaylist("m1pe7z60bsw", map[string]string{"description": ""}); err != nil {
		t.Fatal(err)
	}
	if guard.body != "description=" || guard.last.URL.RawQuery != "" {
		t.Fatal("empty description was omitted or sent in the query string")
	}
	if _, err := client.AddPlaylistItem("m1pe7z60bsw", map[string]string{"podcast_id": "4d3fe717742d4963a85562e9f84d8c79", "notes": ""}); err != nil {
		t.Fatal(err)
	}
	form, _ := url.ParseQuery(guard.body)
	if !form.Has("notes") || form.Has("episode_id") || form.Get("podcast_id") == "" {
		t.Fatal("incorrect add-item form")
	}
}

func TestMockIntegrationMissingRoute(t *testing.T) {
	client, _ := mockIntegrationClient(t)
	client.(*standardHTTPClient).baseURL = BaseURLTest + "/missing-route"
	response, err := client.Search(nil)
	if !errors.Is(err, ErrNotFound) || response == nil || response.StatusCode != 404 {
		t.Fatalf("missing route error lost: %v", err)
	}
}
