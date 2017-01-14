package main

import (
	"fmt"
	"path/filepath"
	"time"

	"strings"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
	"github.com/kardianos/osext"
)

func check(e error) {
	if e != nil {
		raven.CaptureErrorAndWait(e, nil)
		panic(e)
	}
}

func GetExecutablePath() string {
	path, err := osext.ExecutableFolder()
	check(err)
	return path
}

// 파일명을 만들때 시간을 이용하면 시간순정렬을 쓸수있다
// TODO 파일 수정시간을 손보는게 더 좋지 않을까?
func MakePrefix(now time.Time) string {
	return now.Format("20060102-150405-")
}
func MakeTweetFileName(id string, now time.Time, ext string) string {
	tokens := []string{
		//MakePrefix(now),
		id,
		ext,
	}
	return strings.Join(tokens, "")
}
func MakeNormalFileName(filename string, now time.Time) string {
	tokens := []string{
		//MakePrefix(now),
		filename,
	}
	return strings.Join(tokens, "")
}

// 리트윗의 경우 원형을 찾아서 저장하고싶다
func MakeOriginIdStr(t *anaconda.Tweet) string {
	if t.RetweetedStatus != nil {
		return MakeOriginIdStr(t.RetweetedStatus)
	}
	return t.IdStr
}
func MakeOriginScreenName(t *anaconda.Tweet) string {
	if t.RetweetedStatus != nil {
		return MakeOriginScreenName(t.RetweetedStatus)
	}
	return t.User.ScreenName
}

func findURLFromVideo(media anaconda.EntityMedia) string {
	maxBitrate := -1
	selectedVariant := anaconda.Variant{}
	for _, v := range media.VideoInfo.Variants {
		if v.Bitrate > maxBitrate {
			maxBitrate = v.Bitrate
			selectedVariant = v
		}
	}
	return selectedVariant.Url
}

func findURLFromPhoto(media anaconda.EntityMedia) string {
	return media.Media_url
}

func findURLFromAnimatedGif(media anaconda.EntityMedia) string {
	// 사실상 비디오랑 같은 취급
	return findURLFromVideo(media)
}

func FindMediaURL(media anaconda.EntityMedia) string {
	switch media.Type {
	case "video":
		return findURLFromVideo(media)
	case "animated_gif":
		return findURLFromAnimatedGif(media)
	case "photo":
		return findURLFromPhoto(media)
	default:
		return findURLFromPhoto(media)
	}
}

// 트윗에 붙은 이미지가 여러개인 경우와 한개인 경우를 구분
func MakeMediaFileName(tweet *anaconda.Tweet, media anaconda.EntityMedia) string {
	mediaCount := len(tweet.ExtendedEntities.Media)
	if mediaCount <= 1 {
		url := FindMediaURL(media)
		ext := filepath.Ext(url)
		id := MakeOriginIdStr(tweet)
		return fmt.Sprintf("%s%s", id, ext)
	}

	found := -1
	for i := 0; i < mediaCount; i++ {
		m := tweet.ExtendedEntities.Media[i]
		if m.Media_url == media.Media_url {
			found = i
			break
		}
	}

	if found < 0 {
		// not found
		return ""
	}

	num := found + 1
	url := FindMediaURL(media)
	ext := filepath.Ext(url)
	id := MakeOriginIdStr(tweet)
	return fmt.Sprintf("%s_%d%s", id, num, ext)
}
