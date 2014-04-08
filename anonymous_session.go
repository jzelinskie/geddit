// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"encoding/json"
	"fmt"
)

// Session represents an HTTP session with reddit.com without logged into
// an account.
type AnonymousSession struct {
	useragent string
}

func NewAnonymousSession(useragent string) *AnonymousSession {
	return &AnonymousSession{
		useragent: useragent,
	}
}

// DefaultHeadlines returns a slice of headlines on the default reddit frontpage.
// TODO Override this method for AccountSession
func (s AnonymousSession) DefaultHeadlines() ([]*Headline, error) {
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
				Data *Headline
			}
		}
	}
	r := &Response{}
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make([]*Headline, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}

// SubredditHeadlines returns a slice of headlines on the given subreddit.
func (s AnonymousSession) SubredditHeadlines(subreddit string) ([]*Headline, error) {
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
				Data *Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make([]*Headline, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
}

// SortedHeadlines will return headlines from a subreddit (or homepage if "") by popularity and age
// TODO Review this
func (s AnonymousSession) SortedHeadlines(subreddit string, popularity PopularitySort, age AgeSort) ([]*Headline, error) {
	if age != DefaultAge {
		switch popularity {
		case NewHeadlines, RisingHeadlines, HotHeadlines:
			return nil, fmt.Errorf("Cannot sort %s by %s", popularity, age)
		}
	}

	url := "http://reddit.com/"

	if subreddit != "" {
		url = fmt.Sprintf("%sr/%s/", url, subreddit)
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
				Data *Headline
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	headlines := make([]*Headline, len(r.Data.Children))
	for i, child := range r.Data.Children {
		headlines[i] = child.Data
	}

	return headlines, nil
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

func (s AnonymousSession) HeadlineComments(h *Headline) ([]*Comment, error) {
	req := &request{
		url:       fmt.Sprintf("http://www.reddit.com/comments/%s/.json", h.Id),
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
