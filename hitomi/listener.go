package hitomi

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type Listener struct {
	config Config
}

func NewListener(config Config) *Listener {
	detector := &Listener{
		config: config,
	}
	return detector
}

func (d *Listener) OnTweet(tweet *anaconda.Tweet) {
	code := FindReaderNumber(tweet.Text, time.Now())
	if code == notFound {
		return
	}

	log.Printf("Hitomi Found Code %d, %s", code, tweet.IdStr)
	FetchPreview(code, tweet, d.config)
	log.Printf("Hitomi Fetch Preview Complete %s", tweet.IdStr)
}

func (d *Listener) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.config.MyName {
		return
	}

	t := tweet.TargetObject
	code := FindReaderNumber(t.Text, time.Now())
	if code == notFound {
		return
	}

	log.Printf("Hitomi Found Code %d, %s", code, t.IdStr)
	FetchPreview(code, t, d.config)
	log.Printf("Hitomi Fetch Preview Complete %s", t.IdStr)
}
func (d *Listener) OnEvent(ev string, event *anaconda.EventTweet) {
	switch ev {
	case "favorite":
		d.OnFavorite(event)
	}
}

func (d *Listener) OnDirectMessage(dm *anaconda.DirectMessage) {
	if dm.Sender.ScreenName != d.config.MyName {
		return
	}

	text := dm.Text

	reFull := regexp.MustCompile(`hitomi (\d{6})`)
	for _, m := range reFull.FindAllStringSubmatch(text, -1) {
		code, _ := strconv.Atoi(m[1])
		FetchFull(code, d.config)
	}

	rePreview := regexp.MustCompile(`hitomi preview (\d{6})`)
	for _, m := range rePreview.FindAllStringSubmatch(text, -1) {
		code, _ := strconv.Atoi(m[1])
		FetchPreview(code, nil, d.config)
	}
}

func FetchFull(code int, config Config) {
	success, zipfilename := ExecuteHaru(config, code)
	if success {
		log.Printf("Haru Complete %s", zipfilename)
		// upload
		baseZipFileName := filepath.Base(zipfilename)
		config.Accessor.UploadFile(zipfilename, baseZipFileName)
		log.Printf("Haru Upload %s", baseZipFileName)

		body := fmt.Sprintf("Fetch(%d) %s success", code, baseZipFileName)
		config.StatusSender.Send("Makina-Hitomi", body)

	} else {
		body := fmt.Sprintf("Fetch(%d) failed", code)
		log.Printf(body)
		config.StatusSender.Send("Makina-Hitomi", body)
	}
}