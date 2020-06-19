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

func TestClient_LookupGoodResponse(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))

		fmt.Fprintln(w, `{"success":true,"found":true,"macPrefix":"000000","company":"XEROX CORPORATION","address":"M/S 105-50C, WEBSTER NY 14580, US","country":"US","blockStart":"000000000000","blockEnd":"000000FFFFFF","blockSize":16777215,"blockType":"MA-L","updated":"2015-11-17","isRand":false,"isPrivate":false}`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	macInfo, err := client.Lookup("000000")

	assert.Nil(t, err)
	assert.NotEmpty(t, macInfo)
	assert.True(t, macInfo.Found)
	assert.Equal(t, macInfo.MacPrefix, "000000")
	assert.Equal(t, macInfo.Company, "XEROX CORPORATION")
	assert.Equal(t, macInfo.Address, "M/S 105-50C, WEBSTER NY 14580, US")
	assert.Equal(t, macInfo.Country, "US")
	assert.Equal(t, macInfo.BlockStart, "000000000000")
	assert.Equal(t, macInfo.BlockEnd, "000000FFFFFF")
	assert.Equal(t, macInfo.BlockSize, 16777215)
	assert.Equal(t, macInfo.BlockType, "MA-L")
	assert.Equal(t, macInfo.Updated, "2015-11-17")
	assert.False(t, macInfo.IsRand)
	assert.False(t, macInfo.IsRand)
	assert.Equal(t, macInfo.RateLimit.Limit, int64(10))
	assert.Equal(t, macInfo.RateLimit.Remaining, int64(9))
	assert.Equal(t, macInfo.RateLimit.Reset, time.Unix(now.Unix(), 0))

	macInfo, err = client.Lookup("00:00:00")

	assert.Nil(t, err)
	assert.NotEmpty(t, macInfo)
	assert.True(t, macInfo.Found)
	assert.Equal(t, macInfo.MacPrefix, "000000")
	assert.Equal(t, macInfo.Company, "XEROX CORPORATION")
	assert.Equal(t, macInfo.Address, "M/S 105-50C, WEBSTER NY 14580, US")
	assert.Equal(t, macInfo.Country, "US")
	assert.Equal(t, macInfo.BlockStart, "000000000000")
	assert.Equal(t, macInfo.BlockEnd, "000000FFFFFF")
	assert.Equal(t, macInfo.BlockSize, 16777215)
	assert.Equal(t, macInfo.BlockType, "MA-L")
	assert.Equal(t, macInfo.Updated, "2015-11-17")
	assert.False(t, macInfo.IsRand)
	assert.False(t, macInfo.IsRand)
}

func TestClient_LookupGooNotFound(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))

		fmt.Fprintln(w, `{"success":true,"found":false,"isRand":false}`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	macInfo, err := client.Lookup("000000")

	assert.Nil(t, err)
	assert.NotEmpty(t, macInfo)
	assert.False(t, macInfo.Found)
	assert.Empty(t, macInfo.MacPrefix)
	assert.Empty(t, macInfo.Company, "XEROX CORPORATION")
	assert.Empty(t, macInfo.Address, "M/S 105-50C, WEBSTER NY 14580, US")
	assert.Empty(t, macInfo.Country, "US")
	assert.Empty(t, macInfo.BlockStart, "000000000000")
	assert.Empty(t, macInfo.BlockEnd, "000000FFFFFF")
	assert.Empty(t, macInfo.BlockSize, 16777215)
	assert.Empty(t, macInfo.BlockType, "MA-L")
	assert.Empty(t, macInfo.Updated, "2015-11-17")
	assert.False(t, macInfo.IsRand)
	assert.False(t, macInfo.IsPrivate)
	assert.Equal(t, macInfo.RateLimit.Limit, int64(10))
	assert.Equal(t, macInfo.RateLimit.Remaining, int64(9))
	assert.Equal(t, macInfo.RateLimit.Reset, time.Unix(now.Unix(), 0))
}

func TestClient_LookupBadAPIKey(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000?apiKey=BAD_API_KEY", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, `{"success":false,"error":"Unauthorized","errorCode":401,"moreInfo":"https://maclookup.app/api-v2/plans"}`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	client.WithAPIKey("BAD_API_KEY")
	_, err := client.Lookup("000000")
	assert.NotNil(t, err)

	var e *BadAPIKey

	assert.True(t, errors.As(err, &e))
}

func TestClient_LookupBadRequest(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/0000", r.RequestURI)
		w.Header().Add(xRateLimit, "10")
		w.Header().Add(xRateRemaining, "9")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"success":false,"error":"MAC must be grater than 5 chars","errorCode":101,"moreInfo":"https://maclookup.app/api-v2/documentation"}`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	_, err := client.Lookup("0000")
	assert.NotNil(t, err)

	var e *BadAPIRequest

	assert.True(t, errors.As(err, &e))
}

func TestClient_LookupExceededRateLimit(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000", r.RequestURI)
		w.Header().Add(xRateLimit, "2, 2;window=1")
		w.Header().Add(xRateRemaining, "0")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, `{"success":false,"error":"Too Many Requests","errorCode":429,"moreInfo":"https://maclookup.app/api-v2/rate-limits"}`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	resp, err := client.Lookup("000000")
	assert.NotNil(t, err)

	var e *RateLimitsExceeded

	assert.True(t, errors.As(err, &e))
	assert.Equal(t, resp.RateLimit.Limit, int64(2))
	assert.Equal(t, resp.RateLimit.Reset, time.Unix(now.Unix(), 0))
}

func TestClient_LookupBadResponse(t *testing.T) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/macs/000000", r.RequestURI)
		w.Header().Add(xRateLimit, "2, 2;window=1")
		w.Header().Add(xRateRemaining, "3")
		w.Header().Add(xRateReset, fmt.Sprintf("%d", now.Unix()))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `SYSTEM ERROR`)
	}))

	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)
	_, err := client.Lookup("000000")
	assert.NotNil(t, err)

	var e *BadAPIResponse

	assert.True(t, errors.As(err, &e))
}

func ExampleClient_Lookup() {
	//Prevent rate limits error
	time.Sleep(time.Millisecond * 550)

	client := New()
	r, err := client.Lookup("000000")
	fmt.Println(err)
	fmt.Printf("%+v", r.MACInfo)
	//Output:
	//<nil>
	//{Found:true MacPrefix:000000 Company:XEROX CORPORATION Address:M/S 105-50C, WEBSTER NY 14580, US Country:US BlockStart:000000000000 BlockEnd:000000FFFFFF BlockSize:16777215 BlockType:MA-L Updated:2015-11-17 IsRand:false IsPrivate:false}
}

func ExampleClient_Lookup_NotFound() {
	//Prevent rate limits error
	time.Sleep(time.Millisecond * 550)

	client := New()
	r, err := client.Lookup("010000")
	fmt.Println(err)
	fmt.Printf("%+v", r.MACInfo)
	//Output:
	//<nil>
	//{Found:false MacPrefix: Company: Address: Country: BlockStart: BlockEnd: BlockSize:0 BlockType: Updated: IsRand:false IsPrivate:false}
}
