package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
)

type DirectMessageWatcher struct {
	myName  string
	storage *Storage
	sender  Sender
}

func NewDirectMessageWatcher(myName string, storage *Storage, sender Sender) MessageRule {
	r := &DirectMessageWatcher{
		myName:  myName,
		storage: storage,
		sender:  sender,
	}
	return r
}

func (r *DirectMessageWatcher) OnDirectMessage(dm *anaconda.DirectMessage) {
	if dm.Sender.ScreenName != r.myName {
		return
	}
	if len(dm.Text) == 0 {
		return
	}

	errorMsg := "Invalid command"
	if dm.Text == errorMsg {
		return
	}

	cmd, cmdname, args := r.ParseLine(dm.Text)
	if cmd != nil {
		msg := fmt.Sprintf("cmd=%s, args=%s is enqueued.", cmdname, strings.Join(args, ","))
		r.sender.Send(msg)
		go cmd.execute(args, r.sender)

	} else {
		r.sender.Send(errorMsg)
	}
}

func (r *DirectMessageWatcher) findCommand(text string) LineCommand {
	cmds := map[string]LineCommand{
		"hitomi_preview": NewHitomiPreviewCommand(r.storage),
		"status":         NewStatusCommand(),
		"sentry":         NewSentryCommand(),
	}
	helpCmd := NewHelpCommand(cmds)

	for key, cmd := range cmds {
		if strings.HasPrefix(text, key) {
			return cmd
		}
	}
	if strings.HasPrefix(text, "help") {
		return helpCmd
	}
	return nil
}

func (r *DirectMessageWatcher) ParseLine(text string) (LineCommand, string, []string) {
	tokens := strings.Split(text, ",")
	for i, token := range tokens {
		tokens[i] = strings.Trim(token, " ")
	}

	cmdname := strings.ToLower(tokens[0])
	cmd := r.findCommand(cmdname)
	args := tokens[1:]
	return cmd, cmdname, args
}

// 간단한 텍스트 기반의 명령어
type LineCommand interface {
	execute(args []string, sender Sender)
	help() string
}

type HelpCommand struct {
	cmds map[string]LineCommand
}

func NewHelpCommand(cmds map[string]LineCommand) LineCommand {
	return &HelpCommand{
		cmds: cmds,
	}
}
func (c *HelpCommand) help() string {
	return "help"
}
func (c *HelpCommand) execute(args []string, sender Sender) {
	log.Printf("DM: help")
	lines := make([]string, len(c.cmds))

	idx := 0
	for key, cmd := range c.cmds {
		lines[idx] = key + " : " + cmd.help()
		idx++
	}

	help := strings.Join(lines, "\n")
	sender.Send(help)
}

type HitomiPreviewCommand struct {
	storage *Storage
}

func NewHitomiPreviewCommand(storage *Storage) LineCommand {
	return &HitomiPreviewCommand{storage}
}

func (c *HitomiPreviewCommand) execute(args []string, sender Sender) {
	log.Printf("DM: hitomi preview %s", strings.Join(args, ","))

	reCode := regexp.MustCompile(`^(\d+)$`)
	for _, arg := range args {
		if !reCode.MatchString(arg) {
			continue
		}
		go func(code string) {
			ok := FetchHitomiPreview(code, nil, c.storage)
			msg := ""
			if ok {
				msg = fmt.Sprintf("hitomi success : %s", code)
			} else {
				msg = fmt.Sprintf("hitomi fail : %s", code)
			}
			sender.Send(msg)
		}(arg)
	}
}

func (c *HitomiPreviewCommand) help() string {
	return "<hitomi_code>"
}

type StatusCommand struct {
}

func NewStatusCommand() LineCommand {
	return &StatusCommand{}
}

func (c *StatusCommand) help() string {
	return "server status"
}
func (c *StatusCommand) execute(args []string, sender Sender) {
	log.Printf("DM: status")
	now := time.Now()
	msg := now.Format(time.RFC3339)
	sender.Send("still alive : " + msg)
}

type SentryCommand struct {
}

func NewSentryCommand() LineCommand {
	return &SentryCommand{}
}
func (c *SentryCommand) help() string {
	return "send sentry event, for development only"
}
func (c *SentryCommand) execute(args []string, sender Sender) {
	log.Printf("DM: sentry")
	_, err := os.Open("invalid-file-to-raise-error")
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		sender.Send("sentry success")
	} else {
		sender.Send("sentry fail")
	}
}
