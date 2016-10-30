package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"

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
	type Status struct {
		Ok bool `json:"ok"`
	}
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := Status{
			true,
		}
		data, _ := json.Marshal(status)
		var out bytes.Buffer
		json.Indent(&out, data, "", "  ")
		w.Write(out.Bytes())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server started: http://0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
