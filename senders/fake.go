package senders

import "fmt"

type FakeSendStrategy struct {
	msgs []Message
}

func NewFake() SendStrategy {
	return &FakeSendStrategy{
		msgs: []Message{},
	}
}
func (s *FakeSendStrategy) Send(title, body string) {
	msg := Message{
		Title: title,
		Body:  body,
	}
	s.msgs = append(s.msgs, msg)
	fmt.Printf("Fake Send : title=%s, body=%s\n", title, body)
}
