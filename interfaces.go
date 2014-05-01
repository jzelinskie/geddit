// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reddit

type voter interface {
	voteID() string
}

type deleter interface {
	deleteID() string
}

type replier interface {
	replyID() string
}
