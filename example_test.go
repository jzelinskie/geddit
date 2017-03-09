// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit_test

import (
	"fmt"
	"log"

	"github.com/jzelinskie/geddit"
)

func ExampleNewOAuthSession_login() {
	o, err := geddit.NewOAuthSession(
		"client_id",
		"client_secret",
		"Testing OAuth Bot by u/my_user v0.1 see source https://github.com/jzelinskie/geddit",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create new auth token for confidential clients (personal scripts/apps).
	err = o.LoginAuth("my_user", "my_password")
	if err != nil {
		log.Fatal(err)
	}

	// Ready to make API calls!
}

func ExampleNewOAuthSession_url() {
	o, err := geddit.NewOAuthSession(
		"client_id",
		"client_secret",
		"Testing OAuth Bot by u/my_user v0.1 see source https://github.com/jzelinskie/geddit",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Pass a random/unique state string which will be returned to the
	//   redirect URL. Ideally, you should verify that it matches to
	//   avoid CSRF attack.
	url := o.AuthCodeURL("random string", []string{"indentity", "read"})
	fmt.Printf("Visit %s to obtain auth code", url)

	var code string
	fmt.Scanln(&code)

	// Create and set token using given auth code.
	err = o.CodeAuth(code)
	if err != nil {
		log.Fatal(err)
	}

	// Ready to make API calls!
}
