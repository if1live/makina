package main

import (
	"encoding/json"
	"io/ioutil"

	"path"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/senders"
	"github.com/if1live/makina/storages"
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

	PushbulletAccessToken string `json:"pushbullet_access_token"`

	UseDummy         bool `json:"use_dummy"`
	UseMediaArchiver bool `json:"use_media_archiver"`
	UseHaru          bool `json:"use_haru"`

	StorageName string `json:"storage_name"`

	HaruFilePath string `json:"haru_filepath"`
	HaruHostName string `json:"haru_hostname"`

	SenderCategory string `json:"sender_category"`
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

func (config *Config) NewStorageAccessor(rootpath string) storages.Accessor {
	switch config.StorageName {
	case "local":
		return storages.NewLocal()
	case "dropbox":
		return storages.NewDropbox(rootpath, config.DropboxAccessToken)
	default:
		return nil
	}
}

func (config *Config) MakeSender(category string) *senders.Sender {
	switch category {
	case "pushbullet":
		s := senders.NewPushbullet(config.PushbulletAccessToken)
		return senders.New(s)
	case "dm":
		api := config.NewDataSourceAuthConfig().CreateApi()
		s := senders.NewDirectMessage(api, config.DataSourceScreenName)
		return senders.New(s)
	default:
		s := senders.NewFake()
		return senders.New(s)
	}
}

func (config *TwitterAuthConfig) CreateApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	return api
}
