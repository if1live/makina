package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/haru/hitomi"
)

// 썸네일 얻기. hitomi api 에서 획득 가능하겠지?
func peekCoverImageUrls(code string) []string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}

	h := hitomi.New(client)
	metadata, _ := h.Metadata(code)
	return metadata.Covers
}

func FetchHitomiPreview(code string, tweet *anaconda.Tweet, storage *Storage) bool {
	coverUrls := peekCoverImageUrls(code)
	if len(coverUrls) == 0 {
		return false
	}

	now := time.Now()

	// 커버는 1개일 확률이 높으니까 고루틴 굳이 쓸 필요 없을거다
	for i, url := range coverUrls {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		prefix := MakePrefix(now)
		filename := ""
		ext := filepath.Ext(url)
		if len(coverUrls) == 1 {
			filename = fmt.Sprintf("%s%s", code, ext)
		} else {
			num := i + 1
			filename = fmt.Sprintf("%s_%d%s", code, num, ext)
		}
		storage.UploadBytes(body, prefix+filename)
	}

	if tweet != nil && len(coverUrls) > 0 {
		storage.UploadMetadata(tweet, "", now)
	}
	return true
}
