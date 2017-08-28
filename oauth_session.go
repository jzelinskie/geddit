// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package geddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/google/go-querystring/query"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type transport struct {
	http.RoundTripper
	useragent string
}

// Any request headers can be modified here.
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", t.useragent)
	return t.RoundTripper.RoundTrip(req)
}

// OAuthSession represents an OAuth session with reddit.com --
// all authenticated API calls are methods bound to this type.
type OAuthSession struct {
	Client       *http.Client
	ClientID     string
	ClientSecret string
	OAuthConfig  *oauth2.Config
	//TokenExpiry  time.Time
	UserAgent string
	ctx       context.Context
	throttle  *rate.RateLimiter
}

// NewOAuthSession creates a new session for those who want to log into a
// reddit account via OAuth.
func NewOAuthSession(clientID, clientSecret, useragent, redirectURL string) (*OAuthSession, error) {
	o := &OAuthSession{}

	if len(useragent) > 0 {
		o.UserAgent = useragent
	} else {
		o.UserAgent = "Geddit Reddit Bot https://github.com/jzelinskie/geddit"
	}

	// Set OAuth config
	o.OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.reddit.com/api/v1/authorize",
			TokenURL: "https://www.reddit.com/api/v1/access_token",
		},
		RedirectURL: redirectURL,
	}
	// Inject our custom HTTP client so that a user-defined UA can
	// be passed during any authentication requests.
	c := &http.Client{}
	c.Transport = &transport{http.DefaultTransport, o.UserAgent}
	o.ctx = context.WithValue(context.Background(), oauth2.HTTPClient, c)
	return o, nil
}

// Throttle sets the interval of each HTTP request.
// Disable by setting interval to 0. Disabled by default.
// Throttling is applied to invidual OAuthSession types.
func (o *OAuthSession) Throttle(interval time.Duration) {
	if interval == 0 {
		o.throttle = nil
		return
	}
	o.throttle = rate.New(1, interval)
}

// LoginAuth creates the required HTTP client with a new token.
func (o *OAuthSession) LoginAuth(username, password string) error {
	// Fetch OAuth token.
	t, err := o.OAuthConfig.PasswordCredentialsToken(o.ctx, username, password)
	if err != nil {
		return err
	}
	if !t.Valid() {
		msg := "Invalid OAuth token"
		if t != nil {
			if extra := t.Extra("error"); extra != nil {
				msg = fmt.Sprintf("%s: %s", msg, extra)
			}
		}
		return errors.New(msg)
	}
	o.Client = o.OAuthConfig.Client(o.ctx, t)
	return nil
}

// AuthCodeURL creates and returns an auth URL which contains an auth code.
func (o *OAuthSession) AuthCodeURL(state string, scopes []string) string {
	o.OAuthConfig.Scopes = scopes
	return o.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// CodeAuth creates and sets a token using an authentication code returned from AuthCodeURL.
func (o *OAuthSession) CodeAuth(code string) error {
	t, err := o.OAuthConfig.Exchange(o.ctx, code)
	if err != nil {
		return err
	}
	o.Client = o.OAuthConfig.Client(o.ctx, t)
	return nil
}

// NeedsCaptcha check whether CAPTCHAs are needed for the Submit function.
func (o *OAuthSession) NeedsCaptcha() (bool, error) {
	var b bool
	err := o.getBody("https://oauth.reddit.com/api/needs_captcha", &b)
	if err != nil {
		return false, err
	}
	return b, nil
}

// NewCaptcha returns a string used to create CAPTCHA links for users.
func (o *OAuthSession) NewCaptcha() (string, error) {
	// Build form for POST request.
	v := url.Values{
		"api_type": {"json"},
	}

	type captcha struct {
		Json struct {
			Errors [][]string
			Data   struct {
				Iden string
			}
		}
	}
	c := &captcha{}

	err := o.postBody("https://oauth.reddit.com/api/new_captcha", v, c)
	if err != nil {
		return "", err
	}
	return c.Json.Data.Iden, nil
}

func (o *OAuthSession) getBody(link string, d interface{}) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}

	if o.Client == nil {
		return errors.New("OAuth Session lacks HTTP client! Use func (o OAuthSession) LoginAuth() to make one.")
	}

	// Throttle request
	if o.throttle != nil {
		o.throttle.Wait()
	}

	resp, err := o.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, d)
	if err != nil {
		return err
	}

	return nil
}

