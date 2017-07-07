// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"

	"github.com/google/go-querystring/query"
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
func (s Session) DefaultFrontpage(sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	return s.SubredditSubmissions("", sort, params)
}

// SubredditSubmissions returns the submissions on the given subreddit.
func (s Session) SubredditSubmissions(subreddit string, sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	baseUrl := "https://www.reddit.com"

	// If subbreddit given, add to URL
	if subreddit != "" {
		baseUrl += "/r/" + subreddit
	}

	redditUrl := fmt.Sprintf(baseUrl+"/%s.json?%s", sort, v.Encode())

	req := request{
		url:       redditUrl,
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
		url:       fmt.Sprintf("https://www.reddit.com/user/%s/about.json", username),
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
		url:       fmt.Sprintf("https://www.reddit.com/r/%s/about.json", subreddit),
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
		url:       fmt.Sprintf("https://www.reddit.com/comments/%s/.json", h.ID),
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
		url:       fmt.Sprintf("https://www.reddit.com/captcha/%s", iden),
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

// SubredditComments gets all the new comments from a subreddit, returning them in a slice of Comment structs
func (s Session) SubredditComments(subreddit string) ([]*Comment, error) {
	baseURL := "https://www.reddit.com"
	
	if subreddit != "" {
		baseURL += "/r/" + subreddit
	}
	
	subCommentsURL := baseURL + "/comments.json"

	req := request{
		url:       subCommentsURL,
		useragent: s.useragent,
	}

	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}
	
	var comments interface{}
	if err = json.NewDecoder(body).Decode(&comments); err != nil {
		return nil, err
	}
	helper := new(helper)
	helper.buildComments(comments)

	return helper.comments, nil
}
