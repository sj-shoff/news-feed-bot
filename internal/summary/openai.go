package summary

import (
	"context"
	"errors"
	"log/slog"
	"news-feed-bot/internal/logger/sl"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type SummaryConfig struct {
	Prompt      string
	Model       string
	MaxTokens   int
	Temperature float32
	TopP        float32
}

type OpenAISummarizer struct {
	client  *openai.Client
	config  SummaryConfig
	logger  *slog.Logger
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummarizer(client *openai.Client, config SummaryConfig, log *slog.Logger) Summarizer {
	enabled := client != nil
	log = log.With(
		slog.String("component", "OpenAISummarizer"),
		slog.Bool("enabled", enabled),
		slog.String("model", config.Model),
	)

	if enabled {
		log.Info("initializing OpenAI summarizer")
	}

	return &OpenAISummarizer{
		client:  client,
		config:  config,
		logger:  log,
		enabled: enabled,
	}
}

func (s *OpenAISummarizer) Summarize(ctx context.Context, text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log := s.logger
	if !s.enabled {
		log.Error("summarizer disabled")
		return "", errors.New("OpenAI summarizer is disabled")
	}

	log.Debug("starting summarization")

	request := s.createRequest(text)
	log.Debug("sending request to OpenAI")

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Error("API request failed", sl.Err(err))
		return "", errors.New("API request failed: " + err.Error())
	}

	summary, err := s.extractSummary(resp)
	if err != nil {
		log.Error("failed to extract summary", sl.Err(err))
		return "", err
	}

	log.Debug("summarization completed",
		slog.Int("result_length", len(summary)))

	return summary, nil
}

func (s *OpenAISummarizer) createRequest(text string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:       s.config.Model,
		Messages:    s.createMessages(text),
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		TopP:        s.config.TopP,
	}
}

func (s *OpenAISummarizer) createMessages(text string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: s.config.Prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		},
	}
}

func (s *OpenAISummarizer) extractSummary(resp openai.ChatCompletionResponse) (string, error) {
	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from OpenAI")
	}

	rawSummary := strings.TrimSpace(resp.Choices[0].Message.Content)
	result := ensureSentenceEnding(rawSummary)
	return result, nil
}

func ensureSentenceEnding(s string) string {
	if strings.HasSuffix(s, ".") {
		return s
	}

	lastDot := strings.LastIndex(s, ".")
	if lastDot == -1 {
		return s + "."
	}

	return s[:lastDot+1]
}
