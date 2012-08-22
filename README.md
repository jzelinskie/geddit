#reddit
an abstraction for the [reddit.com](http://reddit.com) API

## warning

this is a **work in progress**! future features include

* more public api calls
* auth-only api calls (blame net/http not having a cookiejar)
* comments

## example

     package main
     
     import (
       "fmt"
       "github.com/jzelinskie/reddit"
     )

     func main() {
       s := reddit.NewSession()
       
       headlines, err := s.FrontPage()
       if err != nil {
         panic(err)
       }
       
       for _, hl := range headlines {
         fmt.Println(hl)
       }
     }

## docs

Checkout the [gopkgdoc page](http://go.pkgdoc.org/github.com/jzelinskie/reddit)

## license

Modified BSD License
