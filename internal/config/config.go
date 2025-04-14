package config

import (
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"os"
	"path/filepath"
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

		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			exePath, err := os.Executable()
			if err != nil {
				panic("failed to get executable path: " + err.Error())
			}
			configPath = filepath.Join(filepath.Dir(exePath), "config", "config.yaml")
		}

		slog.Info("Loading configuration", "path", configPath)

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			slog.Error("Config file not found", sl.Err(err))
			panic("config file not found at " + configPath)
		}

		v.SetConfigFile(configPath)

		if err := v.ReadInConfig(); err != nil {
			slog.Error("Failed to read config", sl.Err(err))
			panic(err)
		}

		if err := v.Unmarshal(&cfg); err != nil {
			slog.Error("Failed to unmarshal config", sl.Err(err))
			panic(err)
		}

	})

	return cfg
}
