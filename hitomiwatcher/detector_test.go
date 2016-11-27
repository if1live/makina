package hitomiwatcher

import (
	"reflect"
	"testing"
	"time"
)

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

		{"이거 4728472782번 들었는데", nil},
		{"이거 123456번 들었는데", nil},

		{"abc 123456시간 테스트", nil},

		// black list
		{"131072", nil},
		{"234567", nil},
		{"「導」70%「嘘」15%「夜」15% ポイント:95pt ランキング:366545位", nil},

		// mention + hash
		{"@0c442e114489450", nil},
		{"@123456", nil},
		{"@123456_", nil},
		{"@_123456", nil},
		{"@_123456_", nil},
		{"#hash_161030", nil},

		// url
		{"http://foo.bar/299292", nil},

		{"https://hitomi.la/galleries/1.html", []int{1}},
		{"https://hitomi.la/galleries/123456.html", []int{123456}},
		{"https://hitomi.la/galleries/1234567.html", []int{1234567}},

		{"https://hitomi.la/reader/1.html", []int{1}},
		{"https://hitomi.la/reader/123456.html", []int{123456}},
		{"https://hitomi.la/reader/1234567.html", []int{1234567}},

		// 시간으로는 쓰일지 모른다
		{"분으로는 100800분 초로는 604800초", nil},
		{"650000개", nil},
		{"115000원", nil},
		{"！720710点", nil},

		{"110-111-111111 신한은행입니다.", nil},
		{"하나은행도 있어요. 109 111111 11111", nil},

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
			t.Errorf("FindReaderNumber - input : %s, expected %d, got %d, [%s]", c.text, c.num, actual, c.text)
		}
	}
}
