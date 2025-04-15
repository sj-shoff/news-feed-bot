package config

import (
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	TelegramBotToken     string        `env:"BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `env:"TELEGRAM_CHANNEL_ID" required:"true"`
	FetchInterval        time.Duration `yaml:"fetch_interval"`
	NotificationInterval time.Duration `yaml:"notification_interval"`
	FilterKeywords       []string      `yaml:"filter_keywords"`
	OpenAIKey            string        `env:"OPENAI_KEY" required:"true"`
	OpenAIPrompt         string        `yaml:"openai_prompt"`
	OpenAIModel          string        `yaml:"openai_model"`
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

	cfg.TelegramBotToken = os.Getenv("BOT_TOKEN")
	cfg.OpenAIKey = os.Getenv("OPENAI_KEY")
	cfg.TelegramChannelID, _ = strconv.ParseInt(os.Getenv("TELEGRAM_CHANNEL_ID"), 10, 64)
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")

	return &cfg
}
