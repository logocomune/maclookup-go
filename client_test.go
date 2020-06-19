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

func TestClient_WithPrefixUri(t *testing.T) {
	c := New()
	assert.NotEmpty(t, c.prefixURI)
	assert.Equal(t, c.prefixURI, "https://api.maclookup.app")

	c.WithPrefixURI("https://example.org")
	assert.Equal(t, c.prefixURI, "https://example.org")

	c.WithPrefixURI("https://example.org/")
	assert.Equal(t, c.prefixURI, "https://example.org")

	c.WithPrefixURI("example.org/")
	assert.Equal(t, c.prefixURI, "https://example.org")

	c.WithPrefixURI("http://example.org/")
	assert.Equal(t, c.prefixURI, "http://example.org")

	c.WithPrefixURI("127.0.0.1")
	assert.Equal(t, c.prefixURI, "http://127.0.0.1")

	c.WithPrefixURI("127.0.0.1:8080")
	assert.Equal(t, c.prefixURI, "http://127.0.0.1:8080")

	c.WithPrefixURI("::1")
	assert.Equal(t, c.prefixURI, "http://::1")
}

func TestClient_WithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)

		fmt.Fprintln(w, `{"success":true,"found":true,"macPrefix":"000000","company":"XEROX CORPORATION","address":"M/S 105-50C, WEBSTER NY 14580, US","country":"US","blockStart":"000000000000","blockEnd":"000000FFFFFF","blockSize":16777215,"blockType":"MA-L","updated":"2015-11-17","isRand":false,"isPrivate":false}`)
	}))
	defer ts.Close()

	timeout := 5 * time.Millisecond
	client := New()
	client.WithPrefixURI(ts.URL)
	client.WithTimeout(timeout)
	assert.Equal(t, client.timeOut, timeout)
	_, err := client.Lookup("000000")
	assert.NotNil(t, err)

	var e *HTTPClientError

	assert.True(t, errors.As(err, &e))
}

func Test_cleanMac(t *testing.T) {
	type args struct {
		mac string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Mac #1",
			args: args{mac: "000000"},
			want: "000000",
		},
		{
			name: "Mac #2",
			args: args{mac: "00:00:00"},
			want: "000000",
		},
		{
			name: "Mac #3",
			args: args{mac: "00.00.00"},
			want: "000000",
		},
		{
			name: "Mac #4",
			args: args{mac: "000.000"},
			want: "000000",
		},
		{
			name: "Mac #5",
			args: args{mac: "00-00-00"},
			want: "000000",
		},
		{
			name: "Mac #6",
			args: args{mac: "0A-0C:cc"},
			want: "0A0CCC",
		},
		{
			name: "Mac #7",
			args: args{mac: "0A-0C:cc 0"},
			want: "0A0CCC0",
		},
		{
			name: "Mac #8",
			args: args{mac: "0A-0C:cc 0a"},
			want: "0A0CCC0",
		},
		{
			name: "Mac #8",
			args: args{mac: "0A-0C:cc 0a"},
			want: "0A0CCC0",
		},
		{
			name: "Mac #9",
			args: args{mac: "0A-0C:cc 0aAb"},
			want: "0A0CCC0AA",
		},
		{
			name: "Mac #10",
			args: args{mac: "0A-0C:cc 0aA"},
			want: "0A0CCC0AA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanMac(tt.args.mac); got != tt.want {
				t.Errorf("cleanMac() = %v, want %v", got, tt.want)
			}
		})
	}
}
