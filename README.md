# maclookup-go

[![CircleCI](https://circleci.com/gh/logocomune/maclookup-go/tree/main.svg?style=svg)](https://circleci.com/gh/logocomune/maclookup-go/tree/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/maclookup-go)](https://goreportcard.com/report/github.com/logocomune/maclookup-go)
[![codecov](https://codecov.io/gh/logocomune/maclookup-go/branch/master/graph/badge.svg)](https://codecov.io/gh/logocomune/maclookup-go)
[![CodeFactor](https://www.codefactor.io/repository/github/logocomune/maclookup-go/badge)](https://www.codefactor.io/repository/github/logocomune/maclookup-go)

A Go client library for the [maclookup.app API v2](https://maclookup.app/api-v2/documentation).

Given any MAC address (or just its OUI/prefix), the library resolves the IEEE-registered vendor information — company name, address, country, block type and range — in a single HTTP call. It handles rate-limit headers, API key authentication, timeouts, and all documented error responses.

## Features

- **Full MAC lookup** — resolve MAC prefix → company name, address, country, block start/end/size/type, last-updated date, `isRand`, `isPrivate` flags.
- **Lightweight company-name lookup** — retrieve only the vendor name (faster, cheaper on quota).
- **Rate-limit awareness** — every response exposes `RateLimit.Limit`, `Remaining`, and `Reset` so callers can back off gracefully.
- **Flexible MAC input** — accepts any common notation: `AA:BB:CC:DD:EE:FF`, `AA-BB-CC`, `AABBCC`, `AA.BB.CC`, etc.
- **Typed errors** — distinguish network errors, bad requests, invalid API keys, and rate-limit violations with `errors.As`.
- **Configurable** — custom timeout, API key, or base URL (useful for testing against a local proxy).

## Installation

Requires Go 1.16 or later.

```shell
go get github.com/logocomune/maclookup-go
```

## Quick Start

### Look up full vendor information

```go
package main

import (
    "fmt"
    "log"

    "github.com/logocomune/maclookup-go"
)

func main() {
    client := maclookup.New()

    r, err := client.Lookup("00:00:00:00:00:00")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Found:      ", r.Found)
    fmt.Println("Company:    ", r.Company)
    fmt.Println("Address:    ", r.Address)
    fmt.Println("Country:    ", r.Country)
    fmt.Println("Block type: ", r.BlockType)
    fmt.Println("Updated:    ", r.Updated)
    fmt.Println("Private:    ", r.IsPrivate)
    fmt.Println("Random:     ", r.IsRand)
    fmt.Println("Response time:", r.RespTime)
}
```

### Look up company name only

```go
client := maclookup.New()

r, err := client.CompanyName("00:00:00")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Found:   ", r.Found)
fmt.Println("Private: ", r.IsPrivate)
fmt.Println("Company: ", r.Company)
```

## Configuration

### API key

Free-tier accounts are subject to strict rate limits. Register at <https://maclookup.app/api-v2/plans> and pass your key to the client:

```go
client := maclookup.New()
client.WithAPIKey("your_api_key_here")
```

### Custom timeout

The default per-request timeout is 5 seconds.

```go
client := maclookup.New()
client.WithTimeout(10 * time.Second)
```

### Custom base URL

Useful for routing through a proxy or pointing at a local mock server during tests:

```go
client := maclookup.New()
client.WithPrefixURI("http://localhost:8080")
```

## Error Handling

All errors returned by `Lookup` and `CompanyName` are typed and can be inspected with `errors.As`:

| Error type | When returned |
|---|---|
| `*maclookup.BadAPIRequest` | HTTP 400 — invalid MAC or query parameter |
| `*maclookup.BadAPIKey` | HTTP 401 — missing or invalid API key |
| `*maclookup.RateLimitsExceeded` | HTTP 429 — quota exhausted; check `e.Limit` and `e.Reset` |
| `*maclookup.BadAPIResponse` | HTTP 200 but unreadable response body |
| `*maclookup.HTTPClientError` | Network/transport error or unexpected HTTP status |

```go
import (
    "errors"
    "log"

    "github.com/logocomune/maclookup-go"
)

client := maclookup.New()

r, err := client.Lookup("00:00:00")
if err != nil {
    var rateLimitErr *maclookup.RateLimitsExceeded
    var badKeyErr    *maclookup.BadAPIKey

    switch {
    case errors.As(err, &rateLimitErr):
        log.Printf("rate limited — quota: %d, resets at: %s",
            rateLimitErr.Limit, rateLimitErr.Reset)
    case errors.As(err, &badKeyErr):
        log.Println("check your API key:", badKeyErr)
    default:
        log.Println("lookup failed:", err)
    }
    return
}

// r.Found is false when the prefix is not in the database
if !r.Found {
    log.Println("MAC prefix not found")
    return
}
log.Println(r.Company)
```

## Rate Limits

Every successful response (even when `Found` is `false`) populates `RateLimit`:

```go
r, _ := client.Lookup("00:00:00")
fmt.Println("Quota remaining:", r.RateLimit.Remaining)
fmt.Println("Resets at:      ", r.RateLimit.Reset)
```

To stay within limits, apply a client-side rate limiter (see the [rate-limit example](example/rate-limit/main.go)).

## Examples

| Example | Description |
|---|---|
| [example/lookup](example/lookup/main.go) | Full MAC lookup |
| [example/company-name](example/company-name/main.go) | Vendor name only |
| [example/rate-limit](example/rate-limit/main.go) | Client-side rate limiting with `golang.org/x/time/rate` |

## API Reference

Full Go documentation is available via `go doc` or [pkg.go.dev](https://pkg.go.dev/github.com/logocomune/maclookup-go).

## License

MIT — see [LICENSE](LICENSE).
