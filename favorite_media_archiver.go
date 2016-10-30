package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"path/filepath"

	"path"

	"github.com/ChimeraCoder/anaconda"
	dropy "github.com/tj/go-dropy"
)

const (
	savePath = "/archive-temp"
)

type FavoriteMediaArchiver struct {
	config *Config
	client *dropy.Client
}

func NewFavoriteMediaArchiver(config *Config) *FavoriteMediaArchiver {
	client := config.CreateDropboxClient()
	archiver := &FavoriteMediaArchiver{
		config,
		client,
	}

	// 이미지 저장 경로가 존재하는지 확인
	_, err := archiver.client.Stat(savePath)
	if err != nil {
		archiver.client.Mkdir(savePath)
	}

	return archiver
}

func (ar *FavoriteMediaArchiver) OnTweet(tweet *anaconda.Tweet) {
}

func (ar *FavoriteMediaArchiver) OnUnfavorite(tweet *anaconda.EventTweet) {
}

type MediaResponse struct {
	Response *FetchResponse
	FileName string
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

func fetchMediaCh(tweet *anaconda.Tweet, idx int, totalMediaCount int, media anaconda.EntityMedia, resps chan<- *MediaResponse) {
	url := ""
	if media.Type == "video" {
		url = findURLFromVideo(media)
	} else {
		url = findURLFromPhoto(media)
	}

	// 트윗에 붙은 이미지가 여러개인 경우와 한개인 경우를 구분
	filename := ""
	if totalMediaCount == 1 {
		ext := filepath.Ext(url)
		filename = fmt.Sprintf("%s%s", tweet.IdStr, ext)
	} else {
		num := idx + 1
		ext := filepath.Ext(url)
		filename = fmt.Sprintf("%s_%d%s", tweet.IdStr, num, ext)
	}

	fetcher := HttpFetcher{}
	resp := fetcher.Fetch(url)

	resps <- &MediaResponse{
		resp,
		filename,
	}
}

func (ar *FavoriteMediaArchiver) OnFavorite(tweet *anaconda.EventTweet) {
	t := tweet.TargetObject

	// 남이 favorite한것도 이벤트로 들어오더라. 그래서 무시
	if tweet.Source.ScreenName != ar.config.DataSourceScreenName {
		return
	}

	if len(t.ExtendedEntities.Media) == 0 {
		return
	}

	// 파일명에 prefix가 있으면 저장 시점을 파악하기 쉬울것이다
	log.Printf("favorite : %s, %s\n", t.IdStr, t.Text)
	now := time.Now()
	prefix := now.Format("20060102-150405-")

	mediaCount := len(t.ExtendedEntities.Media)

	mediaRespChannel := make(chan *MediaResponse, mediaCount)
	for idx, media := range t.ExtendedEntities.Media {
		go fetchMediaCh(t, idx, mediaCount, media, mediaRespChannel)
	}

	mediaRespList := make([]*MediaResponse, mediaCount)
	for i := 0; i < mediaCount; i++ {
		mediaRespList[i] = <-mediaRespChannel
	}

	//ar.saveLocal(t, mediaRespList, prefix)
	ar.saveDropbox(t, mediaRespList, prefix)
	log.Printf("FavoriteMediaArchiver Complete %s", t.IdStr)
}

func (ar *FavoriteMediaArchiver) saveLocal(tweet *anaconda.Tweet, mediaRespList []*MediaResponse, prefix string) {
	executablePath := GetExecutablePath()

	jsonFilename := prefix + tweet.IdStr + ".json"
	jsonFilePath := path.Join(executablePath, jsonFilename)
	SaveTweetJsonFile(tweet, jsonFilePath)
	log.Printf("Save Tweet %s  ->%s", tweet.IdStr, jsonFilename)

	for _, resp := range mediaRespList {
		filename := prefix + resp.FileName
		filepath := path.Join(executablePath, filename)
		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		file.Write(resp.Response.Data)
		log.Printf("Save Image %s -> %s", tweet.IdStr, filename)
		file.Close()
	}
}

func (ar *FavoriteMediaArchiver) saveDropbox(tweet *anaconda.Tweet, mediaRespList []*MediaResponse, prefix string) {
	c := ar.client

	// upload tweet metadata
	jsonFilename := prefix + tweet.IdStr + ".json"
	b, err := json.Marshal(tweet)
	check(err)
	var jsonOut bytes.Buffer
	json.Indent(&jsonOut, b, "", "  ")

	r := bytes.NewReader(jsonOut.Bytes())
	uploadFilePath := path.Join(savePath, jsonFilename)
	e := c.Upload(uploadFilePath, r)
	if e != nil {
		log.Fatalf("Upload Tweet Fail! %s -> %s, [%s]", tweet.IdStr, jsonFilename, err.Error())
	} else {
		log.Printf("Upload Tweet %s -> %s", tweet.IdStr, jsonFilename)
	}

	// upload media
	for _, resp := range mediaRespList {
		filename := prefix + resp.FileName
		r := bytes.NewReader(resp.Response.Data)
		uploadFilePath := path.Join(savePath, filename)
		err := c.Upload(uploadFilePath, r)
		if err != nil {
			log.Fatalf("Upload Image Fail! %s -> %s, [%s]", tweet.IdStr, filename, err.Error())
		} else {
			log.Printf("Upload Image %s -> %s", tweet.IdStr, filename)
		}
	}
}
