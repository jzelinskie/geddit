// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package geddit implements an abstraction for the reddit.com API.
package geddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

// LoginSession represents an HTTP session with reddit.com --
// all authenticated API calls are methods bound to this type.
type LoginSession struct {
	username  string
	password  string
	useragent string
	cookie    *http.Cookie
	modhash   string `json:"modhash"`
	Session
}

// NewLoginSession creates a new session for those who want to log into a
// reddit account.
func NewLoginSession(username, password, useragent string) (*LoginSession, error) {
	session := &LoginSession{
		username:  username,
		password:  password,
		useragent: useragent,
		Session:   Session{useragent},
	}

	loginURL := fmt.Sprintf("https://www.reddit.com/api/login/%s", username)
	postValues := url.Values{
		"user":     {username},
		"passwd":   {password},
		"api_type": {"json"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", useragent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	// Get the session cookie.
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "reddit_session" {
			session.cookie = cookie
		}
	}

	// Get the modhash from the JSON.
	type Response struct {
		JSON struct {
			Errors [][]string
			Data   struct {
				Modhash string
			}
		}
	}

	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if len(r.JSON.Errors) != 0 {
		var msg []string
		for _, k := range r.JSON.Errors {
			msg = append(msg, k[1])
		}
		return nil, errors.New(strings.Join(msg, ", "))
	}
	session.modhash = r.JSON.Data.Modhash

	return session, nil
}

// Clear clears all session cookies and updates the current session with a new one.
func (s LoginSession) Clear() error {
	req := &request{
		url: "https://www.reddit.com/api/clear_sessions",
		values: &url.Values{
			"curpass": {s.password},
			"uh":      {s.modhash},
		},
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return err
	}

	if !strings.Contains(body.String(), "all other sessions have been logged out") {
		return errors.New("failed to clear session")
	}
	return nil
}

// Frontpage returns the submissions on the logged-in user's personal frontpage.
func (s LoginSession) Frontpage(sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	var redditUrl string

	if sort == DefaultPopularity {
		redditUrl = fmt.Sprintf("https://www.reddit.com/.json?%s", v.Encode())
	} else {
		redditUrl = fmt.Sprintf("https://www.reddit.com/%s/.json?%s", sort, v.Encode())
	}

	req := request{
		url:       redditUrl,
		cookie:    s.cookie,
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

// Me returns an up-to-date redditor object of the logged-in user.
func (s LoginSession) Me() (*Redditor, error) {
	req := &request{
		url:       "https://www.reddit.com/api/me.json",
		cookie:    s.cookie,
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data Redditor
	}
	r := &Response{}
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

func (s LoginSession) Submit(ns *NewSubmission) error {

	var kind string

	if ns.Self {
		kind = "self"
	} else {
		kind = "link"
	}

	req := &request{
		url: "https://www.reddit.com/api/submit",
		values: &url.Values{
			"title":       {ns.Title},
			"url":         {ns.Content},
			"text":        {ns.Content},
			"sr":          {ns.Subreddit},
			"kind":        {kind},
			"sendreplies": {strconv.FormatBool(ns.SendReplies)},
			"resubmit":    {strconv.FormatBool(ns.Resubmit)},
			"extension":   {"json"},
			"captcha":     {ns.Captcha.Response},
			"iden":        {ns.Captcha.Iden},
			"uh":          {s.modhash},
		},
		cookie:    s.cookie,
		useragent: s.useragent,
	}

	body, err := req.getResponse()
	if err != nil {
		return err
	}
	if strings.Contains(body.String(), "error") {
		return errors.New("failed to submit")
	}
	return nil
}

// Vote either votes or rescinds a vote for a Submission or Comment.
func (s LoginSession) Vote(v Voter, vote Vote) error {
	req := &request{
		url: "https://www.reddit.com/api/vote",
		values: &url.Values{
			"id":  {v.voteID()},
			"dir": {string(vote)},
			"uh":  {s.modhash},
		},
		cookie:    s.cookie,
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return err
	}
	if body.String() != "{}" {
		return errors.New("failed to vote")
	}
	return nil
}

// Reply posts a comment as a response to a Submission or Comment.
func (s LoginSession) Reply(r Replier, comment string) error {
	req := &request{
		url: "https://www.reddit.com/api/comment",
		values: &url.Values{
			"thing_id": {r.replyID()},
			"text":     {comment},
			"uh":       {s.modhash},
		},
		cookie:    s.cookie,
		useragent: s.useragent,
	}

	body, err := req.getResponse()
	if err != nil {
		return err
	}

	if !strings.Contains(body.String(), "data") {
		return errors.New("failed to post comment")
	}

	return nil
}

// Delete deletes a Submission or Comment.
func (s LoginSession) Delete(d Deleter) error {
	req := &request{
		url: "https://www.reddit.com/api/del",
		values: &url.Values{
			"id": {d.deleteID()},
			"uh": {s.modhash},
		},
		cookie:    s.cookie,
		useragent: s.useragent,
	}

	body, err := req.getResponse()
	if err != nil {
		return err
	}

	if !strings.Contains(body.String(), "data") {
		return errors.New("failed to delete item")
	}

	return nil
}

// NeedsCaptcha returns true if captcha is required, false if it isn't
func (s LoginSession) NeedsCaptcha() (bool, error) {
	req := &request{
		url:       "https://www.reddit.com/api/needs_captcha.json",
		cookie:    s.cookie,
		useragent: s.useragent,
	}

	body, err := req.getResponse()

	if err != nil {
		return false, err
	}

	need, err := strconv.ParseBool(body.String())

	if err != nil {
		return false, err
	}

	return need, nil
}

// NewCaptchaIden gets a new captcha iden from reddit
func (s LoginSession) NewCaptchaIden() (string, error) {
	req := &request{
		url: "https://www.reddit.com/api/new_captcha",
		values: &url.Values{
			"api_type": {"json"},
		},
		cookie:    s.cookie,
		useragent: s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return "", err
	}

	// Get the CAPTCHA iden from the JSON.
	type Response struct {
		JSON struct {
			Errors [][]string
			Data   struct {
				Iden string
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return "", err
	}

	return r.JSON.Data.Iden, nil
}

// Listing returns a listing for an user
func (s LoginSession) Listing(username, listing string, sort PopularitySort, after string) ([]*Submission, error) {
	values := &url.Values{}
	if sort != "" {
		values.Set("sort", string(sort))
	}
	if after != "" {
		values.Set("after", after)
	}
	url := fmt.Sprintf("https://www.reddit.com/user/%s/%s.json?%s", username, listing, values.Encode())
	req := &request{
		url:       url,
		cookie:    s.cookie,
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

// Fetch the Overview listing for the logged-in user
func (s LoginSession) MyOverview(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "overview", sort, after)
}

// Fetch the Submitted listing for the logged-in user
func (s LoginSession) MySubmitted(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "submitted", sort, after)
}

// Fetch the Comments listing for the logged-in user
func (s LoginSession) MyComments(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "comments", sort, after)
}

// Fetch the Liked listing for the logged-in user
func (s LoginSession) MyLiked(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "liked", sort, after)
}

// Fetch the Disliked listing for the logged-in user
func (s LoginSession) MyDisliked(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "disliked", sort, after)
}

// Fetch the Hidden listing for the logged-in user
func (s LoginSession) MyHidden(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "hidden", sort, after)
}

// Fetch the Saved listing for the logged-in user
func (s LoginSession) MySaved(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "saved", sort, after)
}

// Fetch the Gilded listing for the logged-in user
func (s LoginSession) MyGilded(sort PopularitySort, after string) ([]*Submission, error) {
	return s.Listing(s.username, "gilded", sort, after)
}
