package notifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider  // поставщик статей
	summarizer       Summarizer       // cаммари
	bot              *tgbotapi.BotAPI // инстанс бота
	sendInterval     time.Duration    // интервал рассылки
	lookupTimeWindow time.Duration    // интервал в котором проверяются новые статьи
	channelID        int64            // id канала куда будет все поститься
}

func New(articles ArticleProvider, summarizer Summarizer, bot *tgbotapi.BotAPI, sendInterval, lookupTimeWindow time.Duration, channelID int64) *Notifier {
	return &Notifier{
		articles:         articles,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
	}
}

// !!! НУЖНО СДЕЛАТЬ РЕАЛИЗАЦИЮ КАК ТРАНЗАКЦИЮ
func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneAtricles, err := n.articles.AllNotPosted(ctx, time.Now().UTC().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return err
	}

	if len(topOneAtricles) == 0 {
		return nil
	}

	article := topOneAtricles[0]

	sumamry, err := n.extractSummary(ctx, article)
	//summary, err := n.summarizer.Summarize(ctx, article.Summary)

	if err := n.sendArticle(article, sumamry); err != nil {
		return err
	}

	return n.articles.MarkAsPosted(ctx, article.ID)
}

// добавили очистку html тегов, (go-readability), выводим отформированный текст
func (n *Notifier) extractSummary(ctx context.Context, article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.Summarize(ctx, cleanText(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil // пустые две строки тк перед самари будет заголовок
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`) // эта регулярка соответствует всем последовательностям  от и более 3х пустых строк

func cleanText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(summary),
		markup.EscapeForMarkdown(article.Link),
	))
	msg.ParseMode = "MarkdownV2"

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
