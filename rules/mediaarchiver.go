package rules

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/storages"
	"github.com/if1live/makina/twutils"
)

type MediaArchiver struct {
	accessor storages.Accessor
	myName   string
}

func NewMediaArchiver(accessor storages.Accessor, myName string) TweetRule {
	archiver := &MediaArchiver{
		accessor: accessor,
		myName:   myName,
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
	log.Printf("retweet : %s, %s\n", t.IdStr, t.Text)
	handleTweet(t, ar.accessor)
	log.Printf("MediaArchiver Complete %s", t.IdStr)
}

func (ar *MediaArchiver) OnFavorite(tweet *anaconda.EventTweet) {
	// 남이 favorite한것도 이벤트로 들어오더라. 그래서 무시
	if tweet.Source.ScreenName != ar.myName {
		return
	}
	t := tweet.TargetObject
	log.Printf("favorite : %s, %s\n", t.IdStr, t.Text)
	handleTweet(t, ar.accessor)
	log.Printf("MediaArchiver Complete %s", t.IdStr)
}

func handleTweet(tweet *anaconda.Tweet, accessor storages.Accessor) {
	if tweet == nil {
		return
	}
	if len(tweet.ExtendedEntities.Media) == 0 {
		return
	}
	twutils.ArchiveMedia(tweet, accessor)
}
