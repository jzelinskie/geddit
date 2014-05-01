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

	// Get reddit's default frontpage
	headlines, err := session.DefaultFrontpage()
	if err != nil {
		panic(err)
	}

	// Get our own personal frontpage
	headlines, err = session.Frontpage()
	if err != nil {
		panic(err)
	}

	// Upvote the first post
	err = session.Vote(headlines[0], reddit.UpVote)
	if err != nil {
		panic(err)
	}
}
```
