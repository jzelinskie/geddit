// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"fmt"
	"testing"
)

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

/*
// Test this manually since spamming login can get you banned
func TestLogin(t *testing.T) {
	sesh := NewSession()
	err := sesh.Login("user", "pass")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sesh)
}
*/
