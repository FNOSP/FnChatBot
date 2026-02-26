package llm

import (
	"context"
	"fmt"

	"fnchatbot/internal/models"
	"fnchatbot/internal/services/memory"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// StreamChat generates a streaming response with tool support
// It does NOT handle history automatically to allow caller (WebSocket) to manage the loop and UI updates.
// The caller is responsible for passing the full conversation history.
func (s *Service) StreamChat(ctx context.Context, provider models.Provider, modelName string, messages []llms.MessageContent, tools []llms.Tool, streamCallback func(ctx context.Context, chunk []byte) error) (*llms.ContentResponse, error) {
	llm, err := s.createLLM(ctx, provider, modelName)
	if err != nil {
		return nil, err
	}

	opts := []llms.CallOption{
		llms.WithStreamingFunc(streamCallback),
	}
	if len(tools) > 0 {
		opts = append(opts, llms.WithTools(tools))
	}

	return llm.GenerateContent(ctx, messages, opts...)
}

// GetHistory returns the chat history for a session
func (s *Service) GetHistory(ctx context.Context, sessionID uint) ([]llms.ChatMessage, error) {
	history := memory.NewSQLiteHistory(s.DB, sessionID)
	return history.Messages(ctx)
}

// SaveUserMessage saves a user message to history
func (s *Service) SaveUserMessage(ctx context.Context, sessionID uint, content string) error {
	history := memory.NewSQLiteHistory(s.DB, sessionID)
	return history.AddUserMessage(ctx, content)
}

// SaveAIMessage saves an AI message to history
func (s *Service) SaveAIMessage(ctx context.Context, sessionID uint, content string) error {
	history := memory.NewSQLiteHistory(s.DB, sessionID)
	return history.AddAIMessage(ctx, content)
}

func (s *Service) createLLM(ctx context.Context, provider models.Provider, modelName string) (llms.Model, error) {
	switch provider.Type {
	case models.ProviderTypeOpenAI, models.ProviderTypeOpenAIResponse:
		opts := []openai.Option{
			openai.WithToken(provider.APIKey),
			openai.WithModel(modelName),
		}
		if provider.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(provider.BaseURL))
		}
		return openai.New(opts...)
	case models.ProviderTypeAnthropic:
		return anthropic.New(anthropic.WithToken(provider.APIKey), anthropic.WithModel(modelName))
	case models.ProviderTypeOllama:
		opts := []ollama.Option{
			ollama.WithModel(modelName),
		}
		if provider.BaseURL != "" {
			opts = append(opts, ollama.WithServerURL(provider.BaseURL))
		}
		return ollama.New(opts...)
	case models.ProviderTypeGemini:
		return googleai.New(ctx, googleai.WithAPIKey(provider.APIKey))
	default:
		// Default to OpenAI compatible for unknown types if they have BaseURL, otherwise error
		if provider.BaseURL != "" {
			opts := []openai.Option{
				openai.WithToken(provider.APIKey),
				openai.WithModel(modelName),
				openai.WithBaseURL(provider.BaseURL),
			}
			return openai.New(opts...)
		}
		return nil, fmt.Errorf("unsupported provider type: %s", provider.Type)
	}
}
