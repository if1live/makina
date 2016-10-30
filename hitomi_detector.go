package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	dropy "github.com/tj/go-dropy"
)

const (
	notFound           = -1
	hitomiSavePath     = "/hitomi-temp"
	haruPath           = "../haru"
	haruExecutableName = "haru"
)

type HitomiDetector struct {
	config *Config
	client *dropy.Client
}

func NewHitomiDetector(config *Config) *HitomiDetector {
	client := config.CreateDropboxClient()
	detector := &HitomiDetector{
		config,
		client,
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
func (d *HitomiDetector) OnUnfavorite(tweet *anaconda.EventTweet) {
}

func (d *HitomiDetector) ProcessText(text string, tweetId string) {
	code := d.FindReaderNumber(text)
	if code == notFound {
		return
	}

	log.Printf("HitomiDetector Found Code %d", code)

	currDir := GetExecutablePath()
	cmd := path.Join(currDir, haruPath, haruExecutableName)
	args := []string{
		fmt.Sprintf("-id=%d", code),
		"-service=hitomi",
		"-cmd=download",
	}
	out, err := exec.Command(cmd, args...).CombinedOutput()

	// dump stdout/stderr
	stderrs := []string{}
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			stderrs = append(stderrs, err.Error())
		}
	}
	stdouts := strings.Split(string(out[:]), "\n")

	for _, line := range stderrs {
		log.Println(line)
	}
	for _, line := range stdouts {
		log.Println(line)
	}

	// 공백이 아닌 가장 마지막 출력 찾기
	// 거기에 파일명이 있을거다
	lastStdout := ""
	for i := len(stdouts) - 1; i >= 0; i-- {
		if len(stdouts[i]) > 0 {
			lastStdout = stdouts[i]
			break
		}
	}

	zipFilename := ""
	re := regexp.MustCompile(` (/.*\.zip)`)
	for _, m := range re.FindAllStringSubmatch(lastStdout, -1) {
		zipFilename = m[1]
	}
	log.Printf("HitomiDetector Haru Complete %s, %s", zipFilename, tweetId)

	// upload
	baseZipFileName := filepath.Base(zipFilename)
	uploadFilePath := path.Join(hitomiSavePath, baseZipFileName)
	file, _ := os.Open(zipFilename)
	r := bufio.NewReader(file)
	d.client.Upload(uploadFilePath, r)

	log.Printf("HitomiDetector Complete %s", tweetId)
}

func (d *HitomiDetector) FindReaderNumber(text string) int {
	if len(text) == 0 {
		return notFound
	}

	if len(text) == 6 {
		m, _ := regexp.MatchString(`([1-9]\d{5})`, text)
		if m {
			val, _ := strconv.Atoi(text)
			return val
		}
		return notFound
	}

	// 7글자 이상이면 7자리 숫자일떄의 예외처리가 필요
	re := regexp.MustCompile(`([1-9]\d{5})[^0-9]`)
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		val, _ := strconv.Atoi(m[1])
		return val
	}
	return notFound
}
