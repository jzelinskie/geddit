// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package reddit

import (
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
	IsDefault bool
	Username  string
	Password  string
	Cookie    *http.Cookie
	Modhash   string `json:"modhash"`
	Settings  *SessionSettings
}

// Change some session settings (unimplemented and wanting to add more) 
type SessionSettings struct {
	ClearPassAfterLogin bool
	AutoReLogin         bool
}

// Use default session settings
func DefaultSessionSettings() *SessionSettings {
	sessionSettings := new(SessionSettings)
	sessionSettings.ClearPassAfterLogin = false
	sessionSettings.AutoReLogin = false
	return sessionSettings
}

// Use a default session
func DefaultSession() *Session {
	Dsession := new(Session)
	Dsession.IsDefault = true
	Dsession.Username = ""
	Dsession.Password = ""
	Dsession.Cookie = nil
	Dsession.Settings = nil
	Dsession.Settings = DefaultSessionSettings()
	return Dsession
}

// Make a new session
func NewSession(username, password string, sessionSettings *SessionSettings) *Session {
	session := new(Session)
	session.IsDefault = false
	session.Username = username
	session.Password = password
	if sessionSettings == nil {
		session.Settings = DefaultSessionSettings()
	}

	return session
}

// String returns a string representation of a session.
func (s *Session) String() string {
	return fmt.Sprintf("%s %s %s %s", s.Username, s.Password, s.Modhash, s.Cookie)
}

// Login returns a new authenticated reddit session.
func (s *Session) Login() error {
	if s.IsDefault {
		return errors.New("Cannot login with default session")
	}
	// Make the login request.
	loginurl := fmt.Sprintf("http://www.reddit.com/api/login/%s", s.Username)
	resp, err := http.PostForm(loginurl,
		url.Values{
			"user":     {s.Username},
			"passwd":   {s.Password},
			"api_type": {"json"},
		})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
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
		return err
	}
	if len(r.Json.Errors) != 0 {
		var msg []string
		for _, k := range r.Json.Errors {
			msg = append(msg, k[1])
		}
		return errors.New(strings.Join(msg, ", "))
	}
	s.Modhash = r.Json.Data.Modhash

	return nil
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
func (s *Session) VoteHeadline(h *Headline, v string) error {
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
