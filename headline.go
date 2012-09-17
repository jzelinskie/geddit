// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"fmt"
)

// Headline represents an individual post from the perspective
// of a subreddit.
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
	// Ups and downs are fake to trick spammers
	Ups        int     `json:"ups"`
	Downs      int     `json:"downs"`
	IsNSFW     bool    `json:"over_18"`
	IsSelf     bool    `json:"is_self"`
	WasClicked bool    `json:"clicked"`
	IsSaved    bool    `json:"saved"`
	BannedBy   *string `json:"banned_by"`
}

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

type Headlines []*Headline

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

// DefaultHeadlines returns a slice of headlines on the default reddit frontpage.
func DefaultHeadlines() (Headlines, error) {
	url := "http://www.reddit.com/.json"
	body, err := getResponse(url, nil, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make(Headlines, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}

// SubredditHeadlines returns a slice of headlines on the given subreddit.
func SubredditHeadlines(subreddit string) (Headlines, error) {
	url := fmt.Sprintf("http://www.reddit.com/r/%s.json", subreddit)
	body, err := getResponse(url, nil, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make(Headlines, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}

// SortedHeadlines will return headlines from a subreddit (or homepage if "") by popularity and age
func SortedHeadlines(subreddit string, popularity PopularitySort, age AgeSort) (Headlines, error) {
	if age != DefaultAge {
		switch popularity {
		case NewHeadlines, RisingHeadlines, HotHeadlines:
			return nil, fmt.Errorf("Cannot sort %s by %s", popularity, age)
		}
	}

	url := "http://reddit.com/"

	if subreddit != "" {
		url = fmt.Sprintf("http://%s.reddit.com/", subreddit)
	}

	if popularity != DefaultPopularity {
		if popularity == NewHeadlines || popularity == RisingHeadlines {
			url = fmt.Sprintf("%snew.json?sort=%s", url, popularity)
		} else {
			url = fmt.Sprintf("%s%s.json?sort=%s", url, popularity, popularity)
		}
	} else {
		url = fmt.Sprintf("%s.json", url)
	}

	if age != DefaultAge {
		if popularity != DefaultPopularity {
			url = fmt.Sprintf("%s&t=%s", url, age)
		} else {
			url = fmt.Sprintf("%s?t=%s", url, age)
		}
	}

	body, err := getResponse(url, nil, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make(Headlines, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil

}
