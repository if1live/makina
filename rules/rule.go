package rules

import "github.com/ChimeraCoder/anaconda"

type Rule interface {
	OnTweet(tweet *anaconda.Tweet)
	OnEvent(ev string, event *anaconda.EventTweet)
	OnDirectMessage(dm *anaconda.DirectMessage)
}
