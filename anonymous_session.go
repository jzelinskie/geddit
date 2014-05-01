// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"fmt"
)

// AnonymousSession represents an HTTP session with reddit.com
// without logging into an account.
type AnonymousSession struct {
	useragent string
}

// NewAnonymousSession creates a new unauthenticated session to reddit.com.
func NewAnonymousSession(useragent string) *AnonymousSession {
	return &AnonymousSession{
		useragent: useragent,
	}
}

// DefaultFrontpage returns the submissions on the default reddit frontpage.
func (s AnonymousSession) DefaultFrontpage() ([]*Submission, error) {
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
func (s AnonymousSession) SubredditSubmissions(subreddit string) ([]*Submission, error) {
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

// SortedSubmissions will return submissions from a subreddit (or homepage if "") by popularity and age
// TODO Review this
func (s AnonymousSession) SortedSubmissions(subreddit string, popularity popularitySort, age ageSort) ([]*Submission, error) {
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
func (s AnonymousSession) AboutRedditor(username string) (*Redditor, error) {
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
func (s AnonymousSession) AboutSubreddit(subreddit string) (*Subreddit, error) {
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
func (s AnonymousSession) Comments(h *Submission) ([]*Comment, error) {
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
