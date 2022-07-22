package models

import (
	"os"

	"github.com/mrjones/oauth"
)

type Config struct {
	URL           string `env:"URL,required"`
	PORT          string `env:"PORT,required"`
	TelegramToken string `env:"TELEGRAM_TOKEN,required"`
	MysqlUri      string `env:"MYSQL_URI,required"`
	Redis         string `env:"REDISCLOUD_URL,required"`

	TwitterKey    string `env:"TWITTER_KEY,required"`
	TwitterSecret string `env:"TWITTER_SECRET,required"`
	TwitterToken  string `env:"TWITTER_TOKEN,required"`

	Debug bool `env:"DEBUG,required"`
}

var (
	TwitterConsumer = oauth.NewConsumer(
		os.Getenv("TWITTER_KEY"),
		os.Getenv("TWITTER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
	Tokens = map[string]*oauth.RequestToken{}
)
