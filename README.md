# geddit
a convenient abstraction for the [reddit.com](http://reddit.com) API in Go

Checkout the [godoc page](http://godoc.org/github.com/jzelinskie/geddit)

## warning

this is a **work in progress**!

## example

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
		"golang reddit example",
	)

	// Get reddit's default frontpage
	submissions, _ := session.DefaultFrontpage()

	// Get our own personal frontpage
	submissions, _ = session.Frontpage()

	// Upvote the first post
	session.Vote(submissions[0], geddit.UpVote)
}
```
