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
}

func (ar *MediaArchiver) OnEvent(ev string, event *anaconda.EventTweet) {
	// Event list
	// reference: https://dev.twitter.com/docs/streaming-apis/messages#User_stream_messages
	switch ev {
	case "favorite":
		ar.OnFavorite(event)
	case "retweeted_retweet":
		ar.OnRetweet(event)
	case "favorited_retweet":
		ar.OnRetweet(event)
	}
}
func (ar *MediaArchiver) OnRetweet(tweet *anaconda.EventTweet) {
	// 내가 RT한것만 저장
	if tweet.Source.ScreenName != ar.myName {
		return
	}

	t := tweet.TargetObject
	id := twutils.ProfitIdStr(t)
	log.Printf("retweet : %s, %s\n", id, t.Text)
	ar.handleTweet(t, ar.accessor, "retweet")
	log.Printf("MediaArchiver Complete %s", id)
}

func (ar *MediaArchiver) OnFavorite(tweet *anaconda.EventTweet) {
	// 남이 favorite한것도 이벤트로 들어오더라. 그래서 무시
	if tweet.Source.ScreenName != ar.myName {
		return
	}
	t := tweet.TargetObject
	id := twutils.ProfitIdStr(t)
	log.Printf("favorite : %s, %s\n", id, t.Text)
	ar.handleTweet(t, ar.accessor, "favorite")
	log.Printf("MediaArchiver Complete %s", id)
}

func (ar *MediaArchiver) handleTweet(tweet *anaconda.Tweet, accessor storages.Accessor, dir string) {
	if tweet == nil {
		return
	}
	if len(tweet.ExtendedEntities.Media) == 0 {
		return
	}

	category := ""

	for _, user := range ar.predefinedUsers {
		// 트위터 계정명은 대소문자를 구분하지 않더라
		s1 := strings.ToLower(user)
		s2 := strings.ToLower(twutils.ProfitScreenName(tweet))
		if s1 == s2 {
			category = "user-" + user
			break
		}
	}

	if category != "" {
		twutils.ArchiveMedia(tweet, accessor, category)
	} else {
		twutils.ArchiveMedia(tweet, accessor, dir)
	}
}
