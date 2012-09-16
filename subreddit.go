// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"fmt"
)

type Subreddit struct {
	Name        string  `json:"display_name"`
	Title       string  `json:"title"`
	Desc        string  `json:"description"`
	PublicDesc  string  `json:"public_description"`
	URL         string  `json:"url"`
	FullId      string  `json:"name"`
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
	loc := fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit)
	body, err := getResponse(loc, nil, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data Subreddit
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}
