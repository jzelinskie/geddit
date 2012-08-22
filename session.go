// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package reddit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Session struct {
	Username string
	password string
	modhash  string
	cookie   string
}

// NewSession returns an empty session.
func NewSession() *Session {
	return &Session{}
}

// Login authenticates the session.
func (s *Session) Login(user, pass string) error {
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

	type Response struct {
		Json struct {
			Data struct {
				Modhash string
				Cookie  string
			}
		}
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return err
	}

	s.Username = user
	s.password = pass
	s.modhash = r.Json.Data.Modhash
	s.cookie = r.Json.Data.Cookie

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
