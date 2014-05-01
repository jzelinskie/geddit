# reddit
a convenient abstraction for the [reddit.com](http://reddit.com) API

Checkout the [godoc page](http://godoc.org/github.com/jzelinskie/reddit)

## warning

this is a **work in progress**!

## example

```Go
package main

import (
	"fmt"

	"github.com/jzelinskie/reddit"
)

// Please don't handle errors this way.
func main() {
	// Login to reddit
	session, _ := reddit.NewAccountSession(
		"novelty_account",
		"password",
		"golang reddit example",
	)

	// Get reddit's default frontpage
	headlines, _ := session.DefaultFrontpage()

	// Get our own personal frontpage
	submissions, _ = session.Frontpage()

	// Upvote the first post
	session.Vote(submissions[0], reddit.UpVote)
}
```
