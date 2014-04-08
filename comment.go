// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

import (
	"fmt"
)

type Comment struct {
	Author              string  //`json:"author"`
	Body                string  //`json:"body"`
	BodyHTML            string  //`json:"body_html"`
	Subreddit           string  //`json:"subreddit"`
	LinkId              string  //`json:"link_id"`
	ParentId            string  //`json:"parent_id"`
	SubredditId         string  //`json:"subreddit_id"`
	FullId              string  //`json:"name"`
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

func (c Comment) String() string {
	return fmt.Sprintf("%s (%d/%d): %s", c.Author, c.UpVotes, c.DownVotes, c.Body)
}

// Does the ugly work of setting the comment fields
func makeComment(cmap map[string]interface{}) *Comment {
	ret := new(Comment)
	ret.Author = cmap["author"].(string)
	ret.Body = cmap["body"].(string)
	ret.BodyHTML = cmap["body_html"].(string)
	ret.Subreddit = cmap["subreddit"].(string)
	ret.LinkId = cmap["link_id"].(string)
	ret.ParentId = cmap["parent_id"].(string)
	ret.SubredditId = cmap["subreddit_id"].(string)
	ret.FullId = cmap["name"].(string)
	ret.UpVotes = cmap["ups"].(float64)
	ret.DownVotes = cmap["downs"].(float64)
	ret.Created = cmap["created_utc"].(float64)

	//These fields commented out because they threw runtime errors in type assertion

	//ret.Edited = cmap["edited"].(bool)
	//ret.BannedBy = cmap["banned_by"].(*string)
	//ret.ApprovedBy = cmap["approved_by"].(*string)
	//ret.AuthorFlairTxt = cmap["author_flair_text"].(*string)
	//ret.AuthorFlairCSSClass = cmap["author_flair_css_class"].(*string)
	//ret.NumReports = cmap["num_reports"].(*int)
	//ret.Likes = cmap["likes"].(*int)

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
