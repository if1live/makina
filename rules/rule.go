package rules

import "github.com/ChimeraCoder/anaconda"

type TweetRule interface {
	OnTweet(tweet *anaconda.Tweet)
	OnEvent(ev string, event *anaconda.EventTweet)
}

type MessageRule interface {
	OnDirectMessage(dm *anaconda.DirectMessage)
}

type DaemonRule interface {
	Execute()
}
