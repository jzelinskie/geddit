package reddit

import "testing"

func TestSubmit(t *testing.T) {

	session, err := NewLoginSession(
		"Username",
		"Password",
		"tester",
	)

	if err != nil {
		t.Error(err)
	}

	subreddit, err := session.AboutSubreddit("mybottester")

	if err != nil {
		t.Error(err)
	}

	err = session.Submit(subreddit, newSubmission{"TESTING SELF", "TEST TEXT", true, true, true})

	if err != nil {
		t.Error(err)
	}

	err = session.Submit(subreddit, newSubmission{"TESTING LINK", "https://github.com/jzelinskie/reddit", false, true, true})

	if err != nil {
		t.Error(err)
	}

}
