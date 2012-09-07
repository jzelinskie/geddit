// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Headline struct {
	Author       string  `json:"author"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Domain       string  `json:"domain"`
	Subreddit    string  `json:"subreddit"`
	SubredditId  string  `json:"subreddit_id"`
	Name         string  `json:"name"`
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

// FullPermalink returns the full URL of a headline.
func (h Headline) FullPermalink() string {
	return "http://reddit.com" + h.Permalink
}

// String returns the string representation of a headline.
func (h Headline) String() string {
	var comments string
	switch h.NumComments {
	case 0:
		comments = "0 comments"
	case 1:
		comments = "1 comment"
	default:
		comments = fmt.Sprintf("%d comments", h.NumComments)
	}
	return fmt.Sprintf("%d - %s (%s)", h.Score, h.Title, comments)
}

// DefaultHeadlines returns a slice of headlines on the default reddit frontpage.
func DefaultHeadlines() ([]Headline, error) {
	url := "http://www.reddit.com/.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make([]Headline, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}

// SubredditHeadlines returns a slice of headlines on the given subreddit.
func SubredditHeadlines(subreddit string) ([]Headline, error) {
	url := fmt.Sprintf("http://www.reddit.com/r/%s.json", subreddit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make([]Headline, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}
