package main

import (
	"log"
	"time"

	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

type HitomiWatcher struct {
	MyName  string
	storage *Storage
}

func NewHitomiWatcher(myName string, storage *Storage) TweetRule {
	detector := &HitomiWatcher{
		MyName:  myName,
		storage: storage,
	}
	return detector
}

func (d *HitomiWatcher) OnTweet(tweet *anaconda.Tweet) {
	codes := FindReaderNumbers(tweet.Text, time.Now())
	if len(codes) == 0 {
		return
	}

	id := MakeOriginIdStr(tweet)
	for _, code := range codes {
		codestr := strconv.Itoa(code)
		log.Printf("Hitomi Found Code %d, %s", code, id)
		FetchHitomiPreview(codestr, tweet, d.storage)
		log.Printf("Hitomi Fetch Preview Complete %s", id)
	}
}

func (d *HitomiWatcher) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.MyName {
		return
	}

	t := tweet.TargetObject
	codes := FindReaderNumbers(t.Text, time.Now())
	if len(codes) == 0 {
		return
	}

	id := MakeOriginIdStr(t)
	for _, code := range codes {
		codestr := strconv.Itoa(code)
		log.Printf("Hitomi Found Code %d, %s", code, id)
		FetchHitomiPreview(codestr, t, d.storage)
		log.Printf("Hitomi Fetch Preview Complete %s", id)
	}
}
func (d *HitomiWatcher) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}
