package twutils

import (
	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
)

func FetchTweet(id int64, api *anaconda.TwitterApi) anaconda.Tweet {
	tweet, err := api.GetTweet(id, nil)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	return tweet
}
