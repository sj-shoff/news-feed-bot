# News Feed Bot

Telegram-бот для управления RSS-источниками новостей.

## Возможности

- Добавление новых RSS-источников
- Удаление источников
- Просмотр информации об источнике
- Просмотр списка всех источников
- Установка приоритета источников

## Команды

### Для администраторов
1. `/addsource {json}` - добавить новый источник
   Пример: `/addsource {"name":"Example","url":"https://example.com/rss","priority":1}`
   
2. `/deletesource {id}` - удалить источник по ID
   Пример: `/deletesource 1`

3. `/getsource {id}` - получить информацию об источнике
   Пример: `/getsource 1`

4. `/listsource` - получить список всех источников

5. `/setpriority {json}` - установить приоритет источника
   Пример: `/setpriority {"source_id":1,"priority":5}`

## Установка

1. Клонировать репозиторий
2. Установить зависимости: `go mod download`
3. Настроить конфигурацию
4. Запустить: `go run cmd/bot/main.go`

## Требования

- Go
- Telegram Bot API токен
- Доступ к базе данных (реализация интерфейсов Storage)

## Конфигурация

Настройки бота задаются через переменные окружения:
- `TELEGRAM_TOKEN` - токен Telegram бота
- `CHANNEL_ID` - ID канала для проверки прав администратора
- `DB_DSN` - строка подключения к базе данных

## Лицензия

MIT
### Для запуска приложения:

```
make build && make run
```
