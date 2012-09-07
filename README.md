#reddit
an abstraction for the [reddit.com](http://reddit.com) API

## warning

this is a **work in progress**! future features include

* more public api calls
* auth'd api calls
* comments (what data structure should I use?)

## example

```Go
package main

import (
 "fmt"
 "github.com/jzelinskie/reddit"
)

func main() {
 // Get the frontpage.
 headlines, err := reddit.DefaultHeadlines()
 if err != nil {
   panic(err)
 }

 for _, hl := range headlines {
   fmt.Println(hl)
 }
}
```

## docs

Checkout the [gopkgdoc page](http://go.pkgdoc.org/github.com/jzelinskie/reddit)

## license

Modified BSD License
