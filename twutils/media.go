package twutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		id := ProfitIdStr(tweet)
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
	id := ProfitIdStr(tweet)
	return fmt.Sprintf("%s_%d%s", id, num, ext)
}

func MakePrefix(now time.Time) string {
	return now.Format("20060102-150405-")
}
func MakeTweetFileName(id string, now time.Time, ext string) string {
	prefix := MakePrefix(now)
	filename := prefix + id + ext
	return filename
}
func MakeNormalFileName(filename string, now time.Time) string {
	return MakePrefix(now) + filename
}

type UploadMetadataResponse struct {
	ID       string
	FileName string
}

func ProfitIdStr(t *anaconda.Tweet) string {
	id := t.IdStr
	if t.RetweetedStatus != nil {
		id = t.RetweetedStatus.IdStr
	}
	return id
}

// TODO 하위 폴더로 정리해서 업로드하는게 필요해질지 모른다
func UploadFullMetadata(t *anaconda.Tweet, accessor storages.Accessor, path string, now time.Time) (UploadMetadataResponse, error) {
	id := ProfitIdStr(t)
	filename := MakeTweetFileName(id, now, ".json")
	e := accessor.UploadJson(t, filename)
	resp := UploadMetadataResponse{
		ID:       id,
		FileName: filename,
	}
	return resp, e
}

func UploadMetadata(t *anaconda.Tweet, accessor storages.Accessor, path string, now time.Time) (UploadMetadataResponse, error) {
	id := ProfitIdStr(t)
	filename := MakeTweetFileName(id, now, ".yaml")
	simpleTweet := NewSimpleTweet(t)
	e := accessor.UploadYaml(simpleTweet, filename)
	resp := UploadMetadataResponse{
		ID:       id,
		FileName: filename,
	}
	return resp, e
}

type MediaResponse struct {
	Data     []byte
	FileName string
}

func ArchiveMedia(tweet *anaconda.Tweet, accessor storages.Accessor) {
	mediaCount := len(tweet.ExtendedEntities.Media)

	mediaRespChannel := make(chan *MediaResponse, mediaCount)
	for _, media := range tweet.ExtendedEntities.Media {
		go fetchMediaCh(tweet, media, mediaRespChannel)
	}

	mediaRespList := make([]*MediaResponse, mediaCount)
	for i := 0; i < mediaCount; i++ {
		mediaRespList[i] = <-mediaRespChannel
	}

	save(tweet, mediaRespList, accessor)
}

func fetchMediaCh(tweet *anaconda.Tweet, media anaconda.EntityMedia, resps chan<- *MediaResponse) {
	url := FindMediaURL(media)
	filename := MakeMediaFileName(tweet, media)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	resps <- &MediaResponse{
		body,
		filename,
	}
}

func save(tweet *anaconda.Tweet, mediaRespList []*MediaResponse, accessor storages.Accessor) {
	now := time.Now()
	id := ProfitIdStr(tweet)

	resp, e := UploadMetadata(tweet, accessor, "", now)
	if e != nil {
		log.Fatalf("Save Tweet Fail! %s -> %s, [%s]", resp.ID, resp.FileName, e.Error())
	} else {
		log.Printf("Save Tweet %s -> %s", resp.ID, resp.FileName)
	}

	// upload media
	for _, resp := range mediaRespList {
		filename := MakeNormalFileName(resp.FileName, now)
		err := accessor.UploadBytes(resp.Data, filename)
		if err != nil {
			log.Fatalf("Save Image Fail! %s -> %s, [%s]", id, filename, err.Error())
		} else {
			log.Printf("Save Image %s -> %s", id, filename)
		}
	}
}
