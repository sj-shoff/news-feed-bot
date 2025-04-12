package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"news-feed-bot/internal/bot"
	"news-feed-bot/internal/bot/middleware"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/config"
	"news-feed-bot/internal/fetcher"
	"news-feed-bot/internal/logger/sl"
	"news-feed-bot/internal/logger/slogpretty"
	"news-feed-bot/internal/notifier"
	"news-feed-bot/internal/storage"
	"news-feed-bot/internal/storage/postgres"
	"news-feed-bot/internal/summary"

	"os/signal"
	"syscall"

	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Get()

	log := setupLogger(cfg.Env)
	log.Info(
		"starting news-feed-bot",
		slog.String("env", cfg.Env),
		slog.String("version", "1.0"),
	)

	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		slog.Info("Failed to create bot", sl.Err(err))
		return
	}

	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Error("Failed to init db", sl.Err(err))
		os.Exit(1)
	}

	dbStorage := db.DB
	var (
		articleStorage = storage.NewArticleStorage(dbStorage)
		sourceStorage  = storage.NewSourceStorage(dbStorage)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		summarizer = summary.NewOpenAISummarizer(
			config.Get().OpenAIKey,
			config.Get().OpenAIModel,
			config.Get().OpenAIPrompt,
		)
		notifier = notifier.New(
			articleStorage,
			summarizer,
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCommand(
		"addsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCommand(
		"setpriority",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdSetPriority(sourceStorage),
		),
	)
	newsBot.RegisterCommand(
		"getsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdGetSource(sourceStorage),
		),
	)
	newsBot.RegisterCommand(
		"listsources",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdListSource(sourceStorage),
		),
	)
	newsBot.RegisterCommand(
		"deletesource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("Failed to run fetcher: %v", sl.Err(err))
				return
			}

			log.Info("Fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("Failed to run notifier: %v", sl.Err(err))
				return
			}

			log.Info("Notifier stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := http.ListenAndServe("0.0.0.0:8088", mux); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("Failed to run http server: %v", sl.Err(err))
				return
			}

			log.Info("Http server stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		log.Error("Failed to run botkit: %v", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
