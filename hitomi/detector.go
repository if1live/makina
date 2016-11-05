package hitomi

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	notFound = -1
)

func FindReaderNumber(text string, now time.Time) int {
	if len(text) == 0 {
		return notFound
	}

	// 금지단어가 포함되어있는 경우
	for _, w := range predefinedBlackListKeyword {
		if strings.Contains(text, w) {
			return notFound
		}
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
	regexp.MustCompile(`\d{6}시간`),
	regexp.MustCompile(`\d{6}개`),
	regexp.MustCompile(`\d{6}cm`),
	regexp.MustCompile(`\d{6}m`),
	regexp.MustCompile(`\d{6}km`),
}

var reGallery = regexp.MustCompile(`/galleries/(\d{6}).html`)
var reReader = regexp.MustCompile(`/reader/(\d{6}).html`)
var reValidCode = regexp.MustCompile(`([1-9]\d{5})`)

var predefinedBlackList = []string{
	// 2**n 중 6자리 숫자 제외
	"131072",
	"262144",
	"524288",

	// 연속된 숫자는 수동으로 입력했을 가능성이 높다
	// 몇개안되니까 하드코딩. 123456은 테스트에서 자주 써서 예외처리
	"234567",
	"345678",
	"456789",
	"567890",
}
var predefinedBlackListKeyword = []string{
	// 은행이라는 단어가 등장하면 계좌번호겠지?
	// 등장했던거 위주로 보충해나가자
	"은행",
	"신한",
	"하나",
	"국민",

	// 설마 이런 문자가 끼겠어?
	"%",
}

func filterBlacklist(word string, blacklist []string) bool {
	for _, b := range blacklist {
		if word == b {
			return false
		}
	}
	return true
}

func findReaderNumberFromText(word string, now time.Time) int {
	blacklist := make([]string, len(predefinedBlackList))
	copy(blacklist, predefinedBlackList)

	// 오늘 +-N일 제외
	// 약간 크게 잡는게 좋을거같다
	for i := -60; i <= 60; i++ {
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
