package adapters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fnchatbot/internal/models"
	"fnchatbot/internal/services"
)

type OllamaAdapter struct {
	baseURL string
	apiKey  string
}

func NewOllamaAdapter() *OllamaAdapter {
	return &OllamaAdapter{}
}

func (a *OllamaAdapter) ProviderType() models.ProviderType {
	return models.ProviderTypeOllama
}

type ollamaTagsResponse struct {
	Models []struct {
		Name     string `json:"name"`
		Modified string `json:"modified_at"`
		Size     int64  `json:"size"`
		Digest   string `json:"digest"`
	} `json:"models"`
}

func (a *OllamaAdapter) FetchModels(ctx context.Context, baseURL, apiKey string) ([]services.ModelInfo, error) {
	url := a.GetBaseURL(baseURL) + "/api/tags"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var tagsResp ollamaTagsResponse
	if err := json.Unmarshal(body, &tagsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]services.ModelInfo, 0, len(tagsResp.Models))
	for _, m := range tagsResp.Models {
		result = append(result, services.ModelInfo{
			ID:      m.Name,
			Object:  "model",
			OwnedBy: "ollama",
		})
	}

	return result, nil
}

type ollamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []ollamaChatMessage    `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Tools    []services.Tool        `json:"tools,omitempty"`
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (a *OllamaAdapter) BuildChatRequest(ctx context.Context, model string, messages []services.ChatMessage, options services.ChatOptions) (*http.Request, error) {
	url := a.GetBaseURL(a.baseURL) + "/api/chat"

	ollamaMessages := make([]ollamaChatMessage, 0, len(messages))
	for _, m := range messages {
		if m.Role == "tool" {
			ollamaMessages = append(ollamaMessages, ollamaChatMessage{
				Role:    "user",
				Content: fmt.Sprintf("Tool result from %s: %s", m.Name, m.Content),
			})
			continue
		}
		ollamaMessages = append(ollamaMessages, ollamaChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	ollamaOptions := make(map[string]interface{})
	if options.Temperature > 0 {
		ollamaOptions["temperature"] = options.Temperature
	}
	if options.MaxTokens > 0 {
		ollamaOptions["num_predict"] = options.MaxTokens
	}

	reqBody := ollamaChatRequest{
		Model:    model,
		Messages: ollamaMessages,
		Stream:   options.Stream,
	}
	if len(ollamaOptions) > 0 {
		reqBody.Options = ollamaOptions
	}
	if len(options.Tools) > 0 {
		reqBody.Tools = options.Tools
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

type ollamaChatResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done          bool   `json:"done"`
	TotalDuration int64  `json:"total_duration,omitempty"`
	LoadDuration  int64  `json:"load_duration,omitempty"`
	PromptCount   int    `json:"prompt_eval_count,omitempty"`
	EvalCount     int    `json:"eval_count,omitempty"`
	DoneReason    string `json:"done_reason,omitempty"`
}

func (a *OllamaAdapter) ParseStreamResponse(ctx context.Context, reader *bufio.Reader) (services.StreamChunk, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return services.StreamChunk{Done: true}, nil
		}
		return services.StreamChunk{}, fmt.Errorf("error reading stream: %w", err)
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return services.StreamChunk{}, nil
	}

	var resp ollamaChatResponse
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return services.StreamChunk{}, fmt.Errorf("failed to parse response: %w", err)
	}

	chunk := services.StreamChunk{
		Content: resp.Message.Content,
		Done:    resp.Done,
	}

	if resp.Done && resp.DoneReason != "" {
		chunk.FinishReason = resp.DoneReason
	}

	return chunk, nil
}

func (a *OllamaAdapter) GetBaseURL(baseURL string) string {
	return strings.TrimSuffix(baseURL, "/")
}

func (a *OllamaAdapter) SetCredentials(baseURL, apiKey string) {
	a.baseURL = baseURL
	a.apiKey = apiKey
}
