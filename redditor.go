// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"fmt"
)

type Redditor struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	LinkKarma    int     `json:"link_karma"`
	CommentKarma int     `json:"comment_karma"`
	Created      float32 `json:"created_utc"`
	Gold         bool    `json:"is_gold"`
	Mod          bool    `json:"is_mod"`
	Mail         *bool   `json:"has_mail"`
	ModMail      *bool   `json:"has_mod_mail"`
}

// String returns the string representation of a reddit user.
func (r *Redditor) String() string {
	return fmt.Sprintf("%s (%d-%d)", r.Name, r.LinkKarma, r.CommentKarma)
}
