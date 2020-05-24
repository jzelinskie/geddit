// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
)

type RewriteTransport struct {
	Transport http.RoundTripper
	URL       *url.URL
}

func (t RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.URL.Scheme
	req.URL.Host = t.URL.Host
	req.URL.Path = path.Join(t.URL.Path, req.URL.Path)
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(req)
}

func testTools(code int, body string) (*httptest.Server, *OAuthSession) {
	// Dummy server to write JSON body provided
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatalf("Failed to parse local server URL: %v", err)
	}
	o := &OAuthSession{Client: http.DefaultClient, UserAgent: "Geddit Test"}
	o.Client.Transport = RewriteTransport{URL: u}

	return server, o
}

// Test defaults o fresh OAuthSession type.
func TestNewOAuthSession(t *testing.T) {
	o, err := NewOAuthSession("user", "pw", "agent", "http://")
	if err != nil {
		t.Fatal(err)
	}

	if o.Client != nil {
		t.Fatal(errors.New("HTTP client created before auth token!"))
	}
}

func TestMe(t *testing.T) {
	server, oauth := testTools(200, `{"has_mail": false, "name": "aggrolite", "is_friend": false, "created": 1278447313.0, "suspension_expiration_utc": null, "hide_from_robots": true, "is_suspended": false, "modhash": "XXX", "created_utc": 1278418513.0, "link_karma": 2327, "comment_karma": 1233, "over_18": true, "is_gold": false, "is_mod": true, "id": "45xiz", "gold_expiration": null, "inbox_count": 0, "has_verified_email": true, "gold_creddits": 0, "has_mod_mail": false}`)
	defer server.Close()

	me, err := oauth.Me()
	if err != nil {
		t.Errorf("Me() Test failed: %v", err)
	}
	// Sanity check just a few fields?
	if me.Name != "aggrolite" {
		t.Fatalf("Me() returned unexpected name: %s", me.Name)
	}
	if me.ID != "45xiz" {
		t.Fatalf("Me() returned unexpected ID: %s", me.ID)
	}
	if me.String() != "aggrolite (2327-1233)" {
		t.Fatalf("Me.String() returns unexpected result: %s", me.String())
	}
	fmt.Println(me)
}

func TestLink(t *testing.T) {
	server, oauth := testTools(200, `{"data": {"children": [{"data": {"name": "t3_12345", "id": "12345", "title": "My Title"}}]}}`)
	defer server.Close()

	link, err := oauth.Link("t3_12345")
	if err != nil {
		t.Errorf("Link() Test failed: %v", err)
	}

	if link.FullID != "t3_12345" {
		t.Fatalf("Link() returned unexpected full ID: %s", link.FullID)
	}
	if link.ID != "12345" {
		t.Fatalf("Link() returned unexpected ID: %s", link.ID)
	}
	if link.Title != "My Title" {
		t.Fatalf("Link() returned unexpected title: %s", link.Title)
	}
}

func TestComment(t *testing.T) {
	server, oauth := testTools(200, `{"data": {"children": [{"data": {"name": "t1_12345", "author": "u/me", "body": "username checks out", "archived": false}}]}}`)
	defer server.Close()

	comment, err := oauth.Comment("", "t1_12345")
	if err != nil {
		t.Errorf("Comment() Test failed: %v", err)
	}

	if comment.FullID != "t1_12345" {
		t.Fatalf("Comment() returned unexpected full ID: %s", comment.FullID)
	}
	if comment.Author != "u/me" {
		t.Fatalf("Comment() returned unexpected ID: %s", comment.Author)
	}
	if comment.Body != "username checks out" {
		t.Fatalf("Comment() returned unexpected body: %s", comment.Body)
	}
	if comment.Archived != false {
		t.Fatalf("Comment() returned wrong archived value: %v", comment.Archived)
	}
}

func TestUserPosts(t *testing.T) {
	server, oauth := testTools(200, `{"data": {"children": [{"data": {"name": "t3_12345", "title": "My Title", "permalink": "www.example.com", "subreddit": "example"}}, {"data": {"name": "t3_56789", "title": "My Title", "permalink": "www.notexample.com", "subreddit": "notexample"}}]}}`)
	defer server.Close()

	posts, err := oauth.UserPosts("example", "u/me", NewSubmissions, ListingOptions{})
	if err != nil {
		t.Errorf("UserPosts() Test failed: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("UserPosts() failed to filter by subreddit (got %d results, want %d results", len(posts), 1)
	}
	p := posts[0]
	if p.FullID != "t3_12345" {
		t.Fatalf("UserPosts() returned unexpected full ID: %s", p.FullID)
	}
	if p.Title != "My Title" {
		t.Fatalf("UserPosts() returned unexpected title: %s", p.Author)
	}
	if p.Permalink != "www.example.com" {
		t.Fatalf("UserPosts() returned unexpected permalink: %s", p.Permalink)
	}
	if p.Subreddit != "example" {
		t.Fatalf("UserPosts() returned wrong subreddit: %v", p.Subreddit)
	}
}
