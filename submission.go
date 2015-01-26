// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Submission represents an individual post from the perspective
// of a subreddit. Remember to check for nil pointers before
// using any pointer fields.
type Submission struct {
	Author       string  `json:"author"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Domain       string  `json:"domain"`
	Subreddit    string  `json:"subreddit"`
	SubredditID  string  `json:"subreddit_id"`
	FullID       string  `json:"name"`
	ID           string  `json:"id"`
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

func (h Submission) voteID() string   { return h.FullID }
func (h Submission) deleteID() string { return h.FullID }
func (h Submission) replyID() string  { return h.FullID }

// FullPermalink returns the full URL of a submission.
func (h *Submission) FullPermalink() string {
	return "http://reddit.com" + h.Permalink
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

// SubredditSubmissions returns the submissions on the given subreddit.
func (s Session) SubredditSubmissions(subreddit string) ([]*Submission, error) {
	req := request{
		url:       fmt.Sprintf("http://www.reddit.com/r/%s.json", subreddit),
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

// Returns top submmissions from a subreddit
// Second argument must have value of one of "hour", "day", "week", "month", "year", "all"
func (s Session) TopSubmissions(subreddit string, time string) ([]*Submission, error) {
	times := []string{"hour", "day", "week", "month", "year", "all"}
	url := ""
	timeValid := ""

	//check if time is valid value
	for _, k := range times {
		if time == k {
			timeValid = time
			url = fmt.Sprintf("http://www.reddit.com/r/%s/top.json?t=%s", subreddit, timeValid)
		}
	}
	if timeValid == "" {
		return nil, errors.New("value must be one of (hour, day, week, month, year, all)")
	}
	req := request{
		url:       url,
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

// Returns a slice of submissions from the hot section of the subreddit
func (s Session) HotSubmissions(subreddit string) ([]*Submission, error) {
	req := request{
		url:       fmt.Sprintf("http://www.reddit.com/r/%s/hot.json", subreddit),
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

// SortedSubmissions will return submissions from a subreddit (or homepage if "") by popularity and age
// TODO Review this
func (s Session) SortedSubmissions(subreddit string, popularity popularitySort, age ageSort) ([]*Submission, error) {
	if age != DefaultAge {
		switch popularity {
		case NewSubmissions, RisingSubmissions, HotSubmissions:
			return nil, fmt.Errorf("cannot sort %s by %s", popularity, age)
		}
	}

	url := "http://reddit.com/"

	if subreddit != "" {
		url = fmt.Sprintf("%sr/%s/", url, subreddit)
	}

	if popularity != DefaultPopularity {
		if popularity == NewSubmissions || popularity == RisingSubmissions {
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

	req := &request{
		url:       url,
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}
