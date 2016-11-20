package rules

import (
	"log"
	"time"

	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/hitomiwatcher"
	"github.com/if1live/makina/storages"
)

type HitomiWatcher struct {
	MyName   string
	Accessor storages.Accessor
}

func NewHitomiWatcher(myName string, accessor storages.Accessor) TweetRule {
	detector := &HitomiWatcher{
		MyName:   myName,
		Accessor: accessor,
	}
	return detector
}

func (d *HitomiWatcher) OnTweet(tweet *anaconda.Tweet) {
	code := hitomiwatcher.FindReaderNumber(tweet.Text, time.Now())
	if code < 0 {
		return
	}

	codestr := strconv.Itoa(code)
	log.Printf("Hitomi Found Code %d, %s", code, tweet.IdStr)
	hitomiwatcher.FetchPreview(codestr, tweet, d.Accessor)
	log.Printf("Hitomi Fetch Preview Complete %s", tweet.IdStr)
}

func (d *HitomiWatcher) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.MyName {
		return
	}

	t := tweet.TargetObject
	code := hitomiwatcher.FindReaderNumber(t.Text, time.Now())
	if code < 0 {
		return
	}

	codestr := strconv.Itoa(code)
	log.Printf("Hitomi Found Code %d, %s", code, t.IdStr)
	hitomiwatcher.FetchPreview(codestr, t, d.Accessor)
	log.Printf("Hitomi Fetch Preview Complete %s", t.IdStr)
}
func (d *HitomiWatcher) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}
