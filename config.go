package main

import (
	"encoding/json"
	"io/ioutil"

	"path"

	"github.com/ChimeraCoder/anaconda"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

type TwitterAuthConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type Config struct {
	DataSourceConsumerKey       string `json:"data_source_consumer_key"`
	DataSourceConsumerSecret    string `json:"data_source_consumer_secret"`
	DataSourceAccessToken       string `json:"data_source_access_token"`
	DataSourceAccessTokenSecret string `json:"data_source_access_token_secret"`
	DataSourceScreenName        string `json:"data_source_screen_name"`

	DropboxAppKey      string `json:"dropbox_app_key"`
	DropboxAppSecret   string `json:"dropbox_app_secret"`
	DropboxAccessToken string `json:"dropbox_access_token"`
}

func LoadConfig() *Config {
	filename := "config.json"
	filepath := path.Join(GetExecutablePath(), filename)

	var config Config
	data, errFile := ioutil.ReadFile(filepath)
	check(errFile)
	errJson := json.Unmarshal(data, &config)
	check(errJson)
	return &config
}

func (config *Config) NewDataSourceAuthConfig() *TwitterAuthConfig {
	return &TwitterAuthConfig{
		config.DataSourceConsumerKey,
		config.DataSourceConsumerSecret,
		config.DataSourceAccessToken,
		config.DataSourceAccessTokenSecret,
	}
}

func (config *TwitterAuthConfig) CreateApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	return api
}

func (config *Config) CreateDropboxClient() *dropy.Client {
	token := config.DropboxAccessToken
	client := dropy.New(dropbox.New(dropbox.NewConfig(token)))
	return client
}
