// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	UpVote     = "1"
	DownVote   = "-1"
	RemoveVote = "0"
)

var DefaultUserAgent = "github.com/jzelinskie/reddit"

func getResponse(url string, vals *url.Values, s *Session) (*bytes.Buffer, error) {
	// Determine the HTTP action.
	var action, finalurl string
	if vals == nil {
		action = "GET"
		finalurl = url
	} else {
		action = "POST"
		finalurl = url + vals.Encode()
	}

	// Create a request and add the proper headers.
	req, err := http.NewRequest(action, finalurl, nil)
	if err != nil {
		return nil, err
	}
	if s == nil {
		req.AddCookie(s.Cookie)
	}
	req.Header.Set("User-Agent", DefaultUserAgent)

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
