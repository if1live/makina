package main

import "testing"

import "reflect"
import "github.com/stretchr/testify/assert"

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
