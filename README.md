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
		"gedditAgent v1",
	)

	// Get reddit's default frontpage
	submissions, _ := session.DefaultFrontpage()

	// Get our own personal frontpage
	submissions, _ = session.Frontpage()

	// Print title and author of each submission
	for _, s := range submissions {
		fmt.Printf("Title: %s\nAuthor: %s\n\n", s.Title, s.Author)
	}

	// Upvote the first post
	session.Vote(submissions[0], geddit.UpVote)
}
```
