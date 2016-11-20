package rules

import (
	"log"
	"regexp"

	"time"

	"fmt"

	"strings"

	"github.com/ChimeraCoder/anaconda"
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

	success := false
	success = success || r.hitomiPreview("hitomi preview", text)
	success = success || r.status("status", text)
	if !success {
		r.StatusSender.SendTitleOnly(errorMsg)
	}
}

var reHitomiPreview = regexp.MustCompile(`^hitomi preview (\d{6})$`)

func (r *DirectMessageWatcher) hitomiPreview(title, text string) bool {
	for _, m := range reHitomiPreview.FindAllStringSubmatch(text, -1) {
		code := m[1]
		log.Printf("DM: %s %s\n", title, code)

		go func() {
			ok := hitomiwatcher.FetchPreview(code, nil, r.Accessor)
			if ok {
				r.StatusSender.Send(title, fmt.Sprintf("success : %s", code))
			} else {
				r.StatusSender.Send(title, fmt.Sprintf("fail : %s", code))
			}
		}()
		return true
	}

	return false
}

func (r *DirectMessageWatcher) status(title, text string) bool {
	if text == "status" {
		log.Printf("DM: %s\n", title)

		now := time.Now()
		msg := now.Format(time.RFC3339)
		r.StatusSender.Send(title, "still alive : "+msg)

		return true
	}

	return false
}
