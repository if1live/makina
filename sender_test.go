package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	sender := NewSender(nil, "name")

	text := "this-is-text"
	err := sender.Send(text)
	assert.Equal(t, nil, err)
}
