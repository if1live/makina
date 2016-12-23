package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path"

	"log"

	"net/url"

	"time"

	"github.com/ChimeraCoder/anaconda"
	raven "github.com/getsentry/raven-go"
)

var cmd string
var logfilename string
var config *Config

func init() {
	flag.StringVar(&cmd, "cmd", "", "command")
	flag.StringVar(&logfilename, "log", "", "log filename")

	config = LoadConfig()
	raven.SetDSN(config.SentryDSN)
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

	switch cmd {
	case "devel":
		mainDevel(config)
	case "":
		mainDefault(config)
	default:
		log.Fatalf("unknown command")
	}
}

func mainDevel(config *Config) {
}

func mainDefault(config *Config) {
	go mainServer(config)
	go mainDaemon(config)
	go mainDirectMessageStreaming(config)
	mainStreaming(config)
}

func mainDaemon(config *Config) {
	rs := config.NewDaemonRules()
	for _, r := range rs {
		go r.Execute()
	}
}

func mainServer(config *Config) {
	type Status struct {
		Ok  bool      `json:"ok"`
		Now time.Time `json:"now"`
	}
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status := Status{
			Ok:  true,
			Now: time.Now(),
		}
		data, _ := json.Marshal(status)
		var out bytes.Buffer
		json.Indent(&out, data, "", "  ")
		w.Write(out.Bytes())
	})

	// tweet id로 tweet 까보는 기능이 있으면 디버깅할떄 편하겠지?
	http.HandleFunc("/tweet/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server started: http://0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mainStreaming(config *Config) {
	handlers := config.NewTweetRules()

	api := config.NewDataSourceAuthConfig().CreateApi()
	v := url.Values{}
	twitterStream := api.UserStream(v)
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
			// pass
		case anaconda.EventTweet:
			evt := tweet.Event.Event
			for _, h := range handlers {
				h.OnEvent(evt, &tweet)
			}
		case anaconda.DirectMessage:
			// pass
		default:
			if x != nil {
				log.Printf("unknown type(%T) : %v \n", x, x)
			}
		}
	}
}

func mainDirectMessageStreaming(config *Config) {
	rs := config.NewMessageRules()
	api := config.CreateTwitterSenderApi()
	v := url.Values{}
	twitterStream := api.UserStream(v)
	for {
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.DirectMessage:
			for _, h := range rs {
				go h.OnDirectMessage(&tweet)
			}
		}
	}
}
