// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"fmt"
)

// Submission represents an individual post from the perspective
// of a subreddit. Remember to check for nil pointers before
// using any pointer fields.
type Submission struct {
	Author        string  `json:"author"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Domain        string  `json:"domain"`
	Subreddit     string  `json:"subreddit"`
	SubredditID   string  `json:"subreddit_id"`
	FullID        string  `json:"name"`
	ID            string  `json:"id"`
	Permalink     string  `json:"permalink"`
	Selftext      string  `json:"selftext"`
	SelftextHTML  string  `json:"selftext_html"`
	ThumbnailURL  string  `json:"thumbnail"`
	DateCreated   float64 `json:"created_utc"`
	NumComments   int     `json:"num_comments"`
	Score         int     `json:"score"`
	Ups           int     `json:"ups"`
	Downs         int     `json:"downs"`
	IsNSFW        bool    `json:"over_18"`
	IsSelf        bool    `json:"is_self"`
	WasClicked    bool    `json:"clicked"`
	IsSaved       bool    `json:"saved"`
	BannedBy      *string `json:"banned_by"`
	LinkFlairText string  `json:"link_flair_text"`
}

func (h Submission) voteID() string   { return h.FullID }
func (h Submission) deleteID() string { return h.FullID }
func (h Submission) replyID() string  { return h.FullID }

// FullPermalink returns the full URL of a submission.
func (h *Submission) FullPermalink() string {
	return "https://reddit.com" + h.Permalink
}

// String returns the string representation of a submission.
func (h *Submission) String() string {
	plural := ""
	if h.NumComments != 1 {
		plural = "s"
	}
	comments := fmt.Sprintf("%d comment%s", h.NumComments, plural)
	return fmt.Sprintf("%d - %s (%s)", h.Score, h.Title, comments)
}
