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
	d.Handle(tweet)
}

func (d *HitomiWatcher) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.MyName {
		return
	}

	t := tweet.TargetObject
	d.Handle(t)
}
func (d *HitomiWatcher) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}

func (d *HitomiWatcher) Handle(tweet *anaconda.Tweet) bool {
	codes := FindReaderNumbers(tweet.Text, time.Now())
	if len(codes) == 0 {
		return false
	}

	id := MakeOriginIdStr(tweet)
	for _, code := range codes {
		codestr := strconv.Itoa(code)
		log.Printf("Hitomi Found Code %d, %s", code, id)
		success := FetchHitomiPreview(codestr, tweet, d.storage)
		if success {
			log.Printf("Hitomi Fetch Preview Complete %s", id)

		} else {
			// 디버깅 목적으로 업로드 검색 실패해도 업로드는 하기
			d.storage.UploadMetadata(tweet, "hitomi-preview-fail", time.Now())
			log.Printf("Hitomi Fetch Preview Failed %s", id)
		}
	}
	return true
}
