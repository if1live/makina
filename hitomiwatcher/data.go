package hitomiwatcher

import "regexp"

var notAllowedPrefixList = []string{
	"0",

	// 특수문자
	"-",
	"+",
	"$",
	"\\",
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
	"GTB",

	// 게임용어?
	"점",
	"위",
	"pt",
	"렙",
	"레벨",

	// 자주 쓰이는 한자
	"点",
	"個",
	"回",
	"人",
	"位",

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
}

var predefinedBlackListKeyword = []string{
	// 은행이라는 단어가 등장하면 계좌번호겠지?
	// 등장했던거 위주로 보충해나가자
	"은행",
	"신한",
	"하나",
	"국민",
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
