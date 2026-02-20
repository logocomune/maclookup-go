package maclookup

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Internal helper benchmarks
// ---------------------------------------------------------------------------

var macInputs = []string{
	"000000",
	"00:00:00",
	"00-00-00",
	"00.00.00",
	"AA:BB:CC:DD:EE:FF",
	"aa-bb-cc-dd-ee-ff",
	"0A-0C:cc 0aAb",
	"  00:11:22:33:44:55  ",
}

func BenchmarkCleanMac(b *testing.B) {
	for _, mac := range macInputs {
		mac := mac
		b.Run(mac, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = cleanMac(mac)
			}
		})
	}
}

func BenchmarkParseLimit(b *testing.B) {
	cases := []string{
		"10",
		"2, 2;window=1",
		"not-a-number",
	}
	for _, c := range cases {
		c := c
		b.Run(c, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = parseLimit(c)
			}
		})
	}
}

func BenchmarkParseIntHeader(b *testing.B) {
	h := http.Header{}
	h.Set("X-RateLimit-Remaining", "42")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseIntHeader(h, "X-RateLimit-Remaining")
	}
}

func BenchmarkParseTimeHeader(b *testing.B) {
	h := http.Header{}
	h.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Unix()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseTimeHeader(h, "X-RateLimit-Reset")
	}
}

func BenchmarkIsIP(b *testing.B) {
	cases := []string{
		"127.0.0.1",
		"127.0.0.1:8080",
		"::1",
		"example.com",
		"example.com:443",
	}
	for _, c := range cases {
		c := c
		b.Run(c, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = isIP(c)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// HTTP endpoint benchmarks (no network — httptest.Server)
// ---------------------------------------------------------------------------

const lookupResponseJSON = `{"success":true,"found":true,"macPrefix":"000000","company":"XEROX CORPORATION","address":"M/S 105-50C, WEBSTER NY 14580, US","country":"US","blockStart":"000000000000","blockEnd":"000000FFFFFF","blockSize":16777215,"blockType":"MA-L","updated":"2015-11-17","isRand":false,"isPrivate":false}`

func newLookupServer(b *testing.B) *httptest.Server {
	b.Helper()
	now := time.Now()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(xRateLimit, "100")
		w.Header().Set(xRateRemaining, "99")
		w.Header().Set(xRateReset, fmt.Sprintf("%d", now.Unix()))
		fmt.Fprintln(w, lookupResponseJSON)
	}))
}

func newCompanyNameServer(b *testing.B) *httptest.Server {
	b.Helper()
	now := time.Now()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(xRateLimit, "100")
		w.Header().Set(xRateRemaining, "99")
		w.Header().Set(xRateReset, fmt.Sprintf("%d", now.Unix()))
		fmt.Fprint(w, "XEROX CORPORATION")
	}))
}

func BenchmarkLookup(b *testing.B) {
	ts := newLookupServer(b)
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Lookup("000000")
	}
}

func BenchmarkLookupParallel(b *testing.B) {
	ts := newLookupServer(b)
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.Lookup("000000")
		}
	})
}

func BenchmarkCompanyName(b *testing.B) {
	ts := newCompanyNameServer(b)
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.CompanyName("000000")
	}
}

func BenchmarkCompanyNameParallel(b *testing.B) {
	ts := newCompanyNameServer(b)
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.CompanyName("000000")
		}
	})
}

// BenchmarkLookupVariousMAC measures the full pipeline including MAC cleaning
// across several input formats.
func BenchmarkLookupVariousMAC(b *testing.B) {
	ts := newLookupServer(b)
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	macs := []string{"000000", "00:00:00", "00-00-00", "00.00.00"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Lookup(macs[i%len(macs)])
	}
}

// BenchmarkLookupNotFound benchmarks the path where the API returns found=false.
func BenchmarkLookupNotFound(b *testing.B) {
	now := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(xRateLimit, "100")
		w.Header().Set(xRateRemaining, "99")
		w.Header().Set(xRateReset, fmt.Sprintf("%d", now.Unix()))
		fmt.Fprintln(w, `{"success":true,"found":false,"isRand":false,"isPrivate":false}`)
	}))
	defer ts.Close()

	client := New()
	client.WithPrefixURI(ts.URL)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.Lookup("010000")
	}
}
