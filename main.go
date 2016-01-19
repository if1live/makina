package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func FetchFavorites(api *anaconda.TwitterApi, screenName string) {
	v := url.Values{}
	// maximum count: 200
	// count value for test: 5
	v.Set("count", "1")
	v.Set("screen_name", screenName)

	result, _ := api.GetFavorites(v)
	for _, tweet := range result {
		// TODO archive tweet
		fmt.Println(tweet.Text)
		fmt.Printf("%#v", tweet)
	}
}

func main() {
	api := DefaultAuthConfig().CreateApi()
	util := TweetUtil{api, nil}

	if os.Getenv("FETCH_SAMPLE") == "1" {
		// sample tweet는 한번만 만들면 충분하다
		// 이후에는 디버깅/테스트 목적으로 쓰인다
		sampleTweet := util.FetchSampleTweet()
		util.SaveJsonFile(&sampleTweet, "./sample-tweet.json")
		fmt.Println("Fetch sample tweet : success")
		os.Exit(0)
	}

	for i := 0; i < 2; i++ {
		t := util.RandomTweet()
		fmt.Println(t.Id, t.Text)
	}

	// TODO favorites 를 받는대로 어딘가에 저장하기
}
