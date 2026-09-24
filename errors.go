package listennotes

import (
	"errors"
	"fmt"
	"net/http"
)

// Known errors can be matched with errors.Is.
var (
	ErrBadRequest          = errors.New("invalid API request")
	ErrUnauthorized        = errors.New("wrong API key, suspended account, or unavailable plan feature")
	ErrForbidden           = errors.New("the API account cannot perform this operation")
	ErrNotFound            = errors.New("endpoint or requested content does not exist")
	ErrTooManyRequests     = errors.New("API quota or rate limit exceeded")
	ErrInternalServerError = errors.New("API server error")
	ErrUnexpectedStatus    = errors.New("unexpected HTTP status")
)

var errMap = map[int]error{
	400: ErrBadRequest,
	401: ErrUnauthorized,
	403: ErrForbidden,
	404: ErrNotFound,
	429: ErrTooManyRequests,
}

// APIError retains the status, headers, and unmodified body of a non-2xx response.
// Use errors.As to access these fields and errors.Is to match a known error.
type APIError struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// Unwrap classifies errors without losing the server's response details.
func (e *APIError) Unwrap() error {
	if mapped := errMap[e.StatusCode]; mapped != nil {
		return mapped
	}
	if e.StatusCode >= 500 && e.StatusCode <= 599 {
		return ErrInternalServerError
	}
	return ErrUnexpectedStatus
}
