package main

import "testing"

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
		// 0으로 시작하진 않을거다
		{"012345", -1},
		// 후보가 여러개면 첫번째
		{"123456 234567", 123456},
		{"012345 123456 234567", 123456},

		// 사용된적이 있는 텍스트
		{"123456 펫이야 펫!!", 123456},
		{"123456 미즈가센세 너무 좋아", 123456},
	}

	detector := HitomiDetector{nil}
	for _, c := range cases {
		actual := detector.FindReaderNumber(c.text)
		if c.num != actual {
			t.Errorf("FindReaderNumber - expected %d, got %d, [%s]", c.num, actual, c.text)
		}
	}
}
