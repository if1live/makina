package rules

import (
	"log"
	"time"

	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/hitomiwatcher"
	"github.com/if1live/makina/storages"
	"github.com/if1live/makina/twutils"
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

	id := twutils.ProfitIdStr(tweet)
	codestr := strconv.Itoa(code)
	log.Printf("Hitomi Found Code %d, %s", code, id)
	hitomiwatcher.FetchPreview(codestr, tweet, d.Accessor)
	log.Printf("Hitomi Fetch Preview Complete %s", id)
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

	id := twutils.ProfitIdStr(t)
	codestr := strconv.Itoa(code)
	log.Printf("Hitomi Found Code %d, %s", code, id)
	hitomiwatcher.FetchPreview(codestr, t, d.Accessor)
	log.Printf("Hitomi Fetch Preview Complete %s", id)
}
func (d *HitomiWatcher) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}
