// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package reddit

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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
func Login(user, pass string, forceValidCert bool) (*Session, error) {
	s := &Session{
		Username: user,
		Password: pass,
	}

	// Skip ssl certificate verification if !forceValidCert
	var client *http.Client
	if !forceValidCert {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}

	// Make the login request.
	resp, err := client.PostForm("https://www.reddit.com/api/login",
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
	//TODO, obviously
	return nil
}

// Clear clears all session cookies and updates the current session with a new one.
func (s *Session) Clear(password string) error {
	loc := "http://www.reddit.com/api/clear_sessions"
	vals := &url.Values{
		"curpass": {password},
		"uh":      {s.Modhash},
	}
	body, err := getResponse(loc, vals, s)
	if err != nil {
		return err
	}
	if !strings.Contains(body.String(), "all other sessions have been logged out") {
		return errors.New("Failed to clear session.")
	}
	return nil
}

// Me returns an up-to-date redditor object of the current user.
func (s *Session) Me() (*Redditor, error) {
	loc := "http://www.reddit.com/api/me.json"
	body, err := getResponse(loc, nil, s)
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

// VoteHeadline either votes or rescinds a vote for the given headline. The second parameter
// expects to one of the following constants: UpVote, DownVote, RemoveVote.
func (s *Session) VoteHeadline(h Headline, v string) error {
	loc := "http://www.reddit.com/api/vote"
	vals := &url.Values{
		"id":  {h.FullId},
		"dir": {v},
		"uh":  {s.Modhash},
	}
	body, err := getResponse(loc, vals, s)
	if err != nil {
		return err
	}
	if body.String() != "{}" {
		return errors.New("Failed to vote on headline.")
	}
	return nil
}
