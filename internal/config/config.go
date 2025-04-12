package config

import (
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Env                  string        `yaml:"env" env-default:"local"`
	TelegramBotToken     string        `yaml:"telegram_bot_token"`
	TelegramChannelID    int64         `yaml:"telegram_channel_id"`
	FetchInterval        time.Duration `yaml:"fetch_interval"`
	NotificationInterval time.Duration `yaml:"notification_interval"`
	FilterKeywords       []string      `yaml:"filter_keywords"`
	OpenAIKey            string        `yaml:"openai_key"`
	OpenAIPrompt         string        `yaml:"openai_prompt"`
	OpenAIModel          string        `yaml:"openai_model"`
	Postgres             DBConfig      `yaml:"postgres"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		v := viper.New()
		v.AutomaticEnv()

		v.SetConfigType("yaml")
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
			slog.Error("Config unmarshal error: %v", sl.Err(err))
		}
	})

	return cfg
}
