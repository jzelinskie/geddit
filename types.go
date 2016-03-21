// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

// vote represents the three possible states of a vote on reddit.
type Vote string

const (
	UpVote     Vote = "1"
	DownVote        = "-1"
	RemoveVote      = "0"
)

type Captcha struct {
	Iden     string
	Response string
}

// NewSubmission contains the data needed to submit
type NewSubmission struct {
	Subreddit   string
	Title       string
	Content     string
	Self        bool
	SendReplies bool
	Resubmit    bool
	Save        bool
	Captcha     *Captcha
}

// NewLinkSubmission returns a NewSubmission with parameters appropriate for a link submission
func NewLinkSubmission(sr, title, link string, replies bool, c *Captcha) *NewSubmission {
	return &NewSubmission{sr, title, link, false, replies, true, true, c}
}

// NewTextSubmission returns a NewSubmission with parameters appropriate for a text submission
func NewTextSubmission(sr, title, text string, replies bool, c *Captcha) *NewSubmission {
	return &NewSubmission{sr, title, text, true, replies, true, true, c}
}

// PopularitySort represents the possible ways to sort submissions by popularity.
type PopularitySort string

const (
	DefaultPopularity        PopularitySort = ""
	HotSubmissions                          = "hot"
	NewSubmissions                          = "new"
	RisingSubmissions                       = "rising"
	TopSubmissions                          = "top"
	ControversialSubmissions                = "controversial"
)

// ageSort represents the possible ways to sort submissions by age.
type ageSort string

const (
	DefaultAge ageSort = ""
	ThisHour           = "hour"
	ThisDay            = "day"
	ThisMonth          = "month"
	ThisYear           = "year"
	AllTime            = "all"
)

type ListingOptions struct {
	Time    string `url:"t,omitempty"`
	Limit   int    `url:"limit,omitempty"`
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	Count   int    `url:"count,omitempty"`
	Show    string `url:"show,omitempty"`
	Article string `url:"article,omitempty"`
}

// Voter represents something that can be voted on reddit.com.
type Voter interface {
	voteID() string
}

// Deleter represents something that can be deleted on reddit.com.
type Deleter interface {
	deleteID() string
}

// Replier represents something that can be replied to on reddit.com.
type Replier interface {
	replyID() string
}
