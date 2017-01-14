package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"io"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
	"github.com/kardianos/osext"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
	yaml "gopkg.in/yaml.v2"
)

type Storage struct {
	strategy storageStrategy
}

func NewStorage(rootpath string, token string) *Storage {
	if rootpath == "" || token == "" {
		return &Storage{newLocal()}
	}
	return &Storage{newDropbox(rootpath, token)}
}
func (s *Storage) UploadYaml(data interface{}, dst string) error {
	d, err := yaml.Marshal(data)
	check(err)
	return s.UploadBytes(d, dst)
}
func (s *Storage) UploadBytes(data []byte, dst string) error {
	return s.strategy.UploadBytes(data, dst)
}
func (s *Storage) UploadFile(src string, dst string) error {
	return s.strategy.UploadFile(src, dst)
}
func (s *Storage) Mkdir(dirname string) {
	s.strategy.Mkdir(dirname)
}

type MetaUploadResponse struct {
	ID       string
	FileName string
}

func (s *Storage) UploadMetadata(t *anaconda.Tweet, dir string, now time.Time) (MetaUploadResponse, error) {
	id := MakeOriginIdStr(t)
	filename := MakeTweetFileName(id, now, ".yaml")
	filename = path.Join(dir, filename)

	simpleTweet := NewSimpleTweet(t)
	e := s.UploadYaml(simpleTweet, filename)
	resp := MetaUploadResponse{
		ID:       id,
		FileName: filename,
	}
	return resp, e
}

type mediaResponse struct {
	Data     []byte
	FileName string
}

func (s *Storage) ArchiveTweet(tweet *anaconda.Tweet, dir string) {
	mediaCount := len(tweet.ExtendedEntities.Media)

	mediaRespChannel := make(chan *mediaResponse, mediaCount)
	for _, media := range tweet.ExtendedEntities.Media {
		go func(tweet *anaconda.Tweet, media anaconda.EntityMedia, resps chan<- *mediaResponse) {
			url := FindMediaURL(media)
			filename := MakeMediaFileName(tweet, media)

			resp, _ := http.Get(url)
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			resps <- &mediaResponse{
				body,
				filename,
			}
		}(tweet, media, mediaRespChannel)
	}

	mediaRespList := make([]*mediaResponse, mediaCount)
	for i := 0; i < mediaCount; i++ {
		mediaRespList[i] = <-mediaRespChannel
	}

	now := time.Now()
	id := MakeOriginIdStr(tweet)

	resp, e := s.UploadMetadata(tweet, dir, now)
	if e != nil {
		raven.CaptureErrorAndWait(e, nil)
		log.Panicf("Save Tweet Fail! %s -> %s, [%s]", resp.ID, resp.FileName, e.Error())
	} else {
		log.Printf("Save Tweet %s -> %s", resp.ID, resp.FileName)
	}

	// upload media
	for _, resp := range mediaRespList {
		filename := MakeNormalFileName(resp.FileName, now)
		filename = path.Join(dir, filename)
		err := s.UploadBytes(resp.Data, filename)
		if err != nil {
			raven.CaptureErrorAndWait(e, nil)
			log.Panicf("Save Image Fail! %s -> %s, [%s]", id, filename, err.Error())
		} else {
			log.Printf("Save Image %s -> %s", id, filename)
		}
	}
}

type storageStrategy interface {
	UploadBytes(data []byte, dst string) error
	UploadFile(src string, dst string) error
	Mkdir(dirname string)
}

type localStorageStrategy struct {
	RootPath string
}

func newLocal() storageStrategy {
	rootpath, _ := osext.ExecutableFolder()
	return &localStorageStrategy{
		RootPath: rootpath,
	}
}

func (c *localStorageStrategy) UploadBytes(data []byte, dst string) error {
	filepath := path.Join(c.RootPath, dst)
	dir := path.Dir(filepath)
	c.Mkdir(dir)

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	file.Write(data)
	defer file.Close()
	return err
}

func (c *localStorageStrategy) UploadFile(src string, dst string) error {
	panic("not implemented")
}

func (c *localStorageStrategy) Mkdir(dirname string) {
	os.MkdirAll(dirname, 0644)
}

type dropboxStorageStrategy struct {
	RootPath string
	client   *dropy.Client
	token    string
}

func newDropbox(rootpath string, token string) storageStrategy {
	client := dropy.New(dropbox.New(dropbox.NewConfig(token)))
	return &dropboxStorageStrategy{
		RootPath: rootpath,
		client:   client,
		token:    token,
	}
}

func (c *dropboxStorageStrategy) now() time.Time {
	localnow := time.Now()
	utcnow := localnow.UTC()
	t := time.Date(utcnow.Year(), utcnow.Month(), utcnow.Day(), utcnow.Hour(), utcnow.Minute(), utcnow.Second(), 0, time.UTC)
	return t
}

func (c *dropboxStorageStrategy) UploadBytes(data []byte, dst string) error {
	r := bytes.NewReader(data)
	return c.UploadReader(r, dst)
}

func (c *dropboxStorageStrategy) UploadReader(r io.Reader, dst string) error {
	uploadFilePath := path.Join(c.RootPath, dst)
	_, err := c.client.Files.Upload(&dropbox.UploadInput{
		Mode:           dropbox.WriteModeAdd,
		Path:           uploadFilePath,
		Reader:         r,
		Mute:           true,
		ClientModified: c.now(),
	})
	return err
}

func (c *dropboxStorageStrategy) UploadFile(src string, dst string) error {
	file, _ := os.Open(src)
	r := bufio.NewReader(file)
	return c.UploadReader(r, dst)
}

func (c *dropboxStorageStrategy) Mkdir(dirname string) {
	client := c.client
	_, err := client.Stat(dirname)
	if err != nil {
		client.Mkdir(dirname)
	}
}
