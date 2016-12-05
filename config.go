package main

import (
	"encoding/json"
	"io/ioutil"

	"path"

	"github.com/ChimeraCoder/anaconda"
	"github.com/if1live/makina/rules"
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

	TwitterSenderAccessToken       string `json:"twitter_sender_access_token"`
	TwitterSenderAccessTokenSecret string `json:"twitter_sender_access_token_secret"`

	MediaArchiverPredefineUsers []string `json:"media_archiver_predefined_users"`

	StorageName string `json:"storage_name"`

	SenderCategory string `json:"sender_category"`

	SentryDSN string `json:"sentry_dsn"`
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
	case "dm":
		api := config.CreateTwitterSenderApi()
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

func (c *Config) CreateTwitterSenderApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(c.DataSourceConsumerKey)
	anaconda.SetConsumerSecret(c.DataSourceConsumerSecret)
	api := anaconda.NewTwitterApi(c.TwitterSenderAccessToken, c.TwitterSenderAccessTokenSecret)
	return api
}

func (c *Config) NewTweetRules() []rules.TweetRule {
	rs := []rules.TweetRule{}
	{
		const savePath = "/archive-temp"
		a := c.NewStorageAccessor(savePath)
		r := rules.NewMediaArchiver(a, c.DataSourceScreenName, c.MediaArchiverPredefineUsers)
		rs = append(rs, r)
	}
	{
		const savePath = "/hitomi-temp"
		a := c.NewStorageAccessor(savePath)
		r := rules.NewHitomiWatcher(c.DataSourceScreenName, a)
		rs = append(rs, r)
	}
	return rs
}
func (c *Config) NewMessageRules() []rules.MessageRule {
	sender := c.MakeSender(c.SenderCategory)
	rs := []rules.MessageRule{}
	{
		const savePath = "/dm-temp"
		a := c.NewStorageAccessor(savePath)
		r := rules.NewDirectMessageWatcher(c.DataSourceScreenName, a, sender)
		rs = append(rs, r)
	}
	return rs
}

func (c *Config) NewDaemonRules() []rules.DaemonRule {
	sender := c.MakeSender(c.SenderCategory)
	rs := []rules.DaemonRule{}
	{
		r := rules.NewPageWatcher(sender)
		rs = append(rs, r)
	}
	return rs
}
