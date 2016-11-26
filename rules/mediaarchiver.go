package rules

import (
	"log"

	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/storages"
	"github.com/if1live/makina/twutils"
)

type MediaArchiver struct {
	accessor        storages.Accessor
	myName          string
	predefinedUsers []string
}

func NewMediaArchiver(accessor storages.Accessor, myName string, users []string) TweetRule {
	archiver := &MediaArchiver{
		accessor:        accessor,
		myName:          myName,
		predefinedUsers: users,
	}
	return archiver
}

func (ar *MediaArchiver) OnTweet(tweet *anaconda.Tweet) {
	if tweet.RetweetedStatus == nil {
		return
	}

	if tweet.User.ScreenName == ar.myName && tweet.RetweetedStatus.User.ScreenName != ar.myName {
		ar.handleTweet(tweet, ar.accessor, "media-retweet")
	}
}

func (ar *MediaArchiver) OnEvent(ev string, event *anaconda.EventTweet) {
	// Event list
	// reference: https://dev.twitter.com/docs/streaming-apis/messages#User_stream_messages
	switch ev {
	case "favorite":
		ar.OnFavorite(event)
		break
	}
}

func (ar *MediaArchiver) OnFavorite(tweet *anaconda.EventTweet) {
	// 남이 favorite한것도 이벤트로 들어오더라. 그래서 무시
	if tweet.Source.ScreenName != ar.myName {
		return
	}
	t := tweet.TargetObject
	ar.handleTweet(t, ar.accessor, "media-favorite")
}

func (ar *MediaArchiver) handleTweet(tweet *anaconda.Tweet, accessor storages.Accessor, dir string) {
	if tweet == nil {
		return
	}
	if len(tweet.ExtendedEntities.Media) == 0 {
		return
	}

	// 기본 카테고리는 tweet or favorite
	// 더 좋은 카테고리가 있을떄 교체하는 식으로 동작한다
	category := dir

	for _, user := range ar.predefinedUsers {
		// 트위터 계정명은 대소문자를 구분하지 않더라
		s1 := strings.ToLower(user)
		s2 := strings.ToLower(twutils.ProfitScreenName(tweet))
		if s1 == s2 {
			category = "user-" + user
			break
		}
	}

	id := twutils.ProfitIdStr(tweet)
	log.Printf("Media Archive %s : %s, %s\n", category, id, tweet.Text)
	twutils.ArchiveMedia(tweet, accessor, category)
	log.Printf("MediaArchiver Complete %s", id)
}
