// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

// Types of votes for headlines
type Vote string

const (
	UpVote     Vote = "1"
	DownVote        = "-1"
	RemoveVote      = "0"
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
