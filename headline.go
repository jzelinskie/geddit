// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"fmt"
)

// Headline represents an individual post from the perspective
// of a subreddit. Remember to check for nil pointers before
// using any pointer fields.
type Headline struct {
	Author       string  `json:"author"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Domain       string  `json:"domain"`
	Subreddit    string  `json:"subreddit"`
	SubredditId  string  `json:"subreddit_id"`
	FullId       string  `json:"name"`
	Id           string  `json:"id"`
	Permalink    string  `json:"permalink"`
	Selftext     string  `json:"selftext"`
	ThumbnailURL string  `json:"thumbnail"`
	DateCreated  float32 `json:"created_utc"`
	NumComments  int     `json:"num_comments"`
	Score        int     `json:"score"`
	Ups          int     `json:"ups"`
	Downs        int     `json:"downs"`
	IsNSFW       bool    `json:"over_18"`
	IsSelf       bool    `json:"is_self"`
	WasClicked   bool    `json:"clicked"`
	IsSaved      bool    `json:"saved"`
	BannedBy     *string `json:"banned_by"`
}

// FullPermalink returns the full URL of a headline.
func (h *Headline) FullPermalink() string {
	return "http://reddit.com" + h.Permalink
}

// String returns the string representation of a headline.
func (h *Headline) String() string {
	plural := ""
	if h.NumComments != 1 {
		plural = "s"
	}
	comments := fmt.Sprintf("%d comment%s", h.NumComments, plural)
	return fmt.Sprintf("%d - %s (%s)", h.Score, h.Title, comments)
}
