package main

import (
	"flag"

	"log"

	"github.com/ChimeraCoder/anaconda"
)

var cmd string

func init() {
	flag.StringVar(&cmd, "cmd", "", "command")
}

func main() {
	flag.Parse()
	config := LoadConfig()

	switch cmd {
	case "fetch_testdata":
		MainFetchTestData(config)
	case "":
		mainDefault(config)
	default:
		log.Fatalf("unknown command")
	}
}

type StreamingHandler interface {
	OnTweet(tweet *anaconda.Tweet)
	OnFavorite(tweet *anaconda.EventTweet)
	OnUnfavorite(twee *anaconda.EventTweet)
}

func mainDefault(config *Config) {
	go mainServer(config)
	mainStreaming(config)
}

func mainServer(config *Config) {

}

func mainStreaming(config *Config) {
	handlers := []StreamingHandler{
		NewFavoriteMediaArchiver(config),
	}

	api := config.NewDataSourceAuthConfig().CreateApi()
	twitterStream := api.UserStream(nil)
	for {
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			for _, h := range handlers {
				h.OnTweet(&tweet)
			}
		case anaconda.StatusDeletionNotice:
			// pass
		case anaconda.FriendsList:
		case anaconda.EventTweet:
			evt := tweet.Event.Event
			switch evt {
			case "favorite":
				for _, h := range handlers {
					h.OnFavorite(&tweet)
				}
			case "unfavorite":
				for _, h := range handlers {
					h.OnUnfavorite(&tweet)
				}
			case "favorited_retweet":
				log.Println("favorited_retweet : skip")
			case "retweeted_retweet":
				log.Println("retweeted_retweet : skip")
			default:
				log.Printf("unknown event(%T) : %v \n", x, evt)
			}
		default:
			log.Printf("unknown type(%T) : %v \n", x, x)
		}
	}
}
