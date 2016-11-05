package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/hitomi"
	"github.com/if1live/makina/storages"
)

const (
	notFound           = -1
	hitomiSavePath     = "/hitomi-temp"
	haruPath           = "../haru"
	haruExecutableName = "haru"
)

type HitomiDetector struct {
	config   *Config
	accessor storages.Accessor
}

func NewHitomiDetector(config *Config) *HitomiDetector {
	detector := &HitomiDetector{
		config,
		config.NewStorageAccessor(hitomiSavePath),
	}
	return detector
}

func (d *HitomiDetector) OnTweet(tweet *anaconda.Tweet) {
	d.ProcessText(tweet.Text, tweet.IdStr)
}
func (d *HitomiDetector) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.config.DataSourceScreenName {
		return
	}
	d.ProcessText(tweet.TargetObject.Text, tweet.TargetObject.IdStr)
}
func (d *HitomiDetector) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}

func (d *HitomiDetector) ProcessText(text string, tweetId string) {
	code := hitomi.FindReaderNumber(text, time.Now())
	if code == notFound {
		return
	}

	log.Printf("HitomiDetector Found Code %d", code)

	config := hitomi.Config{
		ExecutablePath: haruPath,
		ExecutableName: haruExecutableName,
		ShowLog:        true,
	}
	success, zipfilename := hitomi.ExecuteHaru(config, code)
	if success {
		log.Printf("HitomiDetector Haru Complete %s, %s", zipfilename, tweetId)
		// upload
		baseZipFileName := filepath.Base(zipfilename)
		c := storages.NewDropbox(hitomiSavePath, d.config.DropboxAccessToken)
		c.UploadFile(zipfilename, baseZipFileName)

	} else {
		log.Printf("HitomiDetector Haru Failed %s", tweetId)
	}
	log.Printf("HitomiDetector Complete %s", tweetId)
}
