package main

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type Sender interface {
	Send(text string) error
}

type fakeSender struct {
	msgs []string
}

func NewSender(api *anaconda.TwitterApi, myName string) Sender {
	// 생성에 필요한 정보중 일부를 의도적으로 넣지 않으면 fake로 취급
	if api == nil || myName == "" {
		return &fakeSender{
			msgs: []string{},
		}
	}

	return &directMessageSender{
		api:    api,
		myName: myName,
	}
}

func (s *fakeSender) Send(text string) error {
	s.msgs = append(s.msgs, text)
	fmt.Printf("Fake Send : %s\n", text)
	return nil
}

type directMessageSender struct {
	api    *anaconda.TwitterApi
	myName string
}

func (s *directMessageSender) Send(text string) error {
	_, err := s.api.PostDMToScreenName(text, s.myName)
	return err
}
