#reddit
an abstraction for the [reddit.com](http://reddit.com) API

## warning

this is a **work in progress**! future features include

* more public api calls
* auth'd api calls
* comments

## example

```Go
package main

import (
 "fmt"
 "github.com/jzelinskie/reddit"
)

func main() {
  // Login to reddit
  session, err := reddit.Login("novelty_account", "password")
  if err != nil {
    panic(err)
  }

  // Get default frontpage (not our personal one)
  headlines, err := reddit.DefaultHeadlines()
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

## docs

Checkout the [gopkgdoc page](http://go.pkgdoc.org/github.com/jzelinskie/reddit)

## license

Modified BSD License
