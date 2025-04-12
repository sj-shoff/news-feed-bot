# News Feed Bot

Telegram-бот для автоматического сбора новостей из RSS-лент, генерации кратких содержаний через OpenAI и публикации в Telegram-канал.

## Основные функции

- **Управление источниками**: 
```
Добавление новых RSS-источников
Удаление источников
Просмотр информации об источнике
Просмотр списка всех источников
Установка приоритета источников
```
- **AI-суммаризация**: Генерация кратких содержаний статей с помощью OpenAI (GPT-3.5/GPT-4).
- **Автопубликация**: Планирование постов в Telegram-канал с поддержкой Markdown.

## Технологии

- **Backend**: Go, slog, <a href="https://github.com/spf13/viper">spf13/viper</a>
- **Базы данных**: PostgreSQL, <a href="https://github.com/jmoiron/sqlx">sqlx</a>
- **AI**: OpenAI API
- **Инфраструктура**: Docker, Docker-compose, Makefile
- **Интеграции**: Telegram Bot API, RSS 2.0

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

## Для запуска приложения:

```
make build && make run
```

## Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate
```