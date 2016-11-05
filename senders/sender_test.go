package senders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	strategy := NewFake()
	sender := New(strategy)

	title := "this-is-title"
	body := "this-is-body"
	sender.Send(title, body)

	msg := sender.GetLastMessage()
	assert.Equal(t, title, msg.Title)
	assert.Equal(t, body, msg.Body)
}
