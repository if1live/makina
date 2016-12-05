package rules

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"log"

	"fmt"

	"github.com/gregjones/httpcache/diskcache"
	"github.com/if1live/httpforcecache"
	"github.com/if1live/makina/senders"
	"github.com/xrash/smetrics"
)

type PageWatcher struct {
	tasks        []*PageWatchTask
	StatusSender *senders.Sender
}

type PageWatchTask struct {
	url             string
	duration        time.Duration
	allowedDistance float64
	sender          *senders.Sender
	cache           *diskcache.Cache
}

func NewPageWatchTask(url string, duration time.Duration, allowedDistance float64, sender *senders.Sender) *PageWatchTask {
	const cachedir = "./_cache_pagewatcher"
	c := diskcache.New(cachedir)
	return &PageWatchTask{
		url:             url,
		duration:        duration,
		allowedDistance: allowedDistance,
		sender:          sender,
		cache:           c,
	}
}
func (t *PageWatchTask) createClient() *http.Client {
	tp := httpforcecache.NewTransport(t.cache)
	client := &http.Client{Transport: tp}
	return client
}

func (t *PageWatchTask) Watch() {
	url := t.url

	// 쇼핑몰을 크롤링하는 경우
	// <input type="hidden" name="transactionid" value="bf9ba3fc807bcf61528eb4f27f55dcab57e07f6b" />
	// 같은 요소떄문에 100% 일치를 얻기 어렵다
	// 그래서 문자열의 유사도로 비교

	for {
		go func(url string) {
			log.Printf("Watch URL [%s] start...", url)

			prev := t.getCachedResponse(url)
			curr := t.getRealtimeResponse(url)

			dist := smetrics.Jaro(prev, curr)
			if len(curr) == 0 {
				log.Printf("Watch URL [%s] empty page found", url)
				t.sender.Send("Empty Page?", url)

			} else if prev == curr {
				log.Printf("Watch URL [%s] page not changed. wait...", url)

			} else if dist > t.allowedDistance {
				log.Printf("Watch URL [%s] page are similar, %f. wait...", url, dist)

			} else {
				log.Printf("Watch URL [%s] page changed? %f", url, dist)

				// 알림 보내기
				title := fmt.Sprintf("Page chaned (%f)", dist)
				t.sender.Send(title, url)

				// 캐시 다시 쓰기
				// 간단한 구현으로 캐시 날리고 다시 요청하면 되겠지?
				t.cache.Delete(url)
				t.getCachedResponse(url)
			}
		}(url)

		time.Sleep(t.duration)
	}
}

func (t *PageWatchTask) getCachedResponse(url string) string {
	client := t.createClient()
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	text := getResponseText(resp)
	return text
}
func (t *PageWatchTask) getRealtimeResponse(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	text := getResponseText(resp)
	return text
}

func getResponseText(resp *http.Response) string {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, resp.Body)
	if err != nil {
		return ""
	}
	err = resp.Body.Close()
	if err != nil {
		return ""
	}
	text := buf.String()
	return text
}

func NewPageWatcher(sender *senders.Sender) DaemonRule {
	tasks := []*PageWatchTask{
		//NewPageWatchTask("http://127.0.0.1:3000/status", time.Minute*1, 1.0, sender),
		NewPageWatchTask("http://www.cuffs.co.jp/main/event/", time.Hour*6, 0.99, sender),
		NewPageWatchTask("https://cuffs.dchd-ecshop.net/", time.Hour*6, 0.95, sender),
		NewPageWatchTask("https://store.unity.com/kr/download?ref=personal", time.Hour*24, 1.00, sender),
	}
	return &PageWatcher{
		tasks: tasks,
	}
}

func (pw *PageWatcher) Execute() {
	for _, t := range pw.tasks {
		go t.Watch()
	}
}
