package geddit

import "testing"

func TestSubmit(t *testing.T) {

	session, err := NewLoginSession(
		"redditgolang",
		"apitest11",
		"tester",
	)
	if err != nil {
		t.Error(err)
	}

	subreddit, err := session.AboutSubreddit("mybottester")
	if err != nil {
		t.Error(err)
	}

	needsCaptcha, err := session.NeedsCaptcha()
	if err != nil {
		t.Error(err)
	}

	t.Log(needsCaptcha)

	if needsCaptcha {
		iden, err := session.NewCaptchaIden()
		if err != nil {
			t.Error(err)
		}

		_, err = session.CaptchaImage(iden)
		if err != nil {
			t.Error(err)
		}

		err = session.Submit(NewTextSubmission(subreddit.Name, "CAPTCHA TESTING TEXT", "TEST TEXT", true, &Captcha{iden, "test"}))
		if err != nil {
			t.Error(err)
		}

		err = session.Submit(NewLinkSubmission(subreddit.Name, "CAPTCHA TESTING LINK", "https://github.com/jzelinskie/reddit", true, &Captcha{iden, "test"}))
		if err != nil {
			t.Error(err)
		}

	} else {

		err = session.Submit(NewTextSubmission(subreddit.Name, "TESTING TEXT", "TEST TEXT", true, &Captcha{}))
		if err != nil {
			t.Error(err)
		}

		err = session.Submit(NewLinkSubmission(subreddit.Name, "TESTING LINK", "https://github.com/jzelinskie/reddit", true, &Captcha{}))
		if err != nil {
			t.Error(err)
		}
	}

}

func TestListings(t *testing.T) {
	session, err := NewLoginSession(
		"redditgolang",
		"apitest11",
		"tester",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.MySaved(NewSubmissions, "")
	if err != nil {
		t.Fatal(err)
	}
}
