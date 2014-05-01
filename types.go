// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

// vote represents the three possible states of a vote on reddit.
type vote string

const (
	UpVote     vote = "1"
	DownVote        = "-1"
	RemoveVote      = "0"
)

// popularitySort represents the possible ways to sort headlines by popularity.
type popularitySort string

const (
	DefaultPopularity      popularitySort = ""
	HotHeadlines                          = "hot"
	NewHeadlines                          = "new"
	RisingHeadlines                       = "rising"
	TopHeadlines                          = "top"
	ControversialHeadlines                = "controversial"
)

// ageSort represents the possible ways to sort headlines by age.
type ageSort string

const (
	DefaultAge ageSort = ""
	ThisHour           = "hour"
	ThisMonth          = "month"
	ThisYear           = "year"
	AllTime            = "all"
)

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
