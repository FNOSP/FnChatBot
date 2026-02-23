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

type GeminiAdapter struct{}

func NewGeminiAdapter() *GeminiAdapter {
	return &GeminiAdapter{}
}

func (a *GeminiAdapter) ProviderType() models.ProviderType {
	return models.ProviderTypeGemini
}

func (a *GeminiAdapter) FetchModels(ctx context.Context, baseURL, apiKey string) ([]services.ModelInfo, error) {
	return []services.ModelInfo{
		{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash"},
		{ID: "gemini-2.0-flash-lite", Name: "Gemini 2.0 Flash Lite"},
		{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro"},
		{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash"},
		{ID: "gemini-1.5-flash-8b", Name: "Gemini 1.5 Flash 8B"},
		{ID: "gemini-1.0-pro", Name: "Gemini 1.0 Pro"},
	}, nil
}

type geminiRequest struct {
	Contents          []geminiContent    `json:"contents"`
	SystemInstruction *geminiInstruction `json:"systemInstruction,omitempty"`
	GenerationConfig  *geminiGenConfig   `json:"generationConfig,omitempty"`
	Tools             []geminiTool       `json:"tools,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text         string                  `json:"text,omitempty"`
	FunctionCall *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResp *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

type geminiInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiGenConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDecl `json:"functionDeclarations"`
}

type geminiFunctionDecl struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type geminiFunctionResponse struct {
	Name     string `json:"name"`
	Response struct {
		Content string `json:"content"`
	} `json:"response"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text         string              `json:"text,omitempty"`
				FunctionCall *geminiFunctionCall `json:"functionCall,omitempty"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (a *GeminiAdapter) BuildChatRequest(ctx context.Context, model string, messages []services.ChatMessage, options services.ChatOptions) (*http.Request, error) {
	var systemPrompt string
	contents := make([]geminiContent, 0)

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemPrompt = msg.Content
		case "user":
			contents = append(contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: msg.Content}},
			})
		case "assistant":
			parts := make([]geminiPart, 0)
			if msg.Content != "" {
				parts = append(parts, geminiPart{Text: msg.Content})
			}
			for _, tc := range msg.ToolCalls {
				var args map[string]interface{}
				json.Unmarshal([]byte(tc.Function.Arguments), &args)
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						Name: tc.Function.Name,
						Args: args,
					},
				})
			}
			contents = append(contents, geminiContent{
				Role:  "model",
				Parts: parts,
			})
		case "tool":
			contents = append(contents, geminiContent{
				Role: "user",
				Parts: []geminiPart{{
					FunctionResp: &geminiFunctionResponse{
						Name: msg.Name,
						Response: struct {
							Content string `json:"content"`
						}{Content: msg.Content},
					},
				}},
			})
		}
	}

	reqBody := geminiRequest{
		Contents: contents,
	}

	if systemPrompt != "" {
		reqBody.SystemInstruction = &geminiInstruction{
			Parts: []geminiPart{{Text: systemPrompt}},
		}
	}

	if options.Temperature > 0 || options.MaxTokens > 0 {
		reqBody.GenerationConfig = &geminiGenConfig{
			Temperature:     options.Temperature,
			MaxOutputTokens: options.MaxTokens,
		}
	}

	if len(options.Tools) > 0 {
		declarations := make([]geminiFunctionDecl, len(options.Tools))
		for i, tool := range options.Tools {
			declarations[i] = geminiFunctionDecl{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			}
		}
		reqBody.Tools = []geminiTool{{FunctionDeclarations: declarations}}
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := a.GetBaseURL(options.BaseURL)
	method := "generateContent"
	if options.Stream {
		method = "streamGenerateContent"
	}
	url := fmt.Sprintf("%s/models/%s:%s?key=%s", baseURL, model, method, options.APIKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (a *GeminiAdapter) ParseStreamResponse(ctx context.Context, reader *bufio.Reader) (services.StreamChunk, error) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return services.StreamChunk{Done: true}, nil
			}
			return services.StreamChunk{}, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "[" || line == "]" || line == "," {
			continue
		}

		var resp geminiResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			continue
		}

		if resp.Error != nil {
			return services.StreamChunk{}, fmt.Errorf("gemini API error: %s", resp.Error.Message)
		}

		if len(resp.Candidates) == 0 {
			continue
		}

		candidate := resp.Candidates[0]
		chunk := services.StreamChunk{
			FinishReason: candidate.FinishReason,
		}

		if len(candidate.Content.Parts) > 0 {
			part := candidate.Content.Parts[0]
			if part.Text != "" {
				chunk.Content = part.Text
			}
			if part.FunctionCall != nil {
				args, _ := json.Marshal(part.FunctionCall.Args)
				chunk.ToolCalls = []services.ToolCall{{
					ID:   part.FunctionCall.Name,
					Type: "function",
					Function: services.ToolFunc{
						Name:      part.FunctionCall.Name,
						Arguments: string(args),
					},
				}}
			}
		}

		return chunk, nil
	}
}

func (a *GeminiAdapter) GetBaseURL(baseURL string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return baseURL
}
