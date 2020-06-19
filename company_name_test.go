package maclookup

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_CompanyNameGoodResponse(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))

		fmt.Fprint(w, `XEROX CORPORATION`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	cName, err := client.CompanyName("000000")

	assert.Nil(t, err)
	assert.NotEmpty(t, cName)
	assert.True(t, cName.Found)
	assert.Equal(t, cName.Company, "XEROX CORPORATION")
	assert.False(t, cName.IsPrivate)

	cName, err = client.CompanyName("00:00:00")

	assert.Nil(t, err)
	assert.NotEmpty(t, cName)
	assert.True(t, cName.Found)
	assert.Equal(t, cName.Company, "XEROX CORPORATION")
	assert.False(t, cName.IsPrivate)
}

func TestClient_CompanyNameNotFound(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `*NO COMPANY*`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	cName, err := client.CompanyName("000000")

	assert.Nil(t, err)
	assert.NotEmpty(t, cName)
	assert.False(t, cName.Found)
	assert.Equal(t, cName.Company, "")
	assert.False(t, cName.IsPrivate)
}

func TestClient_CompanyNamePrivate(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `*PRIVATE*`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	cName, err := client.CompanyName("000000")

	assert.Nil(t, err)
	assert.NotEmpty(t, cName)
	assert.True(t, cName.Found)
	assert.Equal(t, cName.Company, "")
	assert.True(t, cName.IsPrivate)
}

func TestClient_CompanyName404BadResponse(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, `404. That’s an error.`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	_, err := client.CompanyName("000000")

	var e *HTTPClientError

	assert.True(t, errors.As(err, &e))
}

func TestClient_CompanyNameBadAPIKey(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name?apiKey=BAD_API_KEY", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `Bad APIKey - Unauthorized - more info: https://maclookup.app/api-v2/plans`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	client.WithAPIKey("BAD_API_KEY")
	_, err := client.CompanyName("000000")
	assert.NotNil(t, err)

	var e *BadAPIKey

	assert.True(t, errors.As(err, &e))
}

func TestClient_CompanyNameBadRequest(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/0000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `MAC must be grater than 5 chars`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	_, err := client.CompanyName("0000")

	assert.NotNil(t, err)

	var e *BadAPIRequest

	assert.True(t, errors.As(err, &e))
}

func TestClient_CompanyNameExceededRateLimit(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "2, 2;window=1")
		w.Header().Add(xRateRemaining, "0")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, `Too Many Requests - more info: https://maclookup.app/api-v2/rate-limits`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	cName, err := client.CompanyName("000000")
	assert.NotNil(t, err)

	var e *RateLimitsExceeded

	assert.True(t, errors.As(err, &e))
	assert.Equal(t, cName.RateLimit.Limit, int64(2))
	assert.Equal(t, cName.RateLimit.Reset, time.Unix(now.Unix(), 0))
}

func TestClient_CompanyNameBadResponse(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000/company/name", r.RequestURI)
		w.Header().Add(xRateLimit, "2, 2;window=1")
		w.Header().Add(xRateRemaining, "3")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `SYSTEM ERROR`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	_, err := client.CompanyName("000000")
	assert.NotNil(t, err)

	var e *HTTPClientError

	assert.True(t, errors.As(err, &e))
}

func ExampleClient_CompanyName() {
	//Prevent rate limits error
	time.Sleep(time.Millisecond * 550)

	client := New()
	r, err := client.CompanyName("000000")
	fmt.Println(err)
	fmt.Printf("%s", r.Company)
	//Output:
	//<nil>
	//XEROX CORPORATION
}

func ExampleClient_CompanyNameNotFound() {
	//Prevent rate limits error
	time.Sleep(time.Millisecond * 550)

	client := New()
	r, err := client.CompanyName("010000")
	fmt.Println(err)
	fmt.Printf("%s\n", r.Company)
	fmt.Printf("%t", r.Found)
	//Output:
	//<nil>
	//
	//false
}
