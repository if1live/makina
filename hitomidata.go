package main

import "regexp"

var notAllowedPrefixList = []string{
	"0",

	// 특수문자
	"-",
	"_",
	"+",
	"\\",
	"/",

	// 통화
	"₩",
	"$",
	"円",
}

/* 대소문자 무시 */
var notAllowedSuffixList = []string{
	// 자주 쓰일거같은 postfix
	"번",
	"회",
	"명",
	"밖",
	"개",
	"배",
	"권",
	"짜리",
	"알",
	"gtb",

	"에",
	"로",
	"으로",
	"이었",
	"이엇",
	"였",

	"트윗",
	"블락",
	"팔로",
	"정도",
	"rt",
	"알티",

	"점",
	"위",
	"pt",
	"렙",
	"뎀",
	"레벨",
	"포인트",

	// 자주 쓰이는 외국어
	"点",
	"個",
	"回",
	"人",
	"位",
	"日",
	"の",
	"フ",
	"α",
	"β",

	// 거리
	"cm",
	"m",
	"km",

	// 시간
	"시",
	"초",
	"분",
	"시",
	"개월",

	// 금액
	"원",
	"엔",
	"달러",

	// 글쎼?
	"세대",
	"가루",

	// 특수문자
	"%",
	"_",
	"-",
	"+",
	"/",
}

var predefinedBlackListKeyword = []string{
	// 은행이라는 단어가 등장하면 계좌번호겠지?
	// 등장했던거 위주로 보충해나가자
	"은행",
	"신한",
	"하나",
	"국민",

	"일러스트",
	"팝픈",
	"주식",
	"해금",
	"빌드",
	"체력",
	"방어",
	"test",
	"테스트",
}

var notAllowedCodes = []string{
	// 2**n 중 6자리 숫자 제외
	"131072",
	"262144",
	"524288",
}

var reHitomiGallery = regexp.MustCompile(`/galleries/(\d+).html`)
var reHitomiReader = regexp.MustCompile(`/reader/(\d+).html`)

var simpleIgnoreReList = []*regexp.Regexp{
	regexp.MustCompile(`@.*`),
	regexp.MustCompile(`#.*`),
}
