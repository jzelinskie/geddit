// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Subreddit struct {
	Name        string  `json:"display_name"`
	Title       string  `json:"title"`
	Desc        string  `json:"description"`
	PublicDesc  string  `json:"public_description"`
	URL         string  `json:"url"`
	Id          string  `json:"id"`
	HeaderImg   string  `json:"header_img"`
	DateCreated float32 `json:"created_utc"`
	NumSubs     int     `json:"subscribers"`
	IsNSFW      bool    `json:"over18"`
}

// String returns the string representation of a subreddit.
func (s *Subreddit) String() string {
	var subs string
	switch s.NumSubs {
	case 1:
		subs = "1 subscriber"
	default:
		subs = fmt.Sprintf("%d subscribers", s.NumSubs)
	}
	return fmt.Sprintf("%s (%s)", s.Title, subs)
}

// AboutSubreddit returns a subreddit for the given subreddit name.
func AboutSubreddit(subreddit string) (*Subreddit, error) {
	url := fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type Response struct {
		Data Subreddit
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}
