package listennotes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStandardClientExecuteNewReqFailure(t *testing.T) {
	client := &standardHTTPClient{
		baseURL: "http://localhost:bogus",
	}
	_, err := client.execute(map[string]string{}, "path")
	if err == nil || !strings.Contains(err.Error(), "invalid port ") {
		t.Errorf("Expected url parse failure but got: %v", err)
	}
}

func TestMappedErrors(t *testing.T) {
	type toTest struct {
		code int
		err  error
	}

	// expected code to error mappings -- dups the errMap however this means that the test is actually validating the map
	errs := []toTest{
		{code: 200, err: nil},
		{code: 400, err: ErrBadRequest},
		{code: 401, err: ErrUnauthorized},
		{code: 404, err: ErrNotFound},
		{code: 429, err: ErrTooManyRequests},
		{code: 500, err: ErrInternalServerError},
	}

	for k, v := range errMap {
		errs = append(errs, toTest{code: k, err: v})
	}

	var expectedCode int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedCode)
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	client := &standardHTTPClient{
		httpClient: http.DefaultClient,
		baseURL:    ts.URL,
	}

	for _, e := range errs {
		expectedCode = e.code
		_, err := client.execute(map[string]string{}, "path")
		if (e.err == nil && err != nil) || (e.err != nil && !errors.Is(err, e.err)) {
			t.Errorf("%d reponse code did not result in correct error: %s", e.code, err)
		}
	}
}

func TestDecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not-json"))
	}))
	defer ts.Close()

	client := &standardHTTPClient{
		httpClient: http.DefaultClient,
		baseURL:    ts.URL,
	}
	_, err := client.execute(map[string]string{}, "path")
	if err == nil || !strings.Contains(err.Error(), "failed parsing the response") {
		t.Errorf("Expected json parse failure but got: %v", err)
	}
}

func TestGetQueryArguments(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.URL.RawQuery != "a=b&c=d" {
			t.Errorf("Query parameters were not as expected: %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	client := &standardHTTPClient{
		httpClient: http.DefaultClient,
		baseURL:    ts.URL,
	}
	client.execute(map[string]string{
		"a": "b",
		"c": "d",
	}, "path")

	if !called {
		t.Errorf("Did not call expected httptest url")
	}
}

func TestParsedResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(searchPayload))
	}))
	defer ts.Close()

	client := &standardHTTPClient{
		httpClient: http.DefaultClient,
		baseURL:    ts.URL,
	}
	resp, err := client.execute(map[string]string{
		"a": "b",
		"c": "d",
	}, "path")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if resp == nil {
		t.Fatalf("Expected resp but go nil")
		return
	}

	if v := resp.Data["took"]; v != float64(0.693) {
		t.Errorf("Wrong took value: %v", v)
	}

	resultZero := (resp.Data["results"].([]interface{}))[0].(map[string]interface{})
	if v := resultZero["id"]; v != "ea09b575d07341599d8d5b71f205517b" {
		t.Errorf("Wrong results[0].id value: %v", v)
	}
}

const searchPayload = `{
	"took": 0.693,
	"count": 10,
	"total": 9499,
	"results": [
	  {
		"id": "ea09b575d07341599d8d5b71f205517b"
	  }
	],
	"next_offset": 10
}`
