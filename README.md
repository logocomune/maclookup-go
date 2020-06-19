# maclookup-go
[![Build Status](https://travis-ci.org/logocomune/maclookup-go.svg?branch=master)](https://travis-ci.org/logocomune/maclookup-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/maclookup-go)](https://goreportcard.com/report/github.com/logocomune/maclookup-go)
[![codecov](https://codecov.io/gh/logocomune/maclookup-go/branch/master/graph/badge.svg)](https://codecov.io/gh/logocomune/maclookup-go)



A Go library for interacting with [MACLookup's API v2](https://maclookup.app/api-v2/documentation). This library allows you to:

- Get full info (MAC prefix, company name, address and country) of a MAC address
- Get Company name by MAC

## Installation

You need a working Go environment.

```shell
 go get github.com/logocomune/maclookup-go
````

##Getting Started

```go
package main

import (
	"log"

	"github.com/logocomune/maclookup-go"
)

func main() {
	client := maclookup.New()
    
	r, err := client.CompanyName("000000")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("MAC found in database:", r.Found)
	log.Println("MAC is private (no company name):", r.IsPrivate)
	log.Println("Company name:", r.Company)
	log.Println("Api response in: ", r.RespTime)
	log.Println("Rate limits - remaining request for current time window:", r.RateLimit.Remaining)
	log.Println("Rate limits - next reset", r.RateLimit.Reset)

}


```

###Use custom timout
```go
    client := maclookup.New()
    client.WithTimeout(10*time.Second) 
```

###API Key
Get an API Key [here](https://maclookup.app/api-v2/plans)
```go
    client := maclookup.New()
	client.WithAPIKey("an_api_key")

```


##Example

- [Get full info of a MAC](/example/lookup)  
- [Get company name](/example/company-name)  
- [Rate limit](/example/rate-limit)  
