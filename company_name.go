package maclookup

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//CompanyName returns company name from API.
func (c Client) CompanyName(mac string) (ResponseVendorName, error) {
	url := c.prefixURI + apiMAC + cleanMac(mac) + companyNameSuffix
	if c.apiKey != "" {
		url += apiKeyParam + c.apiKey
	}

	return c.getCompanyName(url)
}

func (c Client) getCompanyName(url string) (ResponseVendorName, error) {
	var response ResponseVendorName

	start := time.Now()
	timeout, cancell := context.WithTimeout(context.Background(), c.timeOut)
	defer cancell()

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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, &HTTPClientError{Err: err}
	}

	body := string(bodyBytes)
	response.RespTime = time.Since(start)

	statusCode := resp.StatusCode

	switch statusCode {
	case http.StatusBadRequest:
		msg := "client request error"
		if body != "" {
			msg = body
		}

		return response, &BadAPIRequest{Err: errors.New(strings.ToLower(msg))}

	case http.StatusUnauthorized:
		msg := "bad api key"
		if body != "" {
			msg = body
		}

		return response, &BadAPIKey{Err: errors.New(msg)}

	case http.StatusTooManyRequests:
		return response, &RateLimitsExceeded{
			Limit: response.RateLimit.Limit,
			Reset: response.RateLimit.Reset,
		}

	case http.StatusNotFound:
		return response, &HTTPClientError{Err: errors.New("endpoint not found")}
	case http.StatusOK:
		response.Found = !(body == "*NO COMPANY*")
		response.IsPrivate = body == "*PRIVATE*"

		if response.Found && !response.IsPrivate {
			response.Company = body
		}

		return response, nil
	}

	return response, &HTTPClientError{Err: errors.New("unexpected http status: " + strconv.Itoa(statusCode))}
}
