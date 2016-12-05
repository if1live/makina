package senders

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
)

type DirectMessageSendStrategy struct {
	api    *anaconda.TwitterApi
	myName string
}

func NewDirectMessage(api *anaconda.TwitterApi, myName string) SendStrategy {
	return &DirectMessageSendStrategy{
		api:    api,
		myName: myName,
	}
}

// 트위터는 제목과 내용의 구분이 없다
func makeContent(title, body string) string {
	if len(body) == 0 {
		return title
	}
	return fmt.Sprintf("title: %s\r\n%s", title, body)
}
func (s *DirectMessageSendStrategy) Send(title, body string) {
	content := makeContent(title, body)
	_, err := s.api.PostDMToScreenName(content, s.myName)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
}
