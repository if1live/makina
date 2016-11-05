package media_archiver

import (
	"log"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/network"
	"github.com/if1live/makina/storages"
	"github.com/if1live/makina/twutils"
)

type MediaResponse struct {
	Response *network.FetchResponse
	FileName string
}

func Archive(tweet *anaconda.Tweet, accessor storages.Accessor) {
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
	url := twutils.FindMediaURL(media)
	filename := twutils.MakeMediaFileName(tweet, media)
	fetcher := network.HttpFetcher{}
	resp := fetcher.Fetch(url)

	resps <- &MediaResponse{
		resp,
		filename,
	}
}

func save(tweet *anaconda.Tweet, mediaRespList []*MediaResponse, accessor storages.Accessor) {
	now := time.Now()

	resp, e := twutils.UploadMetadata(tweet, accessor, "", now)
	if e != nil {
		log.Fatalf("Save Tweet Fail! %s -> %s, [%s]", resp.ID, resp.FileName, e.Error())
	} else {
		log.Printf("Save Tweet %s -> %s", resp.ID, resp.FileName)
	}

	// upload media
	for _, resp := range mediaRespList {
		filename := twutils.MakeNormalFileName(resp.FileName, now)
		err := accessor.UploadBytes(resp.Response.Data, filename)
		if err != nil {
			log.Fatalf("Save Image Fail! %s -> %s, [%s]", tweet.IdStr, filename, err.Error())
		} else {
			log.Printf("Save Image %s -> %s", tweet.IdStr, filename)
		}
	}
}
