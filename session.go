// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
)

// Session represents an HTTP session with reddit.com
// without logging into an account.
type Session struct {
	useragent string
}

// NewSession creates a new unauthenticated session to reddit.com.
func NewSession(useragent string) *Session {
	return &Session{
		useragent: useragent,
	}
}

// DefaultFrontpage returns the submissions on the default reddit frontpage.
func (s Session) DefaultFrontpage() ([]*Submission, error) {
	req := request{
		url:       "http://www.reddit.com/.json",
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
	r := &Response{}
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
	t := []string{"hour", "day", "week", "month", "year", "all"}
	url := ""
	timeValid := ""

	//check if time is valid value
	for _, k := range t {
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

// AboutRedditor returns a Redditor for the given username.
func (s Session) AboutRedditor(username string) (*Redditor, error) {
	req := &request{
		url:       fmt.Sprintf("http://www.reddit.com/user/%s/about.json", username),
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data Redditor
	}
	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

// AboutSubreddit returns a subreddit for the given subreddit name.
func (s Session) AboutSubreddit(subreddit string) (*Subreddit, error) {
	req := &request{
		url:       fmt.Sprintf("http://www.reddit.com/r/%s/about.json", subreddit),
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data Subreddit
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

// Comments returns the comments for a given Submission.
func (s Session) Comments(h *Submission) ([]*Comment, error) {
	req := &request{
		url:       fmt.Sprintf("http://www.reddit.com/comments/%s/.json", h.ID),
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	r := json.NewDecoder(body)
	var interf interface{}
	if err = r.Decode(&interf); err != nil {
		return nil, err
	}
	helper := new(helper)
	helper.buildComments(interf)

	return helper.comments, nil
}

// CaptchaImage gets the png corresponding to the captcha iden and decodes it
func (s Session) CaptchaImage(iden string) (image.Image, error) {
	req := &request{
		url:       fmt.Sprintf("http://www.reddit.com/captcha/%s", iden),
		useragent: s.useragent,
	}

	p, err := req.getResponse()

	if err != nil {
		return nil, err
	}

	m, err := png.Decode(p)

	if err != nil {
		return nil, err
	}

	return m, nil
}
