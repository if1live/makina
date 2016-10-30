package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"path/filepath"

	"path"

	"github.com/ChimeraCoder/anaconda"
	dropy "github.com/tj/go-dropy"
)

const (
	savePath = "/archive-temp"
)

type FavoriteImageArchiver struct {
	config *Config
	client *dropy.Client
}

func NewFavoriteImageArchiver(config *Config) *FavoriteImageArchiver {
	client := config.CreateDropboxClient()
	archiver := &FavoriteImageArchiver{
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

func (ar *FavoriteImageArchiver) OnTweet(tweet *anaconda.Tweet) {
}

func (ar *FavoriteImageArchiver) OnUnfavorite(tweet *anaconda.EventTweet) {
}

type MediaResponse struct {
	Response *FetchResponse
	FileName string
}

func fetchMediaCh(tweet *anaconda.Tweet, idx int, media anaconda.EntityMedia, resps chan<- *MediaResponse) {
	url := media.Media_url

	num := idx + 1
	ext := filepath.Ext(url)
	filename := fmt.Sprintf("%s_%d%s", tweet.IdStr, num, ext)

	fetcher := HttpFetcher{}
	resp := fetcher.Fetch(url)

	resps <- &MediaResponse{
		resp,
		filename,
	}
}

func (ar *FavoriteImageArchiver) OnFavorite(tweet *anaconda.EventTweet) {
	t := tweet.TargetObject
	log.Printf("favorite : %s, %s\n", t.IdStr, t.Text)

	if len(t.ExtendedEntities.Media) == 0 {
		log.Printf("No media attached, skip")
		return
	}

	// 로컬에 저장. 트윗당 이미지는 최대 4개로 제한된다
	// 그래서 코루틴 만들어서 돌려도 특별한 문제 없다
	mediaRespChannel := make(chan *MediaResponse, 4)
	for idx, media := range t.ExtendedEntities.Media {
		fetchMediaCh(t, idx, media, mediaRespChannel)
	}

	mediaCount := len(t.ExtendedEntities.Media)
	mediaRespList := make([]*MediaResponse, mediaCount)
	for i := 0; i < mediaCount; i++ {
		mediaRespList[i] = <-mediaRespChannel
	}

	//ar.saveLocal(t, mediaRespList)
	ar.saveDropbox(t, mediaRespList)
	log.Printf("Complete %s", t.IdStr)
}

func (ar *FavoriteImageArchiver) saveLocal(tweet *anaconda.Tweet, mediaRespList []*MediaResponse) {
	jsonFilename := tweet.IdStr + ".json"
	SaveTweetJsonFile(tweet, jsonFilename)
	log.Printf("Save Tweet %s  ->%s", tweet.IdStr, jsonFilename)

	for _, resp := range mediaRespList {
		file, err := os.OpenFile(resp.FileName, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		file.Write(resp.Response.Data)
		log.Printf("Save Image %s -> %s", tweet.IdStr, resp.FileName)
		file.Close()
	}
}

func (ar *FavoriteImageArchiver) saveDropbox(tweet *anaconda.Tweet, mediaRespList []*MediaResponse) {
	c := ar.client

	// upload tweet metadata
	jsonFilename := tweet.IdStr + ".json"
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
		r := bytes.NewReader(resp.Response.Data)
		uploadFilePath := path.Join(savePath, resp.FileName)
		err := c.Upload(uploadFilePath, r)
		if err != nil {
			log.Fatalf("Upload Image Fail! %s -> %s, [%s]", tweet.IdStr, resp.FileName, err.Error())
		} else {
			log.Printf("Upload Image %s -> %s", tweet.IdStr, resp.FileName)
		}
	}
}
