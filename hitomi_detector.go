package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"path"

	"github.com/ChimeraCoder/anaconda"
)

const (
	NotFound = -1
)

type HitomiDetector struct {
	config *Config
}

func NewHitomiDetector(config *Config) *HitomiDetector {
	detector := &HitomiDetector{
		config,
	}
	return detector
}

func (d *HitomiDetector) OnTweet(tweet *anaconda.Tweet) {
	d.ProcessText(tweet.Text)
}
func (d *HitomiDetector) OnFavorite(tweet *anaconda.EventTweet) {
	if tweet.Source.ScreenName != d.config.DataSourceScreenName {
		return
	}
	d.ProcessText(tweet.TargetObject.Text)
}
func (d *HitomiDetector) OnUnfavorite(tweet *anaconda.EventTweet) {
}

func (d *HitomiDetector) ProcessText(text string) {
	code := d.FindReaderNumber(text)
	if code == NotFound {
		return
	}

	currDir := GetExecutablePath()
	haruPath := path.Join(currDir, "..", "haru")
	haruExeName := "haru"
	haruFilePath := path.Join(haruPath, haruExeName)
	idParam := fmt.Sprintf("-id=%d", code)
	out, err := exec.Command(haruFilePath, idParam, "-service=hitomi", "-cmd=download").CombinedOutput()

	elems := []string{}
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			elems = append(elems, err.Error())
		}
	}
	elems = strings.Split(string(out[:]), "\n")

	for _, line := range elems {
		fmt.Println(line)
	}

}

func (d *HitomiDetector) FindReaderNumber(text string) int {
	if len(text) == 0 {
		return NotFound
	}

	if len(text) == 6 {
		m, _ := regexp.MatchString(`([1-9]\d{5})`, text)
		if m {
			val, _ := strconv.Atoi(text)
			return val
		}
		return NotFound
	}

	// 7글자 이상이면 7자리 숫자일떄의 예외처리가 필요
	re := regexp.MustCompile(`([1-9]\d{5})[^0-9]`)
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		val, _ := strconv.Atoi(m[1])
		return val
	}
	return NotFound
}
