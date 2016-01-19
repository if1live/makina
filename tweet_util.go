package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/drhodes/golorem"
)

type TweetUtil struct {
	Api *anaconda.TwitterApi

	sample *anaconda.Tweet
}

func (util *TweetUtil) FetchSampleTweet() anaconda.Tweet {
	const tweetId = 303777106620452864
	tweet, err := util.Api.GetTweet(tweetId, nil)
	check(err)
	return tweet
}

func (util *TweetUtil) SaveJsonFile(t *anaconda.Tweet, filepath string) {
	b, err := json.Marshal(t)
	check(err)

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	f, err := os.Create(filepath)
	check(err)

	w := bufio.NewWriter(f)
	out.WriteTo(w)

	w.Flush()
}

func (util *TweetUtil) LoadJsonFile(filepath string) anaconda.Tweet {
	// https://gist.github.com/border/775526
	file, e := ioutil.ReadFile(filepath)
	check(e)

	tweet := anaconda.Tweet{}
	json.Unmarshal(file, &tweet)
	return tweet
}

func (util *TweetUtil) SampleTweet() anaconda.Tweet {
	const filepath = "./sample-tweet.json"

	if util.sample == nil {
		t := util.LoadJsonFile(filepath)
		util.sample = &t
	}
	return *util.sample
}

// 테스트 목적으로 가짜 트윗 찍어내기
func (util *TweetUtil) RandomTweet() anaconda.Tweet {
	t := util.SampleTweet()
	t.Id = rand.Int63()
	t.Text = lorem.Sentence(2, 10)
	return t
}
