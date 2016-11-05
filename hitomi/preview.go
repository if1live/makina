package hitomi

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/network"
	"github.com/if1live/makina/twutils"
)

func ExtractCoverImages(jsontext string) []string {
	// golang으로 json을 파싱하는건 귀찮다
	// 필요한건 preview뿐이니까 정규식으로 찾아내는게 가능하겠다
	// json을 한줄로 만들어야 정규식 검색이 잘된다
	jsontext = strings.Replace(jsontext, "\n", "", -1)
	jsontext = strings.Replace(jsontext, "\t", "", -1)

	reCover := regexp.MustCompile(`"covers":\s*\[([^\]]*)\],`)
	for _, m := range reCover.FindAllStringSubmatch(jsontext, -1) {
		line := m[1]
		line = strings.Trim(line, " ")
		urls := strings.Split(line, ",")
		for i := 0; i < len(urls); i++ {
			url := urls[i]
			url = strings.Replace(url, `"`, "", -1)
			url = strings.Replace(url, " ", "", -1)
			urls[i] = url
		}

		return urls
	}

	// else
	return []string{}
}

// 썸네일 얻기. hitomi api 에서 획득 가능하겠지?
func PeekCoverImageUrls(code int, host string) []string {
	url := fmt.Sprintf("http://%s/api/detail/hitomi/%d", host, code)
	fetcher := network.HttpFetcher{}
	apiResp := fetcher.Fetch(url)
	jsontext := string(apiResp.Data)
	urls := ExtractCoverImages(jsontext)
	return urls
}

func FetchPreview(code int, tweet *anaconda.Tweet, config Config) {
	coverUrls := PeekCoverImageUrls(code, config.HaruHostName)
	now := time.Now()

	// 커버는 1개일 확률이 높으니까 고루틴 굳이 쓸 필요 없을거다
	for i, url := range coverUrls {
		fetcher := network.HttpFetcher{}
		resp := fetcher.Fetch(url)

		prefix := twutils.MakePrefix(now)
		filename := ""
		ext := filepath.Ext(url)
		if len(coverUrls) == 1 {
			filename = fmt.Sprintf("%d%s", code, ext)
		} else {
			num := i + 1
			filename = fmt.Sprintf("%d_%d%s", code, num, ext)
		}
		config.Accessor.UploadBytes(resp.Data, prefix+filename)
	}

	if tweet != nil && len(coverUrls) > 0 {
		twutils.UploadMetadata(tweet, config.Accessor, "", now)
	}
}
