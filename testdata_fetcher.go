package main

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
)

func MainFetchTestData(config *Config) {
	api := config.NewDataSourceAuthConfig().CreateApi()
	tweet := fetchSampleTweet(api)
	SaveTweetJsonFile(&tweet, "testdata/sample-tweet.json")
	log.Println("Fetch test data success.")
}

func fetchSampleTweet(api *anaconda.TwitterApi) anaconda.Tweet {
	const tweetId = 303777106620452864
	tweet, err := api.GetTweet(tweetId, nil)
	check(err)
	return tweet
}
