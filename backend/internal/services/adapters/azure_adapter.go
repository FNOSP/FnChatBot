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

type AzureAdapter struct {
	*OpenAIAdapter
}

func NewAzureAdapter() *AzureAdapter {
	return &AzureAdapter{
		OpenAIAdapter: NewOpenAIAdapter(),
	}
}

func (a *AzureAdapter) ProviderType() models.ProviderType {
	return models.ProviderTypeAzureOpenAI
}

type azureModelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

func (a *AzureAdapter) FetchModels(ctx context.Context, baseURL, apiKey string) ([]services.ModelInfo, error) {
	url := baseURL + "/deployments?api-version=2024-02-15-preview"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", apiKey)
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
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var modelsResp azureModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if modelsResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", modelsResp.Error.Message)
	}

	result := make([]services.ModelInfo, 0, len(modelsResp.Data))
	for _, m := range modelsResp.Data {
		result = append(result, services.ModelInfo{
			ID:      m.ID,
			Name:    m.ID,
			OwnedBy: m.OwnedBy,
		})
	}

	return result, nil
}

type azureChatRequest struct {
	DeploymentID string                 `json:"-"`
	Messages     []services.ChatMessage `json:"messages"`
	Stream       bool                   `json:"stream"`
	Temperature  float32                `json:"temperature,omitempty"`
	MaxTokens    int                    `json:"max_tokens,omitempty"`
	Tools        []services.Tool        `json:"tools,omitempty"`
}

func (a *AzureAdapter) BuildChatRequest(ctx context.Context, model string, messages []services.ChatMessage, options services.ChatOptions) (*http.Request, error) {
	reqBody := azureChatRequest{
		DeploymentID: model,
		Messages:     messages,
		Stream:       options.Stream,
		Temperature:  options.Temperature,
		MaxTokens:    options.MaxTokens,
		Tools:        options.Tools,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/deployments/%s/chat/completions?api-version=2024-02-15-preview", options.BaseURL, model)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", options.APIKey)

	return req, nil
}

type azureStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int              `json:"index"`
		Delta        azureStreamDelta `json:"delta"`
		FinishReason string           `json:"finish_reason"`
	} `json:"choices"`
}

type azureStreamDelta struct {
	Role      string              `json:"role,omitempty"`
	Content   string              `json:"content,omitempty"`
	ToolCalls []services.ToolCall `json:"tool_calls,omitempty"`
}

func (a *AzureAdapter) ParseStreamResponse(ctx context.Context, reader *bufio.Reader) (services.StreamChunk, error) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return services.StreamChunk{Done: true}, nil
			}
			return services.StreamChunk{}, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return services.StreamChunk{Done: true}, nil
		}

		var streamResp azureStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) == 0 {
			continue
		}

		choice := streamResp.Choices[0]
		chunk := services.StreamChunk{
			Content:      choice.Delta.Content,
			FinishReason: choice.FinishReason,
			Done:         false,
		}

		if len(choice.Delta.ToolCalls) > 0 {
			chunk.ToolCalls = choice.Delta.ToolCalls
		}

		return chunk, nil
	}
}

func (a *AzureAdapter) GetBaseURL(baseURL string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	return baseURL
}
