package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	code := d.FindReaderNumber(text, time.Now())
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

func (d *HitomiDetector) FindReaderNumber(text string, now time.Time) int {
	if len(text) == 0 {
		return notFound
	}

	// 공백문자로 쪼갠후 검사. 6자리 숫자는 한 단어 들어갈테니까
	// https://play.golang.org/p/cLHpRxZQiG
	words := strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case ' ', '\n', '\t':
			return true
		}
		return false
	})

	for _, word := range words {
		if code := findReaderNumberFromText(word, now); code != notFound {
			return code
		}
	}

	return notFound
}

func filterBlacklist(word string, blacklist []string) bool {
	for _, b := range blacklist {
		if word == b {
			return false
		}
	}
	return true
}

// 7자리 숫자 or 5자리 이하 숫자는 패스
// 그러면 6자리 관련 검사가 간단해진다
var simpleIgnoreReList = []*regexp.Regexp{
	regexp.MustCompile(`\d{7,}`),
	regexp.MustCompile(`@.*\d{6}.*`),
	regexp.MustCompile(`#.*\d{6}.*`),

	// 자주 쓰일거같은 postfix
	regexp.MustCompile(`\d{6}초`),
	regexp.MustCompile(`\d{6}분`),
	regexp.MustCompile(`\d{6}시`),
	regexp.MustCompile(`\d{6}번`),
	regexp.MustCompile(`\d{6}cm`),
	regexp.MustCompile(`\d{6}m`),
	regexp.MustCompile(`\d{6}km`),
}

var reGallery = regexp.MustCompile(`/galleries/(\d{6}).html`)
var reReader = regexp.MustCompile(`/reader/(\d{6}).html`)
var reValidCode = regexp.MustCompile(`([1-9]\d{5})`)

func findReaderNumberFromText(word string, now time.Time) int {
	// 오늘 +-3일 제외
	blacklist := []string{}
	for i := -3; i <= 3; i++ {
		duration := time.Hour * time.Duration(24*i)
		t := now.Add(duration)
		datestr := t.Format("060102")
		blacklist = append(blacklist, datestr)
	}
	if len(word) < 6 {
		return notFound
	}

	if len(word) == 6 {
		if ok := filterBlacklist(word, blacklist); !ok {
			return notFound
		}
		if m := reValidCode.MatchString(word); m {
			val, _ := strconv.Atoi(word)
			return val
		}
		return notFound
	}

	for _, re := range simpleIgnoreReList {
		if m := re.MatchString(word); m {
			return notFound
		}
	}

	// url로 추정?
	url, err := url.Parse(word)
	if err == nil && url.Host != "" {
		if url.Host == "hitomi.la" {
			if m := reGallery.FindStringSubmatch(word); len(m) > 0 {
				val, _ := strconv.Atoi(m[1])
				return val
			}
			if m := reReader.FindStringSubmatch(word); len(m) > 0 {
				val, _ := strconv.Atoi(m[1])
				return val
			}
		}
		// else..
		return notFound
	}

	for _, m := range reValidCode.FindAllStringSubmatch(word, -1) {
		if ok := filterBlacklist(m[1], blacklist); !ok {
			return notFound
		}

		val, _ := strconv.Atoi(m[1])
		return val
	}

	return notFound
}
