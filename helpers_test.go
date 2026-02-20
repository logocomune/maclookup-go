package maclookup

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- isIP ---

func Test_isIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.1:8080", true},
		{"192.168.1.100", true},
		{"::1", true},          // IPv6 loopback (no port, 3 colons → len>2 path)
		{"example.com", false}, // hostname
		{"example.com:443", false},
		{"not-an-ip", false},
		{"", false},
	}

	for _, tt := range tests {
		got := isIP(tt.input)
		assert.Equal(t, tt.want, got, "isIP(%q)", tt.input)
	}
}

// --- parseIntHeader ---

func Test_parseIntHeader(t *testing.T) {
	h := http.Header{}
	h.Set("X-Test", "42")
	assert.Equal(t, int64(42), parseIntHeader(h, "X-Test"))

	h.Set("X-Bad", "not-a-number")
	assert.Equal(t, int64(-1), parseIntHeader(h, "X-Bad"))

	// missing key
	assert.Equal(t, int64(-1), parseIntHeader(h, "X-Missing"))
}

// --- parseLimit ---

func Test_parseLimit(t *testing.T) {
	assert.Equal(t, int64(10), parseLimit("10"))
	assert.Equal(t, int64(2), parseLimit("2, 2;window=1"))
	assert.Equal(t, int64(-1), parseLimit("not-a-number"))
	assert.Equal(t, int64(-1), parseLimit(""))
}

// --- parseTimeHeader ---

func Test_parseTimeHeader(t *testing.T) {
	now := time.Now()
	h := http.Header{}
	h.Set("X-Reset", "0")
	assert.Equal(t, time.Unix(0, 0), parseTimeHeader(h, "X-Reset"))

	h.Set("X-Reset-Now", fmt.Sprint(now.Unix()))
	assert.Equal(t, time.Unix(now.Unix(), 0), parseTimeHeader(h, "X-Reset-Now"))

	h.Set("X-Bad", "not-a-timestamp")
	assert.Equal(t, time.Time{}, parseTimeHeader(h, "X-Bad"))

	assert.Equal(t, time.Time{}, parseTimeHeader(h, "X-Missing"))
}
