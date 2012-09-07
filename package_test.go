// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	//"fmt"
	"testing"
)

func TestPublicAPI(t *testing.T) {
	_, err := AboutSubreddit("golang")
	if err != nil {
		panic(err)
	}
	_, err = DefaultHeadlines()
	if err != nil {
		panic(err)
	}

	_, err = SubredditHeadlines("golang")
	if err != nil {
		panic(err)
	}

	_, err = AboutRedditor("jzelinskie")
	if err != nil {
		panic(err)
	}
}

func TestAuthenticatedAPI(t *testing.T) {
	sesh, err := Login("goreddittest", "test")
	if err != nil {
		panic(err)
	}

	_, err = sesh.Me()
	if err != nil {
		panic(err)
	}
}
