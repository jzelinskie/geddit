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

type Redditor struct {
	Name         string  `json:"name"`
	LinkKarma    int     `json:"link_karma"`
	CommentKarma int     `json:"comment_karma"`
	Created      float32 `json:"created_utc"`
	Gold         bool    `json:"is_gold"`
	Mod          bool    `json:"is_mod"`
	Mail         *bool   `json:"has_mail"`
	ModMail      *bool   `json:"has_mod_mail"`
}

// String returns the string representation of a reddit user.
func (r Redditor) String() string {
	return fmt.Sprintf("%s (%d-%d)", r.Name, r.LinkKarma, r.CommentKarma)
}

// AboutRedditor returns a redditor for the given username.
func AboutRedditor(username string) (*Redditor, error) {
	url := fmt.Sprintf("http://www.reddit.com/user/%s/about.json", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type Response struct {
		Data Redditor
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}
