// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Session represents an HTTP session with reddit.com -- all authenticated API
// calls are methods bound to this type.
type Session struct {
	Username string
	Password string
	Cookie   *http.Cookie
	Modhash  string `json:"modhash"`
}

// String returns a string representation of a session.
func (s *Session) String() string {
	return fmt.Sprintf("%s %s %s %s", s.Username, s.Password, s.Modhash, s.Cookie)
}

// Login returns a new authenticated reddit session.
func Login(user, pass string) (*Session, error) {
	s := &Session{
		Username: user,
		Password: pass,
	}

	// Make the login request.
	resp, err := http.PostForm("http://www.reddit.com/api/login",
		url.Values{
			"user":     {user},
			"passwd":   {pass},
			"api_type": {"json"},
		})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	// Get the session cookie.
	for _, c := range resp.Cookies() {
		if c.Name == "reddit_session" {
			s.Cookie = c
		}
	}

	// Get the modhash from the JSON.
	type Response struct {
		Json struct {
			Errors [][]string
			Data   struct {
				Modhash string
			}
		}
	}
	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	if len(r.Json.Errors) != 0 {
		var msg []string
		for _, k := range r.Json.Errors {
			msg = append(msg, k[1])
		}
		return nil, errors.New(strings.Join(msg, ", "))
	}
	s.Modhash = r.Json.Data.Modhash

	return s, nil
}

// Logout terminates the authentication of the session.
func (s *Session) Logout() error {
	return nil
}

// Clear clears all session cookies and updates the current session with a new one.
func (s *Session) Clear(password string) error {
	formstring := url.Values{
		"curpass": {password},
		"uh":      {s.Modhash},
	}.Encode()
	req, err := http.NewRequest("POST", "http://www.reddit.com/api/clear_sessions?"+formstring, nil)
	req.AddCookie(s.Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

// Me returns an up-to-date redditor object of the current user.
func (s *Session) Me() (*Redditor, error) {
	req, err := http.NewRequest("GET", "http://www.reddit.com/api/me.json", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(s.Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	type Response struct {
		Data Redditor
	}
	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}

// VoteHeadline either votes or rescinds a vote for the given headline. The second parameter
// expects to one of the following constants: UpVote, DownVote, RemoveVote.
func (s *Session) VoteHeadline(h Headline, v string) error {
	formstring := url.Values{
		"id":  {h.FullId},
		"dir": {v},
		"uh":  {s.Modhash},
	}.Encode()
	req, err := http.NewRequest("POST", "http://www.reddit.com/api/vote?"+formstring, nil)
	req.AddCookie(s.Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	if contents, err := ioutil.ReadAll(resp.Body); err != nil || string(contents) != "{}" {
		return errors.New("Failed to vote on headline.")
	}

	return nil
}
