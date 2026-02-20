package maclookup

import (
	"fmt"
	"time"
)

// HTTPClientError wraps a low-level HTTP or network error (e.g. connection
// refused, context deadline exceeded, or an unexpected HTTP status code).
type HTTPClientError struct {
	Err error
}

// Error implements the error interface.
func (c *HTTPClientError) Error() string {
	return c.Err.Error()
}

// Unwrap returns the underlying error, enabling errors.Is / errors.As traversal.
func (c *HTTPClientError) Unwrap() error {
	return c.Err
}

// BadAPIRequest is returned when the API responds with HTTP 400, indicating
// that the supplied MAC address or query parameters are invalid.
type BadAPIRequest struct {
	Err error
}

// Error implements the error interface.
func (c *BadAPIRequest) Error() string {
	return c.Err.Error()
}

// Unwrap returns the underlying error.
func (c *BadAPIRequest) Unwrap() error {
	return c.Err
}

// BadAPIKey is returned when the API responds with HTTP 401, indicating
// that the supplied API key is missing or invalid.
type BadAPIKey struct {
	Err error
}

// Error implements the error interface.
func (c *BadAPIKey) Error() string {
	return c.Err.Error()
}

// Unwrap returns the underlying error.
func (c *BadAPIKey) Unwrap() error {
	return c.Err
}

// RateLimitsExceeded is returned when the API responds with HTTP 429.
// Limit holds the maximum request quota for the current window and Reset
// indicates when the counter will be reset.
type RateLimitsExceeded struct {
	Limit int64
	Reset time.Time
	Err   error
}

// Error implements the error interface.
func (c *RateLimitsExceeded) Error() string {
	return fmt.Sprintf("rate limits exceded. current limit is %d. next reset %s", c.Limit, c.Reset.Format(time.RFC3339))
}

// BadAPIResponse is returned when the API returns HTTP 200 but the response
// body cannot be decoded or does not indicate success.
type BadAPIResponse struct {
	Err error
}

// Error implements the error interface.
func (c *BadAPIResponse) Error() string {
	return c.Err.Error()
}

// Unwrap returns the underlying error.
func (c *BadAPIResponse) Unwrap() error {
	return c.Err
}
