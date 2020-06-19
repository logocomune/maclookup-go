package maclookup

import "time"

type ResponseMACInfo struct {
	RespTime time.Duration
	RateLimit
	MACInfo
}

type ResponseVendorName struct {
	RespTime time.Duration
	RateLimit
	CompanyInfo
}

type RateLimit struct {
	Limit     int64
	Remaining int64
	Reset     time.Time
}

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

type CompanyInfo struct {
	Found     bool
	IsPrivate bool
	Company   string
}
