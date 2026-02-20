package maclookup

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- HTTPClientError ---

func TestHTTPClientError_Error(t *testing.T) {
	inner := errors.New("connection refused")
	e := &HTTPClientError{Err: inner}
	assert.Equal(t, "connection refused", e.Error())
}

func TestHTTPClientError_Unwrap(t *testing.T) {
	inner := errors.New("connection refused")
	e := &HTTPClientError{Err: inner}
	assert.True(t, errors.Is(e, inner))
}

// --- BadAPIRequest ---

func TestBadAPIRequest_Error(t *testing.T) {
	inner := errors.New("mac must be greater than 5 chars")
	e := &BadAPIRequest{Err: inner}
	assert.Equal(t, "mac must be greater than 5 chars", e.Error())
}

func TestBadAPIRequest_Unwrap(t *testing.T) {
	inner := errors.New("mac must be greater than 5 chars")
	e := &BadAPIRequest{Err: inner}
	assert.True(t, errors.Is(e, inner))
}

// --- BadAPIKey ---

func TestBadAPIKey_Error(t *testing.T) {
	inner := errors.New("Unauthorized")
	e := &BadAPIKey{Err: inner}
	assert.Equal(t, "Unauthorized", e.Error())
}

func TestBadAPIKey_Unwrap(t *testing.T) {
	inner := errors.New("Unauthorized")
	e := &BadAPIKey{Err: inner}
	assert.True(t, errors.Is(e, inner))
}

// --- RateLimitsExceeded ---

func TestRateLimitsExceeded_Error(t *testing.T) {
	reset := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	e := &RateLimitsExceeded{Limit: 100, Reset: reset}
	msg := e.Error()
	assert.Contains(t, msg, "100")
	assert.Contains(t, msg, "2024-01-15T12:00:00Z")
}

// --- BadAPIResponse ---

func TestBadAPIResponse_Error(t *testing.T) {
	inner := errors.New("invalid json")
	e := &BadAPIResponse{Err: inner}
	assert.Equal(t, "invalid json", e.Error())
}

func TestBadAPIResponse_Unwrap(t *testing.T) {
	inner := errors.New("invalid json")
	e := &BadAPIResponse{Err: inner}
	assert.True(t, errors.Is(e, inner))
}
