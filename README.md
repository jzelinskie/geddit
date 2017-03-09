# geddit

[![GoDoc](https://godoc.org/github.com/jzelinskie/geddit?status.svg)](https://godoc.org/github.com/jzelinskie/geddit)
[![Build Status](https://api.travis-ci.org/jzelinskie/geddit.svg?branch=master)](https://travis-ci.org/jzelinskie/geddit)

Geddit is a convenient abstraction for the [reddit.com](http://reddit.com) API in Go.
This library is a WIP.
It should have some API coverage, but does not yet include things like the new OAuth model.

## examples

See [godoc](http://godoc.org/github.com/jzelinskie/geddit) for OAuth examples.

Here is an example usage of the old, cookie authentication method:

(NOTE: You will be heavily rate-limited by reddit's API when using cookies. Consider switching to OAuth).

```Go
package main

import (
	"fmt"

	"github.com/jzelinskie/geddit"
)

// Please don't handle errors this way.
func main() {
	// Login to reddit
	session, _ := geddit.NewLoginSession(
		"novelty_account",
		"password",
		"gedditAgent v1",
	)

	// Set listing options
	subOpts := geddit.ListingOptions{
		Limit: 10,
	}

	// Get reddit's default frontpage
	submissions, _ := session.DefaultFrontpage(geddit.DefaultPopularity, subOpts)

	// Get our own personal frontpage
	submissions, _ = session.Frontpage(geddit.DefaultPopularity, subOpts)

	// Get specific subreddit submissions, sorted by new
	submissions, _ = session.SubredditSubmissions("hockey", geddit.NewSubmissions, subOpts)

	// Print title and author of each submission
	for _, s := range submissions {
		fmt.Printf("Title: %s\nAuthor: %s\n\n", s.Title, s.Author)
	}

	// Upvote the first post
	session.Vote(submissions[0], geddit.UpVote)
}
```
