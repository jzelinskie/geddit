// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package geddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// OAuthSession represents an OAuth session with reddit.com --
// all authenticated API calls are methods bound to this type.
type OAuthSession struct {
	Client       *http.Client
	ClientID     string
	ClientSecret string
	OAuthConfig  *oauth2.Config
	UserAgent    string
}

// NewLoginSession creates a new session for those who want to log into a
// reddit account via OAuth.
func NewOAuthSession(clientID, clientSecret, useragent string, limit bool) (*OAuthSession, error) {
	s := &OAuthSession{}

	if useragent != "" {
		s.UserAgent = useragent
	} else {
		s.UserAgent = "Geddit Reddit Bot https://github.com/jzelinskie/geddit"
	}

	// Set OAuth config
	// TODO Set user-defined scopes
	s.OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{}, //"identity", "read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.reddit.com",
			TokenURL: "https://www.reddit.com/api/v1/access_token",
		},
	}

	ctx := context.Background()

	// TODO offer auth code version as well as personal scripts
	t, err := s.OAuthConfig.PasswordCredentialsToken(ctx, s.Username, s.Password)
	if err != nil {
		return nil, err
	}

	s.Client = s.OAuthConfig.Client(ctx, t)
	return s, nil
}

func (s OAuthSession) getBody(url string, d interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// This is needed to avoid rate limits
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(req.Header.Get("User-Agent"))
	defer resp.Body.Close()

	// DEBUG
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("***DEBUG***\nRequest Body: %s\n***DEBUG***\n\n", body)

	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}

	return nil
}
