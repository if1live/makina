package twutils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/storages"
)

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
		return fmt.Sprintf("%s%s", tweet.IdStr, ext)
	}

	found := -1
	for i := 0; i < mediaCount; i++ {
		m := tweet.ExtendedEntities.Media[i]
		if m.Display_url == media.Display_url {
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
	return fmt.Sprintf("%s_%d%s", tweet.IdStr, num, ext)
}

func MakePrefix(now time.Time) string {
	return now.Format("20060102-150405-")
}
func MakeTweetFileName(id string, now time.Time) string {
	prefix := MakePrefix(now)
	filename := prefix + id + ".json"
	return filename
}
func MakeNormalFileName(filename string, now time.Time) string {
	return MakePrefix(now) + filename
}

type UploadMetadataResponse struct {
	ID       string
	FileName string
}

// TODO 하위 폴더로 정리해서 업로드하는게 필요해질지 모른다
func UploadMetadata(t *anaconda.Tweet, accessor storages.Accessor, path string, now time.Time) (UploadMetadataResponse, error) {
	id := t.IdStr
	filename := MakeTweetFileName(id, now)
	e := accessor.UploadJson(t, filename)
	resp := UploadMetadataResponse{
		ID:       id,
		FileName: filename,
	}
	return resp, e
}
