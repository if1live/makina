package rules

import (
	"log"
	"regexp"

	"time"

	"fmt"

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

	text := dm.Text
	r.hitomiPreview(text)
	r.status(text)
}

var reHitomiPreview = regexp.MustCompile(`^hitomi preview (\d{6})$`)

func (r *DirectMessageWatcher) hitomiPreview(text string) {
	for _, m := range reHitomiPreview.FindAllStringSubmatch(text, -1) {
		code := m[1]
		log.Printf("DM: hitomi preview %s\n", code)
		ok := hitomiwatcher.FetchPreview(code, nil, r.Accessor)
		if ok {
			r.StatusSender.Send("hitomi preview", fmt.Sprintf("success : %s", code))
		} else {
			r.StatusSender.Send("hitomi preview", fmt.Sprintf("fail : %s", code))
		}
	}
}

func (r *DirectMessageWatcher) status(text string) {
	if text == "status" {
		log.Println("DM: status")

		now := time.Now()
		msg := now.Format(time.RFC3339)
		r.StatusSender.Send("still alive", msg)
	}
}
