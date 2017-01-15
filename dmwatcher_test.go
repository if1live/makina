package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseLine(t *testing.T) {
	cases := []struct {
		line string

		cmdtype reflect.Type
		args    []string
	}{
		{"asdf", reflect.TypeOf(nil), []string{}},
		{"", reflect.TypeOf(nil), []string{}},

		{"help", reflect.TypeOf(&HelpCommand{}), []string{}},
		{"HELp", reflect.TypeOf(&HelpCommand{}), []string{}},

		{"status", reflect.TypeOf(&StatusCommand{}), []string{}},
		{"sentry", reflect.TypeOf(&SentryCommand{}), []string{}},
		{"hitomi_preview, 123456", reflect.TypeOf(&HitomiPreviewCommand{}), []string{"123456"}},
	}
	watcher := &DirectMessageWatcher{}
	for _, c := range cases {
		cmd, _, args := watcher.ParseLine(c.line)
		assert.Equal(t, c.cmdtype, reflect.TypeOf(cmd))
		assert.Equal(t, true, reflect.DeepEqual(c.args, args))
	}
}

func Test_ParseURL(t *testing.T) {
	cases := []struct {
		text    string
		cmdtype reflect.Type
	}{
		{"https://twitter.com/if1live/status/785285765697720320", reflect.TypeOf(&TweetSaveCommand{})},
		{"http://twitter.com/if1live/status/785285765697720320", reflect.TypeOf(&TweetSaveCommand{})},
		{"//twitter.com/if1live/status/785285765697720320", reflect.TypeOf(&TweetSaveCommand{})},

		{"https://twitter.com/if1live/invalid/785285765697720320", reflect.TypeOf(nil)},

		{"https://google.com", reflect.TypeOf(nil)},
		{"", reflect.TypeOf(nil)},
	}
	watcher := &DirectMessageWatcher{}
	for _, c := range cases {
		cmd := watcher.ParseURL(c.text)
		assert.Equal(t, c.cmdtype, reflect.TypeOf(cmd))
	}
}

func Test_getStatusID(t *testing.T) {
	cases := []struct {
		rawurl string
		id     int64
	}{
		{"https://twitter.com/if1live/status/785285765697720320", 785285765697720320},
		{"invalid", -1},
	}
	cmd := &TweetSaveCommand{}
	for _, c := range cases {
		id := cmd.getStatusID(c.rawurl)
		assert.Equal(t, c.id, id)
	}
}
