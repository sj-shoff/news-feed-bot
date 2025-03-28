package config

import (
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken     string        `hcl:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `hcl:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
	DatabaseDSN          string        `hcl:"database_dsn" env:"DATABASE_DSN" default:"postgres://postgres:postgres@localhost:5432/news_feed_bot?sslmode=disable"`
	FetchInterval        time.Duration `hcl:"fetch_interval" env:"FETCH_INTERVAL" default:"10m"`
	NotificationInterval time.Duration `hcl:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1m"`
	FilterKeywords       []string      `hcl:"filter_keywords" env:"FILTER_KEYWORDS"`
	OpenAIKey            string        `hcl:"openai_key" env:"OPENAI_KEY"`
	OpenAIPrompt         string        `hcl:"openai_prompt" env:"OPENAI_PROMPT"`
	OpenAIModel          string        `hcl:"openai_model" env:"OPENAI_MODEL"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		v := viper.New()
		v.SetEnvPrefix("NFB")
		v.AutomaticEnv()

		v.SetConfigType("hcl")
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/news-feed-bot")

		if err := v.ReadInConfig(); err == nil {
			slog.Info("Using config:", "config_file", v.ConfigFileUsed())
		}

		v.SetConfigName("config.local")
		if err := v.MergeInConfig(); err == nil {
			slog.Info("Using local config:", "config_file", v.ConfigFileUsed())
		}

		if err := v.Unmarshal(&cfg); err != nil {
			slog.Error("config unmarshal error: %v", sl.Err(err))
		}
	})

	return cfg
}
