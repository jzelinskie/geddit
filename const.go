// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

// Vote represents the three possible states of a vote on reddit.
type Vote string

const (
	// UpVote represents a positive vote.
	UpVote Vote = "1"

	// DownVote represents a negative vote.
	DownVote = "-1"

	// RemoveVote represents no vote.
	RemoveVote = "0"
)

// Sort headlines by popularity
type PopularitySort string

const (
	DefaultPopularity      PopularitySort = ""
	HotHeadlines                          = "hot"
	NewHeadlines                          = "new"
	RisingHeadlines                       = "rising"
	TopHeadlines                          = "top"
	ControversialHeadlines                = "controversial"
)

// Sort headlines by age
type AgeSort string

const (
	DefaultAge AgeSort = ""
	ThisHour           = "hour"
	ThisMonth          = "month"
	ThisYear           = "year"
	AllTime            = "all"
)
