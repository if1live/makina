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

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/hitomi"
	"github.com/if1live/makina/media_archiver"
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
	OnEvent(ev string, event *anaconda.EventTweet)
	OnDirectMessage(dm *anaconda.DirectMessage)
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
	handlers := []StreamingHandler{}
	if config.UseMediaArchiver {
		const savePath = "/archive-temp"
		cfg := media_archiver.Config{
			Accessor: config.NewStorageAccessor(savePath),
			MyName:   config.DataSourceScreenName,
		}
		handlers = append(handlers, media_archiver.NewMediaArchiver(cfg))
	}

	if config.UseHaru {
		const savePath = "/hitomi-temp"
		cfg := hitomi.Config{
			MyName:       config.DataSourceScreenName,
			Accessor:     config.NewStorageAccessor(savePath),
			HaruHostName: config.HaruHostName,
		}
		handlers = append(handlers, hitomi.NewListener(cfg))
	}

	api := config.NewDataSourceAuthConfig().CreateApi()
	v := url.Values{}
	twitterStream := api.UserStream(v)
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
			// pass
		case anaconda.EventTweet:
			evt := tweet.Event.Event
			for _, h := range handlers {
				go h.OnEvent(evt, &tweet)
			}
		case anaconda.DirectMessage:
			for _, h := range handlers {
				go h.OnDirectMessage(&tweet)
			}
		default:
			if x != nil {
				log.Printf("unknown type(%T) : %v \n", x, x)
			}
		}
	}
}
