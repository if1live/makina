package twutils

import "github.com/ChimeraCoder/anaconda"

func FetchTweet(id int64, api *anaconda.TwitterApi) anaconda.Tweet {
	tweet, err := api.GetTweet(id, nil)
	if err != nil {
		panic(err)
	}
	return tweet
}
