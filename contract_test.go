package listennotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type sdkOperation struct {
	OperationID string `json:"operationId"`
	Func        string `json:"func"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Parameters  []struct {
		Name string `json:"name"`
		In   string `json:"in"`
	} `json:"parameters"`
	ExampleParams map[string]interface{} `json:"example_params"`
}

func loadSDKContract(t *testing.T) []sdkOperation {
	t.Helper()
	data, err := os.ReadFile("api-contract.json")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Version    string         `json:"version"`
		Operations []sdkOperation `json:"operations"`
	}
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatal(err)
	}
	if contract.Version != Version || len(contract.Operations) != 31 {
		t.Fatalf("unexpected contract version/method count: %s/%d", contract.Version, len(contract.Operations))
	}
	return contract.Operations
}

func exampleArgs(op sdkOperation) map[string]string {
	args := map[string]string{}
	for name, value := range op.ExampleParams {
		if text, ok := value.(string); ok {
			args[name] = text
		} else if value != nil {
			raw, _ := json.Marshal(value)
			args[name] = string(raw)
		}
	}
	return args
}

func TestEveryGeneratedMethodMatchesContract(t *testing.T) {
	clientType := reflect.TypeOf((*HTTPClient)(nil)).Elem()
	ops := loadSDKContract(t)
	if clientType.NumMethod() != len(ops) {
		t.Fatal("public interface and contract differ")
	}
	for _, op := range ops {
		t.Run(op.OperationID, func(t *testing.T) {
			if _, ok := clientType.MethodByName(op.Func); !ok {
				t.Fatalf("missing public method %s", op.Func)
			}
			args := exampleArgs(op)
			before := maps.Clone(args)
			path, query, form := op.Path, url.Values{}, url.Values{}
			for _, param := range op.Parameters {
				value, exists := args[param.Name]
				if !exists {
					continue
				}
				switch param.In {
				case "path":
					path = strings.ReplaceAll(path, "{"+param.Name+"}", url.PathEscape(value))
				case "query":
					query.Set(param.Name, value)
				case "body":
					form.Set(param.Name, value)
				}
			}
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method != op.Method || r.URL.EscapedPath() != "/api/v2"+path || r.URL.RawQuery != query.Encode() {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				raw, _ := io.ReadAll(r.Body)
				if string(raw) != form.Encode() {
					t.Errorf("body = %q, want %q", raw, form.Encode())
				}
				if r.Header.Get(RequestHeaderKeyAPI) != "fixture-key" || r.Header.Get("User-Agent") != "podcast-api-go "+Version {
					t.Error("incorrect API key or User-Agent")
				}
				if (op.Method == "POST" || op.Method == "PUT") && r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
					t.Error("missing form content type")
				}
				w.Header().Set(ResponseHeaderKeyFreeQuota, "25000")
				w.Header().Set(ResponseHeaderKeyUsage, "19231")
				w.Header().Set(ResponseHeaderKeyLatencySeconds, "0.056")
				w.Header().Set(ResponseHeaderKeyNextBillingDate, "2020-09-26T17:27:33.110641+00:00")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `{"id":23,"ok":true}`)
			}))
			defer server.Close()
			response, err := callSDKMethod(NewClient("fixture-key", WithBaseURL(server.URL+"/api/v2/")), op.OperationID, args)
			if err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 1 || !maps.Equal(args, before) {
				t.Fatal("method retried or mutated caller parameters")
			}
			if response.StatusCode != 201 || response.Data["ok"] != true || response.Stats.FreeQuota != 25000 || response.Stats.Usage != 19231 || response.Stats.LatencySeconds != 0.056 || response.Stats.NextBillingDate.Year() != 2020 {
				t.Fatalf("lost response details: %+v", response)
			}
		})
	}
}

func TestNestedPathsAndEmptyValues(t *testing.T) {
	var actualPath, actualQuery, actualBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualPath, actualQuery = r.URL.EscapedPath(), r.URL.RawQuery
		raw, _ := io.ReadAll(r.Body)
		actualBody = string(raw)
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	client := NewClient("", WithBaseURL(server.URL+"/api/v2"))
	id, itemID := "abc/def?# café%", "23/?#"
	if _, err := client.UpdatePlaylistItemNotes(id, itemID, map[string]string{"notes": "", "id": "ignored", "item_id": "ignored"}); err != nil {
		t.Fatal(err)
	}
	if actualPath != "/api/v2/playlists/"+url.PathEscape(id)+"/items/"+url.PathEscape(itemID) || actualQuery != "" || actualBody != "notes=" {
		t.Fatalf("incorrect nested request: %s?%s body %q", actualPath, actualQuery, actualBody)
	}
	if _, err := client.UpdatePlaylist("abc", map[string]string{"description": ""}); err != nil || actualBody != "description=" {
		t.Fatalf("empty description lost: %s, %v", actualBody, err)
	}
	if _, err := client.Search(map[string]string{"q": "a & café +/#?", "offset": "0", "safe_mode": "false"}); err != nil {
		t.Fatal(err)
	}
	query, _ := url.ParseQuery(actualQuery)
	if query.Get("q") != "a & café +/#?" || query.Get("offset") != "0" || query.Get("safe_mode") != "false" || actualBody != "" {
		t.Fatal("query encoding changed values")
	}
	if _, err := client.DeletePodcast(id, map[string]string{"reason": "a & b"}); err != nil || actualQuery != "reason=a+%26+b" || actualBody != "" {
		t.Fatalf("DELETE query incorrect: %s, %v", actualQuery, err)
	}
}

func TestDeletePlaylistEncodesIDWithoutQueryOrBody(t *testing.T) {
	for name, args := range map[string]map[string]string{"nil": nil, "empty": {}, "path field": {"id": "must-not-leak"}} {
		t.Run(name, func(t *testing.T) {
			id := "abc/def?# café%"
			before := maps.Clone(args)
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method != http.MethodDelete || r.URL.EscapedPath() != "/api/v2/playlists/"+url.PathEscape(id) || r.URL.RawQuery != "" {
					t.Errorf("incorrect deletion request: %s %s", r.Method, r.URL)
				}
				body, _ := io.ReadAll(r.Body)
				if len(body) != 0 || r.Header.Get("Content-Type") != "" {
					t.Error("playlist deletion must not send a request body or content type")
				}
				w.Header().Set(ResponseHeaderKeyUsage, "12")
				json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "deleted": true})
			}))
			defer server.Close()
			response, err := NewClient("fixture-key", WithBaseURL(server.URL+"/api/v2")).DeletePlaylist(id, args)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != 200 || response.Data["id"] != id || response.Data["deleted"] != true || response.Stats.Usage != 12 {
				t.Fatalf("lost deletion response details: %+v", response)
			}
			if requests.Load() != 1 || !maps.Equal(args, before) {
				t.Fatal("deletion retried or changed the caller's parameters")
			}
		})
	}
}

func TestWriteQueryAndBodyFieldsAreSeparated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if r.URL.RawQuery != "page=0" || string(raw) != "notes=" {
			t.Errorf("incorrect query/body: %s/%s", r.URL.RawQuery, raw)
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	client := NewClient("", WithBaseURL(server.URL)).(*standardHTTPClient)
	_, err := client.requestAPI("PUT", "/fixture", nil, []string{"page"}, map[string]string{"page": "0", "notes": ""})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHTTPFailuresRetainDetailsAndNeverRetry(t *testing.T) {
	for _, code := range []int{301, 302, 307, 308, 400, 401, 403, 404, 405, 418, 429, 500, 502, 503} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				w.Header().Set("X-Fixture", "retained")
				w.Header().Set(ResponseHeaderKeyUsage, "10")
				w.WriteHeader(code)
				fmt.Fprint(w, "exact server explanation")
			}))
			defer server.Close()
			client := NewClient("", WithBaseURL(server.URL))
			for index, operation := range []string{"createPlaylist", "deletePlaylist"} {
				args := map[string]string{"name": "fixture"}
				if operation == "deletePlaylist" {
					args = map[string]string{"id": "fixture"}
				}
				response, err := callSDKMethod(client, operation, args)
				var apiError *APIError
				if !errors.As(err, &apiError) || apiError.StatusCode != code || apiError.Body != "exact server explanation" || apiError.Headers.Get("X-Fixture") != "retained" {
					t.Fatalf("missing HTTP error context: %v", err)
				}
				if response == nil || response.Stats.Usage != 10 || response.Body != apiError.Body || requests.Load() != int32(index+1) {
					t.Fatal("lost response or retried write")
				}
				expected := errMap[code]
				if expected == nil {
					expected = ErrUnexpectedStatus
					if code >= 500 {
						expected = ErrInternalServerError
					}
				}
				if !errors.Is(err, expected) {
					t.Fatalf("wrong classification: %v", err)
				}
			}
		})
	}
}

func TestAllSuccessStatusesAndEmptyResponses(t *testing.T) {
	for _, code := range []int{200, 201, 202, 204, 299} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(code) }))
		response, err := NewClient("", WithBaseURL(server.URL)).DeletePlaylistItem("playlist", "23", nil)
		server.Close()
		if err != nil || response.StatusCode != code || response.ToJSON() != "{}" {
			t.Fatalf("success status %d failed: %+v/%v", code, response, err)
		}
	}
}

func TestDefaultRedirectsAndClientIsolation(t *testing.T) {
	var leaked atomic.Int32
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1) }))
	defer destination.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, destination.URL, http.StatusTemporaryRedirect)
	}))
	defer redirect.Close()
	if _, err := NewClient("fixture-key", WithBaseURL(redirect.URL)).Search(nil); !errors.Is(err, ErrUnexpectedStatus) || leaked.Load() != 0 {
		t.Fatalf("redirect followed: %v", err)
	}
	if _, err := NewClient("fixture-key", WithBaseURL(redirect.URL)).DeletePlaylist("list", nil); !errors.Is(err, ErrUnexpectedStatus) || leaked.Load() != 0 {
		t.Fatalf("deletion redirect followed: %v", err)
	}
	var seen sync.Map
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.Store(r.URL.Query().Get("q"), r.Header.Get(RequestHeaderKeyAPI))
		if r.Method == http.MethodDelete {
			seen.Store(r.URL.Path, r.Header.Get(RequestHeaderKeyAPI))
		}
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	first := NewClient("first", WithBaseURL(server.URL)).(*standardHTTPClient)
	second := NewClient("second", WithBaseURL(server.URL)).(*standardHTTPClient)
	first.httpClient.Timeout = 15 * time.Second
	if second.httpClient.Timeout != 30*time.Second {
		t.Fatal("default HTTP client config shared")
	}
	var wg sync.WaitGroup
	for i, client := range []HTTPClient{first, second, first, second} {
		wg.Go(func() {
			if _, err := client.Search(map[string]string{"q": fmt.Sprint(i % 2)}); err != nil {
				t.Error(err)
			}
			if _, err := client.DeletePlaylist(fmt.Sprint(i%2), nil); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if key, _ := seen.Load("0"); key != "first" {
		t.Error("first key was overwritten")
	}
	if key, _ := seen.Load("1"); key != "second" {
		t.Error("second key was overwritten")
	}
	if key, _ := seen.Load("/playlists/0"); key != "first" {
		t.Error("first deletion key was overwritten")
	}
	if key, _ := seen.Load("/playlists/1"); key != "second" {
		t.Error("second deletion key was overwritten")
	}
}

func TestTimeoutAndInvalidParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	client := NewClient("", WithBaseURL(server.URL), WithHTTPClient(&http.Client{Timeout: 20 * time.Millisecond}))
	if _, err := client.Search(nil); err == nil {
		t.Fatal("custom timeout ignored")
	} else {
		var urlError *url.Error
		if !errors.As(err, &urlError) || !urlError.Timeout() {
			t.Fatalf("timeout cause lost: %v", err)
		}
	}
	var requests atomic.Int32
	fast := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, `{}`)
	}))
	defer fast.Close()
	client = NewClient("", WithBaseURL(fast.URL), WithHTTPClient(nil), nil)
	for _, id := range []string{"", ".", ".."} {
		if _, err := client.DeletePlaylist(id, nil); err == nil {
			t.Errorf("accepted invalid playlist deletion identifier %q", id)
		}
		if _, err := client.DeletePlaylistItem(id, "23", nil); err == nil {
			t.Errorf("accepted invalid playlist identifier %q", id)
		}
		if _, err := client.UpdatePlaylistItemNotes("valid", id, nil); err == nil {
			t.Errorf("accepted invalid item identifier %q", id)
		}
	}
	for _, base := range []string{"file:///tmp/fixture", fast.URL + "?query=yes", fast.URL + "#fragment", "http://user:pass@localhost"} {
		if _, err := NewClient("", WithBaseURL(base)).Search(nil); err == nil {
			t.Errorf("accepted invalid base URL %q", base)
		}
	}
	if requests.Load() != 0 {
		t.Fatal("invalid parameters caused a request")
	}
	if _, err := client.Search(nil); err != nil {
		t.Fatalf("nil options should preserve defaults: %v", err)
	}
}
