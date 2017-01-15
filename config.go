package main

import (
	"io/ioutil"
	"path"

	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/yaml.v2"
)

type TwitterAuthConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type Config struct {
	DataSourceConsumerKey       string `yaml:"data_source_consumer_key"`
	DataSourceConsumerSecret    string `yaml:"data_source_consumer_secret"`
	DataSourceAccessToken       string `yaml:"data_source_access_token"`
	DataSourceAccessTokenSecret string `yaml:"data_source_access_token_secret"`
	DataSourceScreenName        string `yaml:"data_source_screen_name"`

	DropboxAccessToken string `yaml:"dropbox_access_token"`

	TwitterSenderAccessToken       string `yaml:"twitter_sender_access_token"`
	TwitterSenderAccessTokenSecret string `yaml:"twitter_sender_access_token_secret"`

	MediaArchiverPredefineUsers []string `yaml:"media_archiver_predefined_users"`

	SentryDSN string `yaml:"sentry_dsn"`
}

func LoadConfig() *Config {
	filename := "config.yaml"
	filepath := path.Join(GetExecutablePath(), filename)

	var config Config
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &Config{}
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return &Config{}
	}

	return &config
}

func (c *Config) NewDataSourceAuthConfig() *TwitterAuthConfig {
	return &TwitterAuthConfig{
		c.DataSourceConsumerKey,
		c.DataSourceConsumerSecret,
		c.DataSourceAccessToken,
		c.DataSourceAccessTokenSecret,
	}
}

func (c *Config) NewStorage(rootpath string) *Storage {
	return NewStorage(rootpath, c.DropboxAccessToken)
}

func (c *Config) MakeSender() Sender {
	api := config.CreateTwitterSenderApi()
	return NewSender(api, config.DataSourceScreenName)
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

func (c *Config) NewTweetRules() []TweetRule {
	rs := []TweetRule{}
	{
		const savePath = "/archive-temp"
		a := c.NewStorage(savePath)
		r := NewMediaArchiver(a, c.DataSourceScreenName, c.MediaArchiverPredefineUsers)
		rs = append(rs, r)
	}
	{
		const savePath = "/hitomi-temp"
		a := c.NewStorage(savePath)
		r := NewHitomiWatcher(c.DataSourceScreenName, a)
		rs = append(rs, r)
	}
	return rs
}
func (c *Config) NewMessageRules() []MessageRule {
	sender := c.MakeSender()
	rs := []MessageRule{}
	{
		r := NewDirectMessageWatcher(c.DataSourceScreenName, sender)
		rs = append(rs, r)
	}
	return rs
}

func (c *Config) NewDaemonRules() []DaemonRule {
	rs := []DaemonRule{}
	return rs
}
