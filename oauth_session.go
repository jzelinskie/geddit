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

// OAuthSession represents an OAuth session with reddit.com --
// all authenticated API calls are methods bound to this type.
type OAuthSession struct {
	Client       *http.Client
	ClientID     string
	ClientSecret string
	OAuthConfig  *oauth2.Config
	UserAgent    string
	ctx          context.Context
	throttle     *rate.RateLimiter
}

// NewLoginSession creates a new session for those who want to log into a
// reddit account via OAuth.
func NewOAuthSession(clientID, clientSecret, useragent string) (*OAuthSession, error) {
	s := &OAuthSession{}

	if useragent != "" {
		s.UserAgent = useragent
	} else {
		s.UserAgent = "Geddit Reddit Bot https://github.com/jzelinskie/geddit"
	}

	// Set OAuth config
	// TODO Set user-defined scopes
	s.OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{}, //"identity", "read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.reddit.com",
			TokenURL: "https://www.reddit.com/api/v1/access_token",
		},
	}
	s.ctx = context.Background()
	return s, nil
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
	o.Client = o.OAuthConfig.Client(context.Background(), t) //o.ctx, t)
	return nil
}

func (o *OAuthSession) getBody(link string, d interface{}) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}

	// This is needed to avoid rate limits
	req.Header.Set("User-Agent", o.UserAgent)

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

	// DEBUG
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("***DEBUG***\nRequest Body: %s\n***DEBUG***\n\n", body)

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
func (o *OAuthSession) Listing(username, listing string, sort popularitySort, after string) ([]*Submission, error) {
	values := &url.Values{}
	if sort != "" {
		values.Set("sort", string(sort))
	}
	if after != "" {
		values.Set("after", after)
	}

	type resp struct {
		Data struct {
			Children []struct {
				Data *Submission
			}
		}
	}
	r := &resp{}
	url := fmt.Sprintf("https://oauth.reddit.com/user/%s/%s?%s", username, listing, values.Encode())
	err := o.getBody(url, r)
	if err != nil {
		return nil, err
	}

	submissions := make([]*Submission, len(r.Data.Children))
	for i, child := range r.Data.Children {
		submissions[i] = child.Data
	}

	return submissions, nil
}

func (o *OAuthSession) MyUpvoted(sort popularitySort, after string) ([]*Submission, error) {
	me, err := o.Me()
	if err != nil {
		return nil, err
	}
	return o.Listing(me.Name, "upvoted", sort, after)
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
func (o *OAuthSession) Comments(h *Submission) ([]*Comment, error) {
	var c interface{}
	link := fmt.Sprintf("https://oauth.reddit.com/comments/%s", h.ID)
	err := o.getBody(link, &c)
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
	req.Header.Set("User-Agent", o.UserAgent)

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
	// DEBUG
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("***DEBUG***\nRequest Body: %s\n***DEBUG***\n\n", body)

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
func (o *OAuthSession) SubredditSubmissions(subreddit string, sort popularitySort, params ListingOptions) ([]*Submission, error) {
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
func (o *OAuthSession) Frontpage(sort popularitySort, params ListingOptions) ([]*Submission, error) {
	return o.SubredditSubmissions("", sort, params)
}

// Vote either votes or rescinds a vote for a Submission or Comment using OAuth.
func (o *OAuthSession) Vote(v Voter, dir vote) error {
	// Build form for POST request.
	form := url.Values{
		"id":  {v.voteID()},
		"dir": {string(dir)},
	}
	var foo interface{}

	err := o.postBody("https://oauth.reddit.com/api/vote", form, foo)
	if err != nil {
		return err
	}
	return nil
}
