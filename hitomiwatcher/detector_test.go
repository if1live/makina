package hitomiwatcher

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFindReaderNumber_Skip(t *testing.T) {
	filename := "testcase-skip.txt"
	dat, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(dat), "\n")
	for i, line := range lines {
		lines[i] = strings.Replace(line, "\r", "", -1)
	}
	now := time.Now()
	for _, c := range lines {
		actual := FindReaderNumbers(c, now)
		if actual != nil {
			t.Errorf("TestFindReaderNumber_Skip - input : %s, expected nil, got %d", c, actual)
		}
	}
}

func TestFindReaderNumber(t *testing.T) {
	cases := []struct {
		text string
		num  []int
	}{
		// 공백 문자열 처리
		{"", nil},

		{"123456", []int{123456}},
		{"abcdef", nil},

		// 5자리 자리 숫자 이상부터 취급
		// 낮은 숫자는 검색으로 얻는것에 비해서 노이즈로 문제 생길 가능성이 높다
		{"1", nil},
		{"12", nil},
		{"123", nil},
		{"1234", nil},
		{"12345", []int{12345}},
		{"1234567", []int{1234567}},
		{"12345678", []int{12345678}},

		// 0으로 시작하진 않을거다
		{"012345", nil},

		// 후보가 여러개면 전부
		{"123456 654321", []int{123456, 654321}},

		// 다양한 분리문자
		{"111111\t111112\n111113 111114", []int{111111, 111112, 111113, 111114}},

		// 사용된적이 있는 텍스트
		{"123456 펫이야 펫!!", []int{123456}},
		{"123456 미즈가센세 너무 좋아", []int{123456}},

		{"123456가", []int{123456}},
		{"가123456", []int{123456}},
		{"가123456가", []int{123456}},

		// black list
		{"131072", nil},

		// url
		{"http://foo.bar/299292", nil},

		{"https://hitomi.la/galleries/1.html", []int{1}},
		{"https://hitomi.la/galleries/123456.html", []int{123456}},
		{"https://hitomi.la/galleries/1234567.html", []int{1234567}},

		{"https://hitomi.la/reader/1.html", []int{1}},
		{"https://hitomi.la/reader/123456.html", []int{123456}},
		{"https://hitomi.la/reader/1234567.html", []int{1234567}},

		// 날짜처럼 생기면 무시
		// 1년에 365일, 최근 2년만 씹어도 별 문제 없을테니까
		// 매우 과거 또는 미래의 경우 날짜처럼 생겨도 허용,
		{"161030 Gimpo Airport", nil},
		{"161030", nil},
		{"171030", []int{171030}},
		{"141030", []int{141030}},
	}

	now := time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC)
	for _, c := range cases {
		actual := FindReaderNumbers(c.text, now)
		if !reflect.DeepEqual(c.num, actual) {
			t.Errorf("FindReaderNumber - input : %s, expected %d, got %d", c.text, c.num, actual)
		}
	}
}
