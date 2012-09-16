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
