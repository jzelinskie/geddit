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

func main() {
	// Login to reddit
	session, err := reddit.NewAccountSession("novelty_account", "password", "golang reddit example")
	if err != nil {
		panic(err)
	}

	// Get our personal frontpage
	headlines, err := session.DefaultHeadlines()
	if err != nil {
		panic(err)
	}

	// Get the default frontpage (not logged in)
	headlines, err = session.AnonymousSession.DefaultHeadlines()
	if err != nil {
		panic(err)
	}

	// Upvote the first post
	err = session.VoteHeadline(headlines[0], reddit.UpVote)
	if err != nil {
		panic(err)
	}
}
```
