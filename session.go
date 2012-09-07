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

type Session struct {
	Username string
	Password string
	Cookie   *http.Cookie
	Modhash  string `json:"modhash"`
}

// NewSession returns an empty session.
func NewSession() *Session {
	return &Session{}
}

// String returns a string representation of a session.
func (s Session) String() string {
	return fmt.Sprintf("%s %s %s %s", s.Username, s.Password, s.Modhash, s.Cookie)
}

// Login authenticates the session.
func (s *Session) Login(user, pass string) error {
	s.Username = user
	s.Password = pass

	resp, err := http.PostForm("http://www.reddit.com/api/login",
		url.Values{
			"user":     {user},
			"passwd":   {pass},
			"api_type": {"json"},
		})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "reddit_session" {
			s.Cookie = c
		}
	}

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
	s.Modhash = r.Json.Data.Modhash

	if len(r.Json.Errors) != 0 {
		var msg []string
		for _, k := range r.Json.Errors {
			msg = append(msg, k[1])
		}
		return errors.New(strings.Join(msg, ", "))
	}

	return nil
}

// Logout terminates the authentication of the session.
func (s *Session) Logout() error {
	return nil
}

// Clear clears all session cookies and updates the current session with a new one.
func (s *Session) Clear() error {
	// POST /api/clear_sessions
	return nil
}
