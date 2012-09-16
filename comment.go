// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"errors"
	"net/url"
	"strings"
)

type Comment struct {
	Author              string  `json:"author"`
	Body                string  `json:"body"`
	BodyHTML            string  `json:"body_html"`
	Subreddit           string  `json:"subreddit"`
	LinkId              string  `json:"link_id"`
	ParentId            string  `json:"parent_id"`
	SubredditId         string  `json:"subreddit_id"`
	FullId              string  `json:"name"`
	UpVotes             int     `json:"ups"`
	DownVotes           int     `json:"downs"`
	Created             float32 `json:"created_utc"`
	Edited              bool    `json:"edited"`
	BannedBy            *string `json:"banned_by"`
	ApprovedBy          *string `json:"approved_by"`
	AuthorFlairTxt      *string `json:"author_flair_text"`
	AuthorFlairCSSClass *string `json:"author_flair_css_class"`
	NumReports          *int    `json:"num_reports"`
	Likes               *int    `json:"likes"`
}

func (s *Session) CommentHeadline(h Headline, comment string) error {
	loc := "http://www.reddit.com/api/comment"
	vals := &url.Values{
		"thing_id": {h.FullId},
		"text":     {comment},
		"uh":       {s.Modhash},
	}
	body, err := getResponse(loc, vals, s)
	if err != nil {
		return err
	}

	if !strings.Contains(body.String(), "data") {
		return errors.New("Failed to post comment.")
	}

	return nil
}

func (s *Session) CommentReply(c Comment, comment string) error {
	loc := "http://www.reddit.com/api/comment"
	vals := &url.Values{
		"parent": {c.FullId},
		"text":   {comment},
		"uh":     {s.Modhash},
	}
	body, err := getResponse(loc, vals, s)
	if err != nil {
		return err
	}

	if !strings.Contains(body.String(), "data") {
		return errors.New("Failed to post comment.")
	}

	return nil
}
