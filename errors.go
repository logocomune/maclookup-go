package maclookup

import (
	"fmt"
	"time"
)

type HTTPClientError struct {
	Err error
}

func (c *HTTPClientError) Error() string {
	return c.Err.Error()
}

func (c *HTTPClientError) Unwrap() error {
	return c.Err
}

type BadAPIRequest struct {
	Err error
}

func (c *BadAPIRequest) Error() string {
	return c.Err.Error()
}

func (c *BadAPIRequest) Unwrap() error {
	return c.Err
}

type BadAPIKey struct {
	Err error
}

func (c *BadAPIKey) Error() string {
	return c.Err.Error()
}

func (c *BadAPIKey) Unwrap() error {
	return c.Err
}

type RateLimitsExceeded struct {
	Limit int64
	Reset time.Time
	Err   error
}

func (c *RateLimitsExceeded) Error() string {
	return fmt.Sprintf("rate limits exceded. current limit is %d. next reset %s", c.Limit, c.Reset.Format(time.RFC3339))
}

type BadAPIResponse struct {
	Err error
}

func (c *BadAPIResponse) Error() string {
	return c.Err.Error()
}

func (c *BadAPIResponse) Unwrap() error {
	return c.Err
}
