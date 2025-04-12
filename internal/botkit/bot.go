package botkit

import (
	"context"
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI          // Клиент Telegram API
	cmdViews map[string]CommandHandler // Реестр обработчиков команд в мапе
}

// CommandHandler определяет сигнатуру функции-обработчика команд
type CommandHandler func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

// конструктор для бота
func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api:      api,
		cmdViews: make(map[string]CommandHandler),
	}
}

// RegisterCommand регистрирует обработчик для конкретной команды
func (b *Bot) RegisterCommand(command string, handler CommandHandler) {
	b.cmdViews[command] = handler
}

// Run запускает главный цикл обработки сообщений
func (b *Bot) Run(ctx context.Context) error {
	// Настраиваем канал обновлений
	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = 60 // Long polling timeout

	updates := b.api.GetUpdatesChan(updatesConfig)

	// Главный цикл обработки событий
	for {
		select {
		case update := <-updates:
			go b.processUpdate(update)

		case <-ctx.Done():
			slog.Info("Bot shutting down...")
			return ctx.Err()
		}
	}
}

// processUpdate обрабатывает отдельное обновление
func (b *Bot) processUpdate(update tgbotapi.Update) {
	// Восстановление после паники
	defer recoverFromPanic()

	// Валидация входящего обновления
	if !isValidUpdate(update) {
		return
	}

	// Извлечение команды из сообщения
	command := extractCommand(update)
	if command == "" {
		return
	}

	// Поиск обработчика команды
	handler, exists := b.cmdViews[command]
	if !exists {
		return
	}

	// Создание контекста с таймаутом для обработки
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Выполнение обработчика
	if err := handler(ctx, b.api, update); err != nil {
		handleProcessingError(b.api, update, err)
	}
}

// recoverFromPanic обрабатывает непредвиденные паники горутины
func recoverFromPanic() {
	if r := recover(); r != nil {
		slog.Info(
			"PANIC RECOVERED: %v\nStack Trace:\n%s",
			r,
			string(debug.Stack()),
		)
	}
}

// isValidUpdate проверяет валидность обновления
func isValidUpdate(update tgbotapi.Update) bool {
	return (update.Message != nil && update.Message.IsCommand()) ||
		update.CallbackQuery != nil
}

// extractCommand извлекает команду из сообщения
func extractCommand(update tgbotapi.Update) string {
	if update.Message == nil || !update.Message.IsCommand() {
		return ""
	}
	return update.Message.Command()
}

// handleProcessingError обрабатывает ошибки выполнения
func handleProcessingError(bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {

	slog.Info("Processing error: %v", sl.Err(err))

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		return
	}

	// Отправка сообщения об ошибке пользователю
	msg := tgbotapi.NewMessage(chatID, "Произошла внутренняя ошибка")
	if _, sendErr := bot.Send(msg); sendErr != nil {
		slog.Info("Failed to send error message: %v", sl.Err(sendErr))
	}
}
