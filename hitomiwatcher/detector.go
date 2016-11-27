package hitomiwatcher

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

func FindReaderNumbers(text string, now time.Time) []int {
	if len(text) == 0 {
		return nil
	}

	// 금지단어가 포함되어있는 경우
	for _, w := range predefinedBlackListKeyword {
		if strings.Contains(text, w) {
			return nil
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

	codes := []int{}
	for _, word := range words {
		if code := findReaderNumberFromText(word, now); code != notFound {
			codes = append(codes, code)
		}
	}

	if len(codes) == 0 {
		return nil
	}
	return codes
}

// 7자리 숫자 or 5자리 이하 숫자는 패스
// 그러면 6자리 관련 검사가 간단해진다
var simpleIgnoreReList = []*regexp.Regexp{
	// 히토미 코드 백만 진입
	// 그래서 7자리 숫자도 허용해야한다
	// 아니, 길이 제한이 의미 있을까?
	//regexp.MustCompile(`\d{7,}`),
	regexp.MustCompile(`@.*\d+.*`),
	regexp.MustCompile(`#.*\d+.*`),

	// 자주 쓰일거같은 postfix
	regexp.MustCompile(`\d+초`),
	regexp.MustCompile(`\d+분`),
	regexp.MustCompile(`\d+시`),
	regexp.MustCompile(`\d+번`),
	regexp.MustCompile(`\d+명`),
	regexp.MustCompile(`\d+원`),
	regexp.MustCompile(`\d+시간`),
	regexp.MustCompile(`\d+개`),
	regexp.MustCompile(`\d+점`),
	regexp.MustCompile(`\d+cm`),
	regexp.MustCompile(`\d+m`),
	regexp.MustCompile(`\d+km`),
	regexp.MustCompile(`\d+pt`),
	regexp.MustCompile(`\d+点`),
	regexp.MustCompile(`\d+GTB`),
}

var reGallery = regexp.MustCompile(`/galleries/(\d+).html`)
var reReader = regexp.MustCompile(`/reader/(\d+).html`)

// 실제 히토미 코드는 한자리수도 존재하지만 5자리 코드부터 허용
// 작은 숫자를 허용하면 얻는것에 비해서 쓸데없는 메세지도 코드로 인식할거같아서
var reValidCode = regexp.MustCompile(`(\d{5,})`)

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

func filterRecentDate(word string, now time.Time) bool {
	t, err := time.Parse("060102", word)
	if err != nil {
		return true
	}

	// 올해, 작년 날짜로 추정되면 무시
	invalidYears := []int{
		now.Year() - 1,
		now.Year(),
	}
	for _, year := range invalidYears {
		if t.Year() == year {
			return false
		}
	}
	return true
}

func findReaderNumberFromText(word string, now time.Time) int {
	blacklist := make([]string, len(predefinedBlackList))
	copy(blacklist, predefinedBlackList)

	// 5자리 숫자까진 허용
	if len(word) < 5 {
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
		s := m[1]
		if ok := filterBlacklist(s, blacklist); !ok {
			return notFound
		}
		if ok := filterRecentDate(s, now); !ok {
			return notFound
		}
		// 숫자의 맨 앞자리가 0인것은 제외
		// 정규식으로 짜르려고했는데 잘 안되서 그냥 후처리로 대응했다
		// 문제있던거 : 012345 를 [1-9]\d+ 로 잡으면 12345가 잡히더라
		if s[0] == '0' {
			return notFound
		}

		val, _ := strconv.Atoi(s)
		return val
	}

	return notFound
}
