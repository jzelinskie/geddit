// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"fmt"
)

// Comment represents a reddit comment.
type Comment struct {
	Author              string  //`json:"author"`
	Body                string  //`json:"body"`
	BodyHTML            string  //`json:"body_html"`
	Subreddit           string  //`json:"subreddit"`
	LinkID              string  //`json:"link_id"`
	ParentID            string  //`json:"parent_id"`
	SubredditID         string  //`json:"subreddit_id"`
	FullID              string  //`json:"name"`
	Permalink           string  //`json:"permalink"`
	Score               float64 //`json:"score"`
	UpVotes             float64 //`json:"ups"`
	DownVotes           float64 //`json:"downs"`
	Created             float64 //`json:"created_utc"`
	Edited              bool    //`json:"edited"`
	BannedBy            *string //`json:"banned_by"`
	ApprovedBy          *string //`json:"approved_by"`
	AuthorFlairTxt      *string //`json:"author_flair_text"`
	AuthorFlairCSSClass *string //`json:"author_flair_css_class"`
	NumReports          *int    //`json:"num_reports"`
	Likes               *int    //`json:"likes"`
	Replies             []*Comment
}

func (c Comment) voteID() string   { return c.FullID }
func (c Comment) deleteID() string { return c.FullID }
func (c Comment) replyID() string  { return c.FullID }

// FullPermalink returns the full URL of a Comment.
func (c Comment) FullPermalink() string {
	return "https://reddit.com" + c.Permalink
}

func (c Comment) String() string {
	return fmt.Sprintf("%s (%d/%d): %s", c.Author, c.UpVotes, c.DownVotes, c.Body)
}

// makeComment tries its best to fill as many fields as possible of a Comment.
func makeComment(cmap map[string]interface{}) *Comment {
	ret := new(Comment)
	ret.Author, _ = cmap["author"].(string)
	ret.Body, _ = cmap["body"].(string)
	ret.BodyHTML, _ = cmap["body_html"].(string)
	ret.Subreddit, _ = cmap["subreddit"].(string)
	ret.LinkID, _ = cmap["link_id"].(string)
	ret.ParentID, _ = cmap["parent_id"].(string)
	ret.SubredditID, _ = cmap["subreddit_id"].(string)
	ret.FullID, _ = cmap["name"].(string)
	ret.Permalink, _ = cmap["permalink"].(string)
	ret.Score, _ = cmap["score"].(float64)
	ret.UpVotes, _ = cmap["ups"].(float64)
	ret.DownVotes, _ = cmap["downs"].(float64)
	ret.Created, _ = cmap["created_utc"].(float64)
	ret.Edited, _ = cmap["edited"].(bool)
	ret.BannedBy, _ = cmap["banned_by"].(*string)
	ret.ApprovedBy, _ = cmap["approved_by"].(*string)
	ret.AuthorFlairTxt, _ = cmap["author_flair_text"].(*string)
	ret.AuthorFlairCSSClass, _ = cmap["author_flair_css_class"].(*string)
	ret.NumReports, _ = cmap["num_reports"].(*int)
	ret.Likes, _ = cmap["likes"].(*int)

	helper := new(helper)
	helper.buildComments(cmap["replies"])
	ret.Replies = helper.comments

	return ret
}

//Helper struct to keep our interesting stuff
type helper struct {
	comments []*Comment
}

//Recursive function to find the fields we want and build the Comments
//Way too hackish for my likes
func (h *helper) buildComments(inf interface{}) {
	switch tp := inf.(type) {
	case []interface{}: //Maybe array for base comments
		for _, k := range tp {
			h.buildComments(k)
		}
	case map[string]interface{}: //Maybe comment data
		if tp["body"] == nil {
			for _, k := range tp {
				h.buildComments(k)
			}
		} else {
			h.comments = append(h.comments, makeComment(tp))
		}
	}
}
