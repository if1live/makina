package main

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type SimpleTweet struct {
	ID         int64  `yaml:"id"`
	ScreenName string `yaml:"screen_name"`
	UserName   string `yaml:"user_name"`
	Text       string `yaml:"text"`
	URL        string `yaml:"url"`
	CreatedAt  string `yaml:"created_at"`
    Extra      string `yaml:"extra"`
}

func NewSimpleTweet(t *anaconda.Tweet) SimpleTweet {
	// 리트윗의 경우 원형을 저장한다
	if t.RetweetedStatus != nil {
		return NewSimpleTweet(t.RetweetedStatus)
	}

	url := fmt.Sprintf("https://twitter.com/%s/status/%s", t.User.ScreenName, t.IdStr)
	return SimpleTweet{
		Text:       t.Text,
		ID:         t.Id,
		ScreenName: t.User.ScreenName,
		UserName:   t.User.Name,
		URL:        url,
		CreatedAt:  t.CreatedAt,
        Extra:      "",
	}
}
