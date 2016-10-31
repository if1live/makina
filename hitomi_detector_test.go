package main

import (
	"testing"
	"time"
)

func TestFindReaderNumber(t *testing.T) {
	cases := []struct {
		text string
		num  int
	}{
		// 공백 문자열 처리
		{"", -1},

		{"123456", 123456},
		{"abcdef", -1},

		// 6자리 아니면 무시
		{"12345", -1},
		{"1234567", -1},
		{"12345678", -1},
		// 0으로 시작하진 않을거다
		{"012345", -1},
		// 후보가 여러개면 첫번째
		{"123456 234567", 123456},
		{"012345 123456 234567", 123456},

		// 다양한 분리문자
		{"012345\t012345\n012345 123456", 123456},

		// 사용된적이 있는 텍스트
		{"123456 펫이야 펫!!", 123456},
		{"123456 미즈가센세 너무 좋아", 123456},

		{"123456가", 123456},
		{"가123456", 123456},
		{"가123456가", 123456},

		{"이거 4728472782번 들었는데", -1},
		{"이거 123456번 들었는데", -1},

		// mention + hash
		{"@0c442e114489450", -1},
		{"@123456", -1},
		{"@123456_", -1},
		{"@_123456", -1},
		{"@_123456_", -1},
		{"#hash_161030", -1},

		{"1. 762606000, in time_t", -1},

		// url
		{"http://foo.bar/299292", -1},
		{"https://hitomi.la/galleries/123456.html", 123456},
		{"https://hitomi.la/reader/123456.html", 123456},

		// 오늘 +-N일은 제외. 확률적으로 설마 이게 코드겠어?
		{"161030 Gimpo Airport", -1},

		// 시간으로는 쓰일지 모른다
		{"분으로는 100800분 초로는 604800초", -1},

		// TODO 계좌번호 포맷으로 추정될 경우?
	}

	detector := HitomiDetector{nil, nil}
	now := time.Date(2016, 10, 30, 0, 0, 0, 0, time.UTC)
	for _, c := range cases {
		actual := detector.FindReaderNumber(c.text, now)
		if c.num != actual {
			t.Errorf("FindReaderNumber - expected %d, got %d, [%s]", c.num, actual, c.text)
		}
	}
}
