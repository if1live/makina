package rules

import (
	"log"
	"os"
	"regexp"

	"time"

	"fmt"

	"strings"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
	"github.com/if1live/makina/hitomiwatcher"
	"github.com/if1live/makina/senders"
	"github.com/if1live/makina/storages"
)

type DirectMessageWatcher struct {
	MyName       string
	Accessor     storages.Accessor
	StatusSender *senders.Sender
}

func NewDirectMessageWatcher(MyName string, Accessor storages.Accessor, StatusSender *senders.Sender) Rule {
	r := &DirectMessageWatcher{
		MyName:       MyName,
		Accessor:     Accessor,
		StatusSender: StatusSender,
	}
	return r
}

func (r *DirectMessageWatcher) OnTweet(tweet *anaconda.Tweet) {

}
func (r *DirectMessageWatcher) OnEvent(ev string, event *anaconda.EventTweet) {

}
func (r *DirectMessageWatcher) OnDirectMessage(dm *anaconda.DirectMessage) {
	if dm.Sender.ScreenName != r.MyName {
		return
	}

	errorMsg := "Invalid command"

	text := strings.Trim(dm.Text, " ")
	if text == errorMsg {
		return
	}

	cmds := []LineCommand{
		NewHitomiPreviewCommand(r.Accessor),
		NewStatusCommand(),
		NewSentryCommand(),
	}
	helpCmd := NewHelpCommand(cmds)

	executeCount := 0
	for _, cmd := range cmds {
		if cmd.match(text) {
			cmd.execute(text, r.StatusSender)
			executeCount++
		}
	}
	if helpCmd.match(text) {
		helpCmd.execute(text, r.StatusSender)
		executeCount++
	}

	if executeCount == 0 {
		r.StatusSender.SendTitleOnly(errorMsg)
	}
}

// 간단한 텍스트 기반의 명령어
type LineCommand interface {
	match(text string) bool
	execute(text string, sender *senders.Sender)
	help() string
}

type HelpCommand struct {
	title string
	cmds  []LineCommand
}

func NewHelpCommand(cmds []LineCommand) LineCommand {
	return &HelpCommand{
		title: "help",
		cmds:  cmds,
	}
}
func (c *HelpCommand) help() string {
	return "help"
}
func (c *HelpCommand) match(text string) bool {
	return text == "help"
}
func (c *HelpCommand) execute(text string, sender *senders.Sender) {
	lines := make([]string, len(c.cmds))
	for i, cmd := range c.cmds {
		lines[i] = cmd.help()
	}
	help := strings.Join(lines, "\n")
	sender.Send(c.title, help)
}

type HitomiPreviewCommand struct {
	title    string
	accessor storages.Accessor
}

func NewHitomiPreviewCommand(accessor storages.Accessor) LineCommand {
	return &HitomiPreviewCommand{
		title:    "Hitomi Preview",
		accessor: accessor,
	}
}

var reHitomiPreview = regexp.MustCompile(`^hitomi preview (\d{6})$`)

func (c *HitomiPreviewCommand) match(text string) bool {
	return reHitomiPreview.MatchString(text)
}
func (c *HitomiPreviewCommand) execute(text string, sender *senders.Sender) {
	title := c.title
	for _, m := range reHitomiPreview.FindAllStringSubmatch(text, -1) {
		code := m[1]
		log.Printf("DM: %s %s\n", title, code)

		go func() {
			ok := hitomiwatcher.FetchPreview(code, nil, c.accessor)
			if ok {
				sender.Send(title, fmt.Sprintf("success : %s", code))
			} else {
				sender.Send(title, fmt.Sprintf("fail : %s", code))
			}
		}()
	}
}

func (c *HitomiPreviewCommand) help() string {
	return "hitomi preview 123456"
}

type StatusCommand struct {
	title string
	cmd   string
}

func NewStatusCommand() LineCommand {
	return &StatusCommand{
		title: "status",
		cmd:   "check status",
	}
}

func (c *StatusCommand) help() string {
	return c.cmd
}
func (c *StatusCommand) execute(text string, sender *senders.Sender) {
	title := c.title

	now := time.Now()
	msg := now.Format(time.RFC3339)
	sender.Send(title, "still alive : "+msg)
}
func (c *StatusCommand) match(text string) bool {
	return text == c.cmd
}

type SentryCommand struct {
	title string
	cmd   string
}

func NewSentryCommand() LineCommand {
	return &SentryCommand{
		title: "sentry",
		cmd:   "check sentry",
	}
}
func (c *SentryCommand) help() string {
	return c.cmd
}
func (c *SentryCommand) execute(text string, sender *senders.Sender) {
	_, err := os.Open("invalid-file-to-raise-error")
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		sender.Send(c.title, "sentry success")
	} else {
		sender.Send(c.title, "sentry fail")
	}
}

func (c *SentryCommand) match(text string) bool {
	return text == c.cmd
}
