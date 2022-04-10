// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package links provides a higher level interface to shortening links using
// the links.wtf shortener service (with support for using your custom links.wtf
// service as well.)
package links

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// URI is the endpoint which we use to shorten the link. Must be stripped of
// trailing slashes. Defaults to "https://links.wtf".
var URI = "https://links.wtf"

type apiResponse struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Shorten shortens a supplied URL with the given links.wtf service, optionally
// with a given encryption string. If you run a custom links.wtf service,
// make sure to update the links.URL variable with the URI of where that
// service is located.
//
// If passwd is blank, no encryption is used. Supply your own httpClient
// to utilize a proxy, or change the timeout (defaults to 4 seconds if no
// config is supplied.)
func Shorten(link, passwd string, httpClient *http.Client) (uri *url.URL, err error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 4 * time.Second}
	}

	params := url.Values{}

	params.Set("url", link)
	if passwd != "" {
		params.Set("encrypt", passwd)
	}

	resp, err := httpClient.PostForm(URI+"/add", params)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &apiResponse{}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		if result.Message == "" {
			return nil, errors.New("api returned unknown unsuccessful response")
		}

		return nil, fmt.Errorf("api returned error: %s", result.Message)
	}

	if uri, err = url.Parse(result.URL); err != nil {
		return nil, fmt.Errorf("api returned invalid uri: %s", result.URL)
	}

	return uri, nil
}
