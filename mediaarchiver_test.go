package main

import "testing"
import "github.com/stretchr/testify/assert"

func Test_FindCategory(t *testing.T) {
	cases := []struct {
		dir        string
		screenName string
		expected   string
	}{
		{"retweet", "normal", "retweet"},
		{"retweet", "save-target", "user-save-target"},
		{"retweet", "SAVE-TARGET", "user-save-target"},
	}

	ar := &MediaArchiver{
		storage:         nil,
		myName:          "myName",
		predefinedUsers: []string{"save-target"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, ar.FindCategory(c.dir, c.screenName))
	}
}
