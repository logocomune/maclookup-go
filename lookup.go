package maclookup

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Lookup retrieve MAC information from API.
func (c Client) Lookup(mac string) (ResponseMACInfo, error) {
	url := c.prefixURI + apiMAC + cleanMac(mac)
	if c.apiKey != "" {
		url += apiKeyParam + c.apiKey
	}

	return c.getMacInfo(url)
}

func (c Client) getMacInfo(url string) (ResponseMACInfo, error) {
	var response ResponseMACInfo

	start := time.Now()
	timeout, cancel := context.WithTimeout(context.Background(), c.timeOut)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, "GET", url, nil)
	if err != nil {
		return response, &HTTPClientError{Err: err}
	}

	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*")
	resp, err := c.client.Do(req)

	if err != nil {
		return response, &HTTPClientError{Err: err}
	}
	defer resp.Body.Close()

	response.RateLimit = RateLimit{
		Limit:     parseLimit(resp.Header.Get(xRateLimit)),
		Remaining: parseIntHeader(resp.Header, xRateRemaining),
		Reset:     parseTimeHeader(resp.Header, xRateReset),
	}

	if err := checkStatusMacInfo(resp.StatusCode, resp.Body, response.Limit, response.Reset); err != nil {
		return response, err
	}

	var apiRespose maclookupResponseAPIV2

	err = json.NewDecoder(resp.Body).Decode(&apiRespose)

	response.RespTime = time.Since(start)

	if err != nil || !apiRespose.Success {
		return response, &BadAPIResponse{Err: err}
	}

	//Decoupling api response
	response.Found = apiRespose.Found
	response.MacPrefix = apiRespose.MacPrefix
	response.Company = apiRespose.Company
	response.Address = apiRespose.Address
	response.Country = apiRespose.Country
	response.BlockStart = apiRespose.BlockStart
	response.BlockEnd = apiRespose.BlockEnd
	response.BlockSize = apiRespose.BlockSize
	response.BlockType = apiRespose.BlockType
	response.Updated = apiRespose.Updated
	response.IsRand = apiRespose.IsRand
	response.IsPrivate = apiRespose.IsPrivate

	return response, nil
}

func checkStatusMacInfo(statusCode int, body io.Reader, rateLimit int64, rateLimitReset time.Time) error {
	switch statusCode {
	case http.StatusBadRequest:
		var e errorResponseAPIV2
		err := json.NewDecoder(body).Decode(&e)
		msg := "client request error"

		if err == nil {
			msg = e.Error
		}

		return &BadAPIRequest{Err: errors.New(strings.ToLower(msg))}

	case http.StatusUnauthorized:
		var e errorResponseAPIV2
		err := json.NewDecoder(body).Decode(&e)
		msg := "bad api key"

		if err == nil {
			msg = e.Error
		}

		return &BadAPIKey{Err: errors.New(msg)}
	case http.StatusTooManyRequests:
		return &RateLimitsExceeded{
			Limit: rateLimit,
			Reset: rateLimitReset,
		}

	case http.StatusOK:
		return nil
	}

	return &HTTPClientError{Err: errors.New("unexpected http status: " + strconv.Itoa(statusCode))}
}
