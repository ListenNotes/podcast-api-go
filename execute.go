package listennotes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Response contains parsed data, usage statistics, and the original HTTP details.
// A non-2xx response is returned alongside its APIError; Data is not parsed then.
type Response struct {
	Stats      ResponseStatistics
	Data       map[string]interface{}
	StatusCode int
	Headers    http.Header
	Body       string
}

// ToJSON encodes Data, returning an empty string for nil receivers or marshal errors.
func (r *Response) ToJSON() string {
	if r == nil {
		return ""
	}
	result, err := json.Marshal(r.Data)
	if err != nil {
		return ""
	}
	return string(result)
}

type pathParam struct{ name, value string }

func (c *standardHTTPClient) requestAPI(method, path string, pathParams []pathParam, queryNames []string, args map[string]string) (*Response, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid API base URL: %w", err)
	}
	if (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" {
		return nil, fmt.Errorf("API base URL must be HTTP(S) without credentials, query, or fragment")
	}
	pathNames := make([]string, 0, len(pathParams))
	for _, param := range pathParams {
		if param.value == "" || param.value == "." || param.value == ".." {
			return nil, fmt.Errorf("invalid path parameter %s", param.name)
		}
		path = strings.ReplaceAll(path, "{"+param.name+"}", url.PathEscape(param.value))
		pathNames = append(pathNames, param.name)
	}
	query, form := url.Values{}, url.Values{}
	for name, value := range args {
		if slices.Contains(pathNames, name) {
			continue
		}
		if (method == http.MethodPost || method == http.MethodPut) && !slices.Contains(queryNames, name) {
			form.Set(name, value)
		} else {
			query.Set(name, value)
		}
	}
	var body io.Reader
	if method == http.MethodPost || method == http.MethodPut {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, strings.TrimRight(base.String(), "/")+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", "podcast-api-go "+Version)
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set(RequestHeaderKeyAPI, c.apiKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	result := &Response{Stats: parseStats(resp), StatusCode: resp.StatusCode, Headers: resp.Header.Clone(), Body: string(raw)}
	if err != nil {
		return result, fmt.Errorf("failed reading API response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, &APIError{StatusCode: resp.StatusCode, Headers: resp.Header.Clone(), Body: string(raw)}
	}
	if len(strings.TrimSpace(result.Body)) == 0 {
		result.Data = map[string]interface{}{}
	} else if err := json.Unmarshal(raw, &result.Data); err != nil {
		return result, fmt.Errorf("failed parsing the response: %w", err)
	}
	return result, nil
}
