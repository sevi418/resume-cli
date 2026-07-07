package ai

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sevi418/resume-cli/internal/model"
)

const (
	defaultAPIBase = "https://api.openai.com/v1"
	defaultModel   = "gpt-4o-mini"
)

type AIService interface {
	ExtractResume(ctx context.Context, text string) (*model.Resume, error)
	ScoreMatch(ctx context.Context, resumeText, jdText string) (*model.Score, error)
}

type RealAIService struct {
	client *openai.Client
	model  string
}

func NewRealAIServiceFromEnv() (*RealAIService, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set; set it or use --mock")
	}

	apiBase := os.Getenv("OPENAI_API_BASE")
	if apiBase == "" {
		apiBase = defaultAPIBase
	}

	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = defaultModel
	}

	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = apiBase

	slog.Debug("configured real AI service", "model", modelName, "base_url", apiBase)
	return &RealAIService{
		client: openai.NewClientWithConfig(cfg),
		model:  modelName,
	}, nil
}

func (s *RealAIService) chatJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	start := time.Now()
	slog.Debug(
		"AI request started",
		"model", s.model,
		"system_chars", len([]rune(systemPrompt)),
		"prompt_chars", len([]rune(userPrompt)),
	)

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		Temperature: 0.1,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("AI returned an empty response")
	}
	slog.Debug(
		"AI request completed",
		"elapsed", time.Since(start),
		"choices", len(resp.Choices),
		"response_chars", len([]rune(resp.Choices[0].Message.Content)),
		"prompt_tokens", resp.Usage.PromptTokens,
		"completion_tokens", resp.Usage.CompletionTokens,
		"total_tokens", resp.Usage.TotalTokens,
	)
	return resp.Choices[0].Message.Content, nil
}
