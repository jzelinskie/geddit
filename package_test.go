// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	//"fmt"
	"testing"
)

func TestPublicAPI(t *testing.T) {
	sub, err := AboutSubreddit("golang")
	if err != nil {
		t.Fatal(err)
	}
	if sub.FullId != "t5_2rc7j" {
		t.Fatalf("Subreddit id: %s, Expected: t5_2rc7j", sub.FullId)
	}

	_, err = DefaultHeadlines()
	if err != nil {
		t.Fatal(err)
	}
	_, err = SubredditHeadlines("golang")
	if err != nil {
		t.Fatal(err)
	}

	me, err := AboutRedditor("jzelinskie")
	if err != nil {
		t.Fatal(err)
	}
	if me.Id != "5pi0h" {
		t.Fatalf("Redditor id: %s, Expected: 5pi0h", me.Id)
	}
}

func TestAuthenticatedAPI(t *testing.T) {
	// Disable ssl certificate verification
	sesh, err := Login("goreddittest", "test", false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sesh.Me()
	if err != nil {
		t.Fatal(err)
	}

	hl, err := DefaultHeadlines()
	if err != nil {
		t.Fatal(err)
	}
	err = sesh.VoteHeadline(hl[0], UpVote)
	if err != nil {
		t.Fatal(err)
	}
	err = sesh.Clear("test")
	if err != nil {
		t.Fatal(err)
	}
}
