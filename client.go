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
	ua      = "MACLookupClient/1.0.0 (https://maclookup.app)"

	xRateLimit     = "X-RateLimit-Limit"
	xRateRemaining = "X-RateLimit-Remaining"
	xRateReset     = "X-RateLimit-Reset"
)

type Client struct {
	client    *http.Client
	apiKey    string
	prefixURI string
	timeOut   time.Duration
}

//New creates a new client for maclookup.app API.
func New() *Client {
	client := http.DefaultClient

	return &Client{
		client:    client,
		prefixURI: apiURIPrefix,
		timeOut:   timeOut,
	}
}

//WithAPIKey adds apiKey to client.
func (c *Client) WithAPIKey(apiKey string) {
	c.apiKey = apiKey
}

//WithTimeout defines a new timeout value for every request.
func (c *Client) WithTimeout(timeout time.Duration) {
	c.timeOut = timeout
}

//WithPrefixURI changes the default API prefix url.
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

func isIP(host string) bool {
	h := strings.Split(host, ":")
	if len(h) <= 2 {
		h = strings.Split(h[0], "/")
		return net.ParseIP(h[0]) != nil
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

func parseLimit(limit string) int64 {
	l := strings.Split(limit, ", ")
	parseInt, err := strconv.ParseInt(l[0], 10, 64)

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

func cleanMac(mac string) string {
	chars := []string{":", ".", "-", " "}
	m := strings.TrimSpace(mac)

	for _, c := range chars {
		m = strings.Replace(m, c, "", -1)
	}

	m = strings.ToUpper(m)

	if len(m) >= 9 {
		return m[0:9]
	}

	if len(m) >= 7 {
		return m[0:7]
	}

	if len(m) >= 6 {
		return m[0:6]
	}

	return strings.ToUpper(m)
}
