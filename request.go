// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type request struct {
	url       string
	values    *url.Values
	cookie    *http.Cookie
	useragent string
}

func (r request) getResponse() (*bytes.Buffer, error) {
	// Determine the HTTP action.
	var action, finalurl string
	if r.values == nil {
		action = "GET"
		finalurl = r.url
	} else {
		action = "POST"
		finalurl = r.url + "?" + r.values.Encode()
	}

	// Create a request and add the proper headers.
	req, err := http.NewRequest(action, finalurl, nil)
	if err != nil {
		return nil, err
	}
	if r.cookie != nil {
		req.AddCookie(r.cookie)
	}
	req.Header.Set("User-Agent", r.useragent)

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(respbytes), nil
}
