package main

import (
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type AuthConfig struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}

func (config *AuthConfig) CreateApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.consumerKey)
	anaconda.SetConsumerSecret(config.consumerSecret)
	api := anaconda.NewTwitterApi(config.accessToken, config.accessTokenSecret)

	// use automatic throttling
	//api.SetDelay(1 * time.Second)

	return api
}

func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
	}
}
