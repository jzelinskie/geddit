// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
func TestAboutSubreddit(t *testing.T) {
	sesh := NewSession()
	_, err := sesh.AboutSubreddit("golang")
	if err != nil {
		panic(err)
	}
}

func TestFrontPage(t *testing.T) {
	sesh := NewSession()
	_, err := sesh.FrontPage()
	if err != nil {
		panic(err)
	}
}

func TestSubreddit(t *testing.T) {
	sesh := NewSession()
	_, err := sesh.Subreddit("golang")
	if err != nil {
		panic(err)
	}
}

func TestAboutRedditor(t *testing.T) {
	sesh := NewSession()
	r, err := sesh.AboutRedditor("jzelinskie")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
*/

// Test this manually since spamming login can get you banned
func TestLogin(t *testing.T) {
	sesh := NewSession()
	err := sesh.Login("goreddittest", "test")
	if err != nil {
		fmt.Println(err)
	}

	r, err := http.NewRequest("GET", "http://www.reddit.com/api/me.json", nil)
	if err != nil {
		fmt.Println(err)
	}
	r.AddCookie(sesh.Cookie)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
