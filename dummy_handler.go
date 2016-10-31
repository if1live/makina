package main

import "github.com/ChimeraCoder/anaconda"

type DummyHandler struct {
	config *Config
}

func NewDummyHandler(config *Config) *DummyHandler {
	return &DummyHandler{config}
}

func (h *DummyHandler) OnTweet(tweet *anaconda.Tweet) {
}

func (h *DummyHandler) OnFavorite(tweet *anaconda.EventTweet) {
}

func (h *DummyHandler) OnUnfavorite(tweet *anaconda.EventTweet) {
}
