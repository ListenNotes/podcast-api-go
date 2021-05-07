package listennotes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Response is the standard response for all client functions
type Response struct {
	Stats ResponseStatistics
	Data  map[string]interface{}
}

// ToJSON will encode the response data as JSON.
// Note: JSON marshal errors are swallowed here on purpose.  This is for ease of use.
// Considering this data marshalled from JSON, the risk here is low.  On failure, "" will be returned.
func (r *Response) ToJSON() string {
	if r == nil {
		return ""
	}
	jsonResult, err := json.Marshal(r.Data)
	if err != nil {
		log.Printf("failed to marshal response data to json: %s", err)
		return ""
	}
	return string(jsonResult)
}

func (c *standardHTTPClient) get(path string, args map[string]string) (*Response, error) {
	return c.exec("GET", path, args, url.Values{})
}

func (c *standardHTTPClient) post(path string, args map[string]string, formFields url.Values) (*Response, error) {
	return c.exec("POST", path, args, formFields)
}

func (c *standardHTTPClient) delete(path string, args map[string]string) (*Response, error) {
	return c.exec("DELETE", path, args, url.Values{})
}

func (c *standardHTTPClient) exec(
	method string,
	path string,
	args map[string]string,
	formFields url.Values,
) (*Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	var body io.Reader
	if len(formFields) > 0 {
		body = strings.NewReader(formFields.Encode())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", path, err)
	}
	req.Header.Add(RequestHeaderKeyAPI, c.apiKey)

	if len(formFields) > 0 {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	q := req.URL.Query()
	for k, v := range args {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to executing request to %s: %w", path, err)
	}
	defer resp.Body.Close()

	// map any generic status code errors
	if mappedError, ok := errMap[resp.StatusCode]; ok && mappedError != nil {
		return nil, mappedError
	}

	// generic body parsing
	var genericJSON map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&genericJSON); err != nil {
		return nil, fmt.Errorf("failed parsing the response from %s: %w", path, err)
	}

	// gather the header statistics
	stats := parseStats(resp)

	return &Response{
		Stats: stats,
		Data:  genericJSON,
	}, nil
}
