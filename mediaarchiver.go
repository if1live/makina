package main

import (
	"log"

	"strings"

	"github.com/ChimeraCoder/anaconda"
)

type MediaArchiver struct {
	storage         *Storage
	myName          string
	predefinedUsers []string
}

func NewMediaArchiver(storage *Storage, myName string, users []string) TweetRule {
	return &MediaArchiver{
		storage:         storage,
		myName:          myName,
		predefinedUsers: users,
	}
}

func (ar *MediaArchiver) OnTweet(tweet *anaconda.Tweet) {
	if tweet.RetweetedStatus == nil {
		return
	}

	if tweet.User.ScreenName == ar.myName {
		ar.handleTweet(tweet, ar.storage, "media-retweet")
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
	ar.handleTweet(t, ar.storage, "media-favorite")
}

func (ar *MediaArchiver) handleTweet(tweet *anaconda.Tweet, storage *Storage, dir string) {
	if !ar.SaveRequired(tweet) {
		return
	}

	writerScreenName := MakeOriginScreenName(tweet)
	category := ar.FindCategory(dir, writerScreenName)

	id := MakeOriginIdStr(tweet)
	log.Printf("Media Archive %s : %s, %s\n", category, id, tweet.Text)
	storage.ArchiveTweet(tweet, category)
	log.Printf("MediaArchiver Complete %s", id)
}

func (ar *MediaArchiver) SaveRequired(tweet *anaconda.Tweet) bool {
	if tweet == nil {
		return false
	}
	if len(tweet.ExtendedEntities.Media) == 0 {
		return false
	}
	return true
}
func (ar *MediaArchiver) FindCategory(dir, writerScreenName string) string {
	// 기본 카테고리는 tweet or favorite
	// 더 좋은 카테고리가 있을떄 교체하는 식으로 동작한다
	category := dir

	name := strings.ToLower(writerScreenName)
	for _, user := range ar.predefinedUsers {
		// 트위터 계정명은 대소문자를 구분하지 않더라
		s1 := strings.ToLower(user)
		if s1 == name {
			category = "user-" + user
			break
		}
	}

	return category
}
