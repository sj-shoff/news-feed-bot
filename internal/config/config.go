package config

import (
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
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
	TelegramBotToken     string        `yaml:"telegram_bot_token" env:"BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `yaml:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
	FetchInterval        time.Duration `yaml:"fetch_interval" env:"FETCH_INTERVAL" default:"10m"`
	NotificationInterval time.Duration `yaml:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1m"`
	FilterKeywords       []string      `yaml:"filter_keywords" env:"FILTER_KEYWORDS"`
	OpenAIKey            string        `yaml:"openai_key" env:"OPENAI_KEY"`
	OpenAIPrompt         string        `yaml:"openai_prompt" env:"OPENAI_PROMPT"`
	OpenAIModel          string        `yaml:"openai_model" env:"OPENAI_MODEL" default:"gpt-3.5-turbo"`
	Postgres             DBConfig      `yaml:"postgres"`
	Env                  string        `yaml:"env" env-default:"local"`
}

func Get() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		slog.Error("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist: %s", configPath, sl.Err(err))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Cannot read config: %s", sl.Err(err))
	}

	return &cfg
}
