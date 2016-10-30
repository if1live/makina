package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type FetchResponse struct {
	Url        string
	Data       []byte
	Date       time.Time
	StatusCode int
}

func (r *FetchResponse) IsSuccess() bool {
	return 200 <= r.StatusCode && r.StatusCode < 300
}

func NewErrorResponse(rawurl string, code int) *FetchResponse {
	return &FetchResponse{
		Url:        rawurl,
		Data:       []byte{},
		Date:       time.Now(),
		StatusCode: code,
	}
}

type HttpFetcher struct {
}

func (f *HttpFetcher) handleError(err error, rawurl string) *FetchResponse {
	// http://stackoverflow.com/questions/22761562/portable-way-to-detect-different-kinds-of-network-error-in-golang
	switch err.(type) {
	case *url.Error:
		switch err.(*url.Error).Err.(type) {
		case *net.OpError:
			return NewErrorResponse(rawurl, 602)

		default:
			// err.Err가 *http.httpError 인 경우가 있는데
			// private라서 쓸수없는 클래스다. 그래서 이름으로 분기
			if reflect.TypeOf(err.(*url.Error).Err).String() == "*http.httpError" {
				return NewErrorResponse(rawurl, 601)
			}
		}
	}
	fmt.Printf("unknown error type: %Q\n", err)
	return NewErrorResponse(rawurl, 600)
}

func (f *HttpFetcher) Fetch(rawurl string) *FetchResponse {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	// http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
	// 오래 걸리는 작업은 워커로 넘길거니까 타임아웃을 굳이 설정할 필요는 없을듯

	// http://stackoverflow.com/questions/12122159/golang-how-to-do-a-https-request-with-bad-certificate
	client := &http.Client{}
	switch u.Scheme {
	case "http":
		client = &http.Client{}
	case "https":
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
		}
	default:
		// 없으면 http로 취급
		panic("unknown scheme")
	}
	resp := f.fetch(client, rawurl)

	if resp.IsSuccess() {
		log.Printf("HttpFetcher: %s -> success\n", rawurl)
	} else {
		log.Printf("HttpFetcher: %s -> fail\n", rawurl)
	}
	return resp
}

func (f *HttpFetcher) fetch(client *http.Client, rawurl string) *FetchResponse {
	response, err := client.Get(rawurl)
	if err != nil {
		return f.handleError(err, rawurl)
	}
	defer response.Body.Close()

	createSuccessFunc := func() *FetchResponse {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		return &FetchResponse{
			Url:        rawurl,
			Data:       buf.Bytes(),
			Date:       time.Now(),
			StatusCode: response.StatusCode,
		}
	}

	switch response.StatusCode {
	case 404:
		return NewErrorResponse(rawurl, 404)
	case 200:
		return createSuccessFunc()
	default:
		return createSuccessFunc()
	}
}
