package main

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

	lowertext := strings.ToLower(text)

	// 금지단어가 포함되어있는 경우
	for _, w := range predefinedBlackListKeyword {
		if strings.Contains(lowertext, w) {
			return nil
		}
	}

	// 공백문자로 쪼갠후 검사. 6자리 숫자는 한 단어 들어갈테니까
	// https://play.golang.org/p/cLHpRxZQiG
	words := strings.FieldsFunc(lowertext, func(r rune) bool {
		switch r {
		case ' ', '\n', '\t':
			return true
		}
		return false
	})

	codes := []int{}
	for _, word := range words {
		if code := findReaderNumberFromWord(word, now); code != notFound {
			codes = append(codes, code)
		}
	}

	if len(codes) == 0 {
		return nil
	}
	return codes
}

func filterBlacklist(word string, blacklist []string) bool {
	for _, b := range blacklist {
		if word == b {
			return false
		}
	}
	return true
}

func filterFullDate(word string) bool {
	// 20170102 규격에 맞아 떨어지면 무시
	// 무시할 년도 미리 지정. 설마 makina가 5년씩 돌아가겠어?
	ignoreYears := []int{
		2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022,
	}
	t, err := time.Parse("20060102", word)
	if err != nil {
		return true
	}
	for _, year := range ignoreYears {
		if t.Year() == year {
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

// 실제 히토미 코드는 한자리수도 존재하지만 5자리 코드부터 허용
// 작은 숫자를 허용하면 얻는것에 비해서 쓸데없는 메세지도 코드로 인식할거같아서
var reValidCode = regexp.MustCompile(`([^1-9]*)([1-9]\d{4,})([^0-9]*)`)

func findReaderNumberFromWord(word string, now time.Time) int {
	// 5자리 숫자까진 허용
	if len(word) < 5 {
		return notFound
	}

	// url로 추정?
	url, err := url.Parse(word)
	if err == nil && url.Host != "" {
		if url.Host == "hitomi.la" {
			if m := reHitomiGallery.FindStringSubmatch(word); len(m) > 0 {
				val, _ := strconv.Atoi(m[1])
				return val
			}
			if m := reHitomiReader.FindStringSubmatch(word); len(m) > 0 {
				val, _ := strconv.Atoi(m[1])
				return val
			}
		}
		// else..
		return notFound
	}

	// 간단한 정규식으로 걸러낼수 있는거
	for _, re := range simpleIgnoreReList {
		if m := re.MatchString(word); m {
			return notFound
		}
	}

	// 금지어 목록 이용
	for _, bw := range predefinedBlackListKeyword {
		if strings.Contains(word, bw) {
			return notFound
		}
	}

	for _, m := range reValidCode.FindAllStringSubmatch(word, -1) {
		prefix := m[1]
		s := m[2]
		suffix := m[3]

		for _, notAllowedPrefix := range notAllowedPrefixList {
			if strings.HasSuffix(prefix, notAllowedPrefix) {
				return notFound
			}
		}
		for _, notAllowedSuffix := range notAllowedSuffixList {
			if strings.HasPrefix(suffix, notAllowedSuffix) {
				return notFound
			}
		}

		for _, notAllowed := range notAllowedCodes {
			if s == notAllowed {
				return notFound
			}
		}
		if ok := filterFullDate(s); !ok {
			return notFound
		}

		if ok := filterRecentDate(s, now); !ok {
			return notFound
		}

		//fmt.Println(prefix, s, suffix)
		val, _ := strconv.Atoi(s)
		return val
	}

	return notFound
}
