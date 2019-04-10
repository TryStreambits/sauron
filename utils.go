package sauron

import (
	"net/http"
	"net/url"
	"time"
)

// This file contains various utilities for Sauron

// NewHTTPClient will create a new request-specific client, with our defined user agent, for the purposes of page fetching.
// If successful, it will return both the client and the request for use
func NewHTTPClient(u *url.URL) (client http.Client, request http.Request) {
	client = http.Client{
		Timeout: time.Second * 15, // 15 seconds
	}

	var requestHeaders = make(http.Header)
	requestHeaders.Set("User-Agent", "Sauron Bot 0.1")

	request = http.Request{
		Header: requestHeaders,
		Method: "GET",
		URL:    u,
	}

	return
}