func (o *OAuthSession) Me() (*Redditor, error) {
	r := &Redditor{}
	err := o.getBody("https://oauth.reddit.com/api/v1/me", r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (o *OAuthSession) MyKarma() ([]Karma, error) {
	type karma struct {
		Data []Karma
	}
	k := &karma{}
	err := o.getBody("https://oauth.reddit.com/api/v1/me/karma", k)
	if err != nil {
		return nil, err
	}
	return k.Data, nil
}

func (o *OAuthSession) MyPreferences() (*Preferences, error) {
	p := &Preferences{}
	err := o.getBody("https://oauth.reddit.com/api/v1/me/prefs", p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (o *OAuthSession) MyFriends() ([]Friend, error) {
	type friends struct {
		Data struct {
			Children []Friend
		}
	}
	f := &friends{}
	err := o.getBody("https://oauth.reddit.com/api/v1/me/friends", f)
	if err != nil {
		return nil, err
	}
	return f.Data.Children, nil
}

func (o *OAuthSession) MyTrophies() ([]*Trophy, error) {
	type trophyData struct {
		Data struct {
			Trophies []struct {
				Data Trophy
			}
		}
	}

	t := &trophyData{}
	err := o.getBody("https://oauth.reddit.com/api/v1/me/trophies", t)
	if err != nil {
		return nil, err
	}

	var myTrophies []*Trophy
	for _, trophy := range t.Data.Trophies {
		myTrophies = append(myTrophies, &trophy.Data)
	}
	return myTrophies, nil
}

// Listing returns a slice of Submission pointers.
// See https://www.reddit.com/dev/api#listings for documentation.
func (o *OAuthSession) Listing(username, listing string, sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	p, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	if sort != "" {
		p.Set("sort", string(sort))
	}

	type resp struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}
	r := &resp{}
	url := fmt.Sprintf("https://oauth.reddit.com/user/%s/%s?%s", username, listing, p.Encode())
	err = o.getBody(url, r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

func (o *OAuthSession) Upvoted(username string, sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	return o.Listing(username, "upvoted", sort, params)
}

func (o *OAuthSession) MyUpvoted(sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	me, err := o.Me()
	if err != nil {
		return nil, err
	}
	return o.Listing(me.Name, "upvoted", sort, params)
}

// AboutRedditor returns a Redditor for the given username using OAuth.
func (o *OAuthSession) AboutRedditor(user string) (*Redditor, error) {
	type redditor struct {
		Data Redditor
	}
	r := &redditor{}
	link := fmt.Sprintf("https://oauth.reddit.com/user/%s/about", user)

	err := o.getBody(link, r)
	if err != nil {
		return nil, err
	}
	return &r.Data, nil
}

func (o *OAuthSession) UserTrophies(user string) ([]*Trophy, error) {
	type trophyData struct {
		Data struct {
			Trophies []struct {
				Data Trophy
			}
		}
	}

	t := &trophyData{}
	url := fmt.Sprintf("https://oauth.reddit.com/api/v1/user/%s/trophies", user)
	err := o.getBody(url, t)
	if err != nil {
		return nil, err
	}

	var trophies []*Trophy
	for _, trophy := range t.Data.Trophies {
		trophies = append(trophies, &trophy.Data)
	}
	return trophies, nil
}

// AboutSubreddit returns a subreddit for the given subreddit name using OAuth.
func (o *OAuthSession) AboutSubreddit(name string) (*Subreddit, error) {
	type subreddit struct {
		Data Subreddit
	}
	sr := &subreddit{}
	link := fmt.Sprintf("https://oauth.reddit.com/r/%s/about", name)

	err := o.getBody(link, sr)
	if err != nil {
		return nil, err
	}
	return &sr.Data, nil
}

// Comments returns the comments for a given Submission using OAuth.
func (o *OAuthSession) Comments(h *Submission, sort PopularitySort, params ListingOptions) ([]*Comment, error) {
	p, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	if sort != "" {
		p.Set("sort", string(sort))
	}

	var c interface{}
	link := fmt.Sprintf("https://oauth.reddit.com/comments/%s?%s", h.ID, p.Encode())
	err = o.getBody(link, &c)
	if err != nil {
		return nil, err
	}
	helper := new(helper)
	helper.buildComments(c)
	return helper.comments, nil
}

func (o *OAuthSession) postBody(link string, form url.Values, d interface{}) error {
	req, err := http.NewRequest("POST", link, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// This is needed to avoid rate limits
	//req.Header.Set("User-Agent", o.UserAgent)

	// POST form provided
	req.PostForm = form

	if o.Client == nil {
		return errors.New("OAuth Session lacks HTTP client! Use func (o OAuthSession) LoginAuth() to make one.")
	}

	// Throttle request
	if o.throttle != nil {
		o.throttle.Wait()
	}

	resp, err := o.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// The caller may want JSON decoded, or this could just be an update/delete request.
	if d != nil {
		err = json.Unmarshal(body, d)
		if err != nil {
			return err
		}
	}

	return nil
}

// Submit accepts a NewSubmission type and submits a new link using OAuth.
// Returns a Submission type.
func (o *OAuthSession) Submit(ns *NewSubmission) (*Submission, error) {

	// Build form for POST request.
	v := url.Values{
		"title":       {ns.Title},
		"url":         {ns.Content},
		"text":        {ns.Content},
		"sr":          {ns.Subreddit},
		"sendreplies": {strconv.FormatBool(ns.SendReplies)},
		"resubmit":    {strconv.FormatBool(ns.Resubmit)},
		"api_type":    {"json"},
		// TODO implement captchas for OAuth types
		//"captcha":     {ns.Captcha.Response},
		//"iden":        {ns.Captcha.Iden},
	}
	if ns.Self {
		v.Add("kind", "self")
	} else {
		v.Add("kind", "link")
	}

	type submission struct {
		Json struct {
			Errors [][]string
			Data   Submission
		}
	}
	submit := &submission{}

	err := o.postBody("https://oauth.reddit.com/api/submit", v, submit)
	if err != nil {
		return nil, err
	}
	// TODO check s.Errors and do something useful?
	return &submit.Json.Data, nil
}

// Delete deletes a link or comment using the given full name ID.
func (o *OAuthSession) Delete(d Deleter) error {
	// Build form for POST request.
	v := url.Values{}
	v.Add("id", d.deleteID())

	return o.postBody("https://oauth.reddit.com/api/del", v, nil)
}

// SubredditSubmissions returns the submissions on the given subreddit using OAuth.
func (o *OAuthSession) SubredditSubmissions(subreddit string, sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	baseUrl := "https://oauth.reddit.com"

	// If subbreddit given, add to URL
	if subreddit != "" {
		baseUrl += "/r/" + subreddit
	}

	redditURL := fmt.Sprintf(baseUrl+"/%s.json?%s", sort, v.Encode())

	type Response struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}

	r := new(Response)
	err = o.getBody(redditURL, r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

// Frontpage returns the submissions on the default reddit frontpage using OAuth.
func (o *OAuthSession) Frontpage(sort PopularitySort, params ListingOptions) ([]*Submission, error) {
	return o.SubredditSubmissions("", sort, params)
}

// Vote either votes or rescinds a vote for a Submission or Comment using OAuth.
func (o *OAuthSession) Vote(v Voter, dir Vote) error {
	// Build form for POST request.
	form := url.Values{
		"id":  {v.voteID()},
		"dir": {string(dir)},
	}
	var vo interface{}

	err := o.postBody("https://oauth.reddit.com/api/vote", form, vo)
	if err != nil {
		return err
	}
	return nil
}

// Reply posts a comment as a response to a Submission or Comment using OAuth.
func (o OAuthSession) Reply(r Replier, comment string) (*Comment, error) {
	// Build form for POST request.
	form := url.Values{
		"api_type": {"json"},
		"thing_id": {r.replyID()},
		"text":     {comment},
	}

	type response struct {
		JSON struct {
			Errors [][]string
			Data   struct {
				Things []struct {
					Data map[string]interface{}
				}
			}
		}
	}

	res := &response{}

	err := o.postBody("https://oauth.reddit.com/api/comment", form, res)
	if err != nil {
		return nil, err
	}

	if len(res.JSON.Errors) != 0 {
		var msg []string
		for _, k := range res.JSON.Errors {
			msg = append(msg, k[1])
		}
		return nil, errors.New(strings.Join(msg, ", "))
	}

	c := makeComment(res.JSON.Data.Things[0].Data)

	return c, nil
}

// Save saves a link or comment using OAuth.
func (o *OAuthSession) Save(v Voter, category string) error {
	// Build form for POST request.
	form := url.Values{
		"id":       {v.voteID()},
		"category": {category},
	}
	var s interface{}

	err := o.postBody("https://oauth.reddit.com/api/save", form, s)
	if err != nil {
		return err
	}
	return nil
}

// Unsave saves a link or comment using OAuth.
func (o *OAuthSession) Unsave(v Voter, category string) error {
	// Build form for POST request.
	form := url.Values{
		"id":       {v.voteID()},
		"category": {category},
	}
	var u interface{}

	err := o.postBody("https://oauth.reddit.com/api/unsave", form, u)
	if err != nil {
		return err
	}
	return nil
}

// SavedLinks fetches links saved by given username using OAuth.
func (o *OAuthSession) SavedLinks(username string, params ListingOptions) ([]*Submission, error) {
	return o.Listing(username, "saved", "", params)
}

// MySavedLinks fetches links saved by current user using OAuth.
func (o *OAuthSession) MySavedLinks(params ListingOptions) ([]*Submission, error) {
	me, err := o.Me()
	if err != nil {
		return nil, err
	}
	return o.Listing(me.Name, "saved", "", params)
}

// SavedComments fetches comments saved by given username using OAuth.
func (o *OAuthSession) SavedComments(user string, params ListingOptions) ([]*Comment, error) {
	var s interface{}
	url := fmt.Sprintf("https://oauth.reddit.com/user/%s/saved", user)
	err := o.getBody(url, &s)
	if err != nil {
		return nil, err
	}

	helper := new(helper)
	helper.buildComments(s)
	return helper.comments, nil
}

// MySavedComments fetches comments saved by current user using OAuth.
func (o *OAuthSession) MySavedComments(params ListingOptions) ([]*Comment, error) {
	me, err := o.Me()
	if err != nil {
		return nil, err
	}
	return o.SavedComments(me.Name, params)
}

// MySubreddits fetchs subreddits the current user subscribes to.
// TODO support other endpoints https://www.reddit.com/dev/api/#GET_subreddits_mine_{where}
func (o *OAuthSession) MySubreddits() ([]*Subreddit, error) {
	type Response struct {
		Data struct {
			Children []struct {
				Data *Subreddit
			}
		}
	}
	r := new(Response)
	err := o.getBody("https://oauth.reddit.com/subreddits/mine/subscriber", r)
	if err != nil {
		return nil, err
	}
	s := make([]*Subreddit, len(r.Data.Children))
	for i, child := range r.Data.Children {
		s[i] = child.Data
	}
	return s, nil
}

// SubredditComments fetches all the new comments in a subreddit, and returns them in a slice of Comment structs
// This function uses www.reddit.com instead of the OAuth API as the latter doesn't have an endpoint for a particular subreddit's comments
func (o *OAuthSession) SubredditComments(subreddit string) ([]*Comment, error) {
	baseURL := "https://www.reddit.com"

	if subreddit != "" {
		baseURL += "/r/" + subreddit
	}

	subCommentsURL := baseURL + "/comments.json"

	req := request{
		url:       subCommentsURL,
		useragent: o.UserAgent,
	}

	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	var comments interface{}
	if err = json.NewDecoder(body).Decode(&comments); err != nil {
		return nil, err
	}

	helper := new(helper)
	helper.buildComments(comments)
	return helper.comments, nil
}
