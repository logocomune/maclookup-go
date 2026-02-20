package maclookup

import "time"

// ResponseMACInfo is the full response returned by Client.Lookup.
// It embeds RateLimit (current quota information) and MACInfo (vendor data).
// RespTime records the total round-trip duration of the HTTP request.
type ResponseMACInfo struct {
	RespTime time.Duration
	RateLimit
	MACInfo
}

// ResponseVendorName is the response returned by Client.CompanyName.
// It embeds RateLimit (current quota information) and CompanyInfo (vendor name).
// RespTime records the total round-trip duration of the HTTP request.
type ResponseVendorName struct {
	RespTime time.Duration
	RateLimit
	CompanyInfo
}

// RateLimit carries the rate-limit metadata returned by the API on every response.
// Limit is the maximum number of requests allowed in the current window.
// Remaining is how many requests are still available.
// Reset is the UTC time at which the counter resets.
type RateLimit struct {
	Limit     int64
	Remaining int64
	Reset     time.Time
}

// MACInfo contains the vendor registration data associated with a MAC prefix.
// Found is false when the prefix is not present in the database; in that case
// all other string and numeric fields are zero values.
type MACInfo struct {
	Found      bool
	MacPrefix  string
	Company    string
	Address    string
	Country    string
	BlockStart string
	BlockEnd   string
	BlockSize  int
	BlockType  string
	Updated    string
	IsRand     bool
	IsPrivate  bool
}

// CompanyInfo contains the vendor name associated with a MAC prefix as returned
// by the lightweight company-name endpoint.
// Found is false when the prefix is not in the database.
// IsPrivate is true when the block is marked as a private address range.
// Company holds the registrant name; it is empty when Found is false or IsPrivate is true.
type CompanyInfo struct {
	Found     bool
	IsPrivate bool
	Company   string
}
