package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"

	"log"

	"path"

	"github.com/ChimeraCoder/anaconda"
)

var cmd string
var logfilename string

func init() {
	flag.StringVar(&cmd, "cmd", "", "command")
	flag.StringVar(&logfilename, "log", "", "log filename")
}

func main() {
	flag.Parse()

	// initialize logger
	// http: //stackoverflow.com/questions/19965795/go-golang-write-log-to-file
	// logger 초기화를 별도 함수에서 할 경우 defer 로 파일이 닫혀서 로그작성이 안된다
	// 그래서 그냥 메인함수에서 처리
	if logfilename != "" {
		filepath := path.Join(GetExecutablePath(), logfilename)
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

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

	ignorableEvents := []string{
		"favorited_retweet",
		"retweeted_retweet",
	}

	api := config.NewDataSourceAuthConfig().CreateApi()
	twitterStream := api.UserStream(nil)
	for {
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			for _, h := range handlers {
				go h.OnTweet(&tweet)
			}
		case anaconda.StatusDeletionNotice:
			// pass
		case anaconda.FriendsList:
		case anaconda.EventTweet:
			evt := tweet.Event.Event
			switch evt {
			case "favorite":
				for _, h := range handlers {
					go h.OnFavorite(&tweet)
				}
			case "unfavorite":
				for _, h := range handlers {
					go h.OnUnfavorite(&tweet)
				}
			default:
				ignorable := false
				for _, evtname := range ignorableEvents {
					if evtname == evt {
						ignorable = true
						break
					}
				}

				if ignorable {
					log.Printf("event = %s : skip\n", evt)
				} else {
					log.Printf("unknown event(%T) : %v \n", x, evt)
				}
			}
		default:
			log.Printf("unknown type(%T) : %v \n", x, x)
		}
	}
}
