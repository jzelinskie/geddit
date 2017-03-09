// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"fmt"
)

type Redditor struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Created float64 `json:"created_utc"`
	Gold    bool    `json:"is_gold"`
	Mod     bool    `json:"is_mod"`
	Mail    bool    `json:"has_mail"`
	ModMail bool    `json:"has_mod_mail"`
	Karma
}

type Preferences struct {
	Research               bool   `json:"research"`
	ShowStylesheets        bool   `json:"show_stylesheets"`
	ShowLinkFlair          bool   `json:"show_link_flair"`
	ShowTrending           bool   `json:"show_trending"`
	PrivateFeeds           bool   `json:"private_feeds"`
	IgnoreSuggestedSort    bool   `json:"ignore_suggested_sort"`
	Media                  string `json:"media"`
	ClickGadget            bool   `json:"clickgadget"`
	LabelNSFW              bool   `json:"label_nsfw"`
	Over18                 bool   `json:"over_18"`
	EmailMessages          bool   `json:"email_messages"`
	HighlightControversial bool   `json:"highlight_controversial"`
	ForceHTTPS             bool   `json:"force_https"`
	Language               string `json:"lang"`
	HideFromRobots         bool   `json:"hide_from_robots"`
	PublicVotes            bool   `json:"public_votes"`
	ShowFlair              bool   `json:"show_flair"`
	HideAds                bool   `json:"hide_ads"`
	Beta                   bool   `json:"beta"`
	NewWindow              bool   `json:"newwindow"`
	LegacySearch           bool   `json:"legacy_search"`
}

type Friend struct {
	Date float32 `json:"date"`
	Name string  `json:"name"`
	ID   string  `json:"id"`
}

type Karma struct {
	CommentKarma int `json:"comment_karma"`
	LinkKarma    int `json:"link_karma"`
}

type Trophy struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon_70"`
}

// String returns the string representation of a reddit user.
func (r *Redditor) String() string {
	return fmt.Sprintf("%s (%d-%d)", r.Name, r.LinkKarma, r.CommentKarma)
}
