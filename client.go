package maclookup

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	timeOut = 5 * time.Second
	ua      = "MACLookupClient-go/1.0.0 (https://maclookup.app)"

	xRateLimit     = "X-RateLimit-Limit"
	xRateRemaining = "X-RateLimit-Remaining"
	xRateReset     = "X-RateLimit-Reset"
)

// Client is the maclookup.app API client.
// Create one with New() and optionally configure it with WithAPIKey, WithTimeout, or WithPrefixURI.
type Client struct {
	client    *http.Client
	apiKey    string
	prefixURI string
	timeOut   time.Duration
}

// New creates a new Client for the maclookup.app API v2 using the default HTTP
// client and a 5-second timeout. No API key is set by default; free-tier rate
// limits apply until one is provided via WithAPIKey.
func New() *Client {
	client := http.DefaultClient

	return &Client{
		client:    client,
		prefixURI: apiURIPrefix,
		timeOut:   timeOut,
	}
}

// WithAPIKey sets the API key used for authenticated requests.
// Obtain a key at https://maclookup.app/api-v2/plans.
func (c *Client) WithAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// WithTimeout overrides the per-request HTTP timeout (default: 5 s).
func (c *Client) WithTimeout(timeout time.Duration) {
	c.timeOut = timeout
}

// WithPrefixURI replaces the default API base URL (https://api.maclookup.app).
// Useful for testing or when routing through a proxy. If the supplied string
// starts with an IP address (and no explicit scheme), http:// is used;
// otherwise https:// is added when no scheme is present.
func (c *Client) WithPrefixURI(prefixURI string) {
	prefix := strings.TrimRight(prefixURI, "/")

	if strings.HasPrefix(prefixURI, "http://") || strings.HasPrefix(prefixURI, "https://") {
		c.prefixURI = prefix
		return
	}

	c.prefixURI = "https://" + prefix

	if isIP(prefix) {
		c.prefixURI = "http://" + prefix
	}
}

// isIP reports whether host (without scheme) is an IP address, optionally
// followed by a port (host:port) or path (host/path). Two or more colons mean
// IPv6.
func isIP(host string) bool {
	colons := strings.Count(host, ":")
	if colons > 1 {
		// IPv6 literal — pass as-is to net.ParseIP.
		return net.ParseIP(host) != nil
	}
	if colons == 1 {
		// Strip port suffix.
		host = host[:strings.IndexByte(host, ':')]
	}
	// Strip path suffix.
	if i := strings.IndexByte(host, '/'); i >= 0 {
		host = host[:i]
	}
	return net.ParseIP(host) != nil
}

func parseIntHeader(header http.Header, property string) int64 {
	parseInt, err := strconv.ParseInt(header.Get(property), 10, 64)
	if err != nil {
		return -1
	}

	return parseInt
}

// parseLimit extracts the first numeric token from an X-RateLimit-Limit header
// value, which may carry an optional policy string (e.g. "2, 2;window=1").
// Using IndexByte avoids the slice allocation of strings.Split.
func parseLimit(limit string) int64 {
	if i := strings.IndexByte(limit, ','); i >= 0 {
		limit = limit[:i]
	}
	parseInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil {
		return -1
	}

	return parseInt
}

func parseTimeHeader(header http.Header, property string) time.Time {
	parseInt, err := strconv.ParseInt(header.Get(property), 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(parseInt, 0)
}

// cleanMac normalises a MAC address string to a bare uppercase hex prefix of
// the appropriate length (6, 7 or 9 chars) by stripping separators (:, ., -,
// space) in a single pass over the bytes. A fixed [9]byte stack buffer avoids
// heap allocation for the common case.
func cleanMac(mac string) string {
	var buf [9]byte
	n := 0

	for i := 0; i < len(mac); i++ {
		c := mac[i]
		switch {
		case c == ':' || c == '.' || c == '-' || c == ' ':
			continue
		case c >= 'a' && c <= 'z':
			buf[n] = c - 32 // to upper
			n++
		default:
			buf[n] = c
			n++
		}
		if n == 9 {
			return string(buf[:9])
		}
	}

	switch {
	case n >= 7:
		return string(buf[:7])
	case n >= 6:
		return string(buf[:6])
	default:
		return string(buf[:n])
	}
}
