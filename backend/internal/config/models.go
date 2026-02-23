package config

import "fnchatbot/internal/models"

type ModelDefinition struct {
	ModelID                string                   `json:"model_id"`
	Name                   string                   `json:"name"`
	Group                  string                   `json:"group"`
	Description            string                   `json:"description"`
	OwnedBy                string                   `json:"owned_by"`
	Capabilities           []models.ModelCapability `json:"capabilities"`
	SupportedEndpointTypes []models.EndpointType    `json:"supported_endpoint_types"`
	EndpointType           models.EndpointType      `json:"endpoint_type"`
	MaxTokens              int                      `json:"max_tokens"`
	InputPrice             float64                  `json:"input_price"`
	OutputPrice            float64                  `json:"output_price"`
}

type ProviderModels struct {
	ProviderID string            `json:"provider_id"`
	Models     []ModelDefinition `json:"models"`
}

var SystemModels = map[string][]ModelDefinition{
	"openai": {
		{ModelID: "gpt-5.1", Name: "GPT 5.1", Group: "GPT 5.1", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5", Name: "GPT 5", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5-mini", Name: "GPT 5 Mini", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5-nano", Name: "GPT 5 Nano", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gpt-5-pro", Name: "GPT 5 Pro", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5-chat", Name: "GPT 5 Chat", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gpt-image-1", Name: "GPT Image 1", Group: "GPT Image", Capabilities: []models.ModelCapability{models.CapabilityImageGeneration}},
		{ModelID: "gpt-4o", Name: "GPT-4o", Group: "GPT 4o", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4o-mini", Name: "GPT-4o-mini", Group: "GPT 4o", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "o3", Name: "o3", Group: "o3", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "o4-mini", Name: "o4-mini", Group: "o4", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"anthropic": {
		{ModelID: "claude-opus-4-6", Name: "Claude Opus 4.6", Group: "Claude 4.6", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5", Group: "Claude 4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-haiku-4-5", Name: "Claude Haiku 4.5", Group: "Claude 4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-opus-4-5", Name: "Claude Opus 4.5", Group: "Claude 4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", Group: "Claude 4", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-opus-4-20250514", Name: "Claude Opus 4", Group: "Claude 4", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
	},
	"gemini": {
		{ModelID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Group: "Gemini 2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Group: "Gemini 2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-flash-image-preview", Name: "Gemini 2.5 Flash Image", Group: "Gemini 2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "gemini-3-pro-preview", Name: "Gemini 3 Pro Preview", Group: "Gemini 3", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-3-pro-image-preview", Name: "Gemini 3 Pro Image Preview", Group: "Gemini 3", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
	},
	"deepseek": {
		{ModelID: "deepseek-chat", Name: "DeepSeek Chat", Group: "DeepSeek Chat", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "deepseek-reasoner", Name: "DeepSeek Reasoner", Group: "DeepSeek Reasoner", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"grok": {
		{ModelID: "grok-4", Name: "Grok 4", Group: "Grok", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "grok-3", Name: "Grok 3", Group: "Grok", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "grok-3-fast", Name: "Grok 3 Fast", Group: "Grok", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "grok-3-mini", Name: "Grok 3 Mini", Group: "Grok", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "grok-3-mini-fast", Name: "Grok 3 Mini Fast", Group: "Grok", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"mistral": {
		{ModelID: "pixtral-12b-2409", Name: "Pixtral 12B", Group: "Pixtral", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "pixtral-large-latest", Name: "Pixtral Large", Group: "Pixtral", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "ministral-3b-latest", Name: "Mistral 3B", Group: "Mistral Mini", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "ministral-8b-latest", Name: "Mistral 8B", Group: "Mistral Mini", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "codestral-latest", Name: "Mistral Codestral", Group: "Mistral Code", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "mistral-large-latest", Name: "Mistral Large", Group: "Mistral Chat", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "mistral-small-latest", Name: "Mistral Small", Group: "Mistral Chat", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "open-mistral-nemo", Name: "Mistral Nemo", Group: "Mistral Chat", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "mistral-embed", Name: "Mistral Embedding", Group: "Mistral Embed", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"groq": {
		{ModelID: "llama3-8b-8192", Name: "LLaMA3 8B", Group: "Llama3", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "llama3-70b-8192", Name: "LLaMA3 70B", Group: "Llama3", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "mistral-saba-24b", Name: "Mistral Saba 24B", Group: "Mistral", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gemma-9b-it", Name: "Gemma 9B", Group: "Gemma", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"perplexity": {
		{ModelID: "sonar-reasoning-pro", Name: "sonar-reasoning-pro", Group: "Sonar", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "sonar-reasoning", Name: "sonar-reasoning", Group: "Sonar", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "sonar-pro", Name: "sonar-pro", Group: "Sonar", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityWebSearch}},
		{ModelID: "sonar", Name: "sonar", Group: "Sonar", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityWebSearch}},
		{ModelID: "sonar-deep-research", Name: "sonar-deep-research", Group: "Sonar", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityWebSearch}},
	},
	"nvidia": {
		{ModelID: "01-ai/yi-large", Name: "yi-large", Group: "Yi", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "meta/llama-3.1-405b-instruct", Name: "llama-3.1-405b-instruct", Group: "llama-3.1", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"together": {
		{ModelID: "meta-llama/Llama-3.2-11B-Vision-Instruct-Turbo", Name: "Llama-3.2-11B-Vision", Group: "Llama-3.2", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "meta-llama/Llama-3.2-90B-Vision-Instruct-Turbo", Name: "Llama-3.2-90B-Vision", Group: "Llama-3.2", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "google/gemma-2-27b-it", Name: "gemma-2-27b-it", Group: "Gemma", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "google/gemma-2-9b-it", Name: "gemma-2-9b-it", Group: "Gemma", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"fireworks": {
		{ModelID: "accounts/fireworks/models/mythomax-l2-13b", Name: "mythomax-l2-13b", Group: "Gryphe", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "accounts/fireworks/models/llama-v3-70b-instruct", Name: "Llama-3-70B-Instruct", Group: "Llama3", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"hyperbolic": {
		{ModelID: "Qwen/Qwen2-VL-72B-Instruct", Name: "Qwen2-VL-72B-Instruct", Group: "Qwen2-VL", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "Qwen/Qwen2-VL-7B-Instruct", Name: "Qwen2-VL-7B-Instruct", Group: "Qwen2-VL", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "mistralai/Pixtral-12B-2409", Name: "Pixtral-12B-2409", Group: "Pixtral", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "meta-llama/Meta-Llama-3.1-405B", Name: "Meta-Llama-3.1-405B", Group: "Meta-Llama-3.1", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"cerebras": {
		{ModelID: "gpt-oss-120b", Name: "GPT oss 120B", Group: "openai", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "zai-glm-4.6", Name: "GLM 4.6", Group: "zai", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen-3-235b-a22b-instruct-2507", Name: "Qwen 3 235B A22B Instruct", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"zhipu": {
		{ModelID: "glm-5", Name: "GLM-5", Group: "GLM-5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "glm-4.7", Name: "GLM-4.7", Group: "GLM-4.7", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "glm-4.6", Name: "GLM-4.6", Group: "GLM-4.6", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "glm-4.6v", Name: "GLM-4.6V", Group: "GLM-4.6V", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "glm-4.6v-flash", Name: "GLM-4.6V-Flash", Group: "GLM-4.6V", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "glm-4.5", Name: "GLM-4.5", Group: "GLM-4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "glm-4.5-flash", Name: "GLM-4.5-Flash", Group: "GLM-4.5", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "glm-4.5v", Name: "GLM-4.5V", Group: "GLM-4.5V", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "embedding-3", Name: "Embedding-3", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"moonshot": {
		{ModelID: "moonshot-v1-auto", Name: "moonshot-v1-auto", Group: "moonshot-v1", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "kimi-k2-0711-preview", Name: "kimi-k2-0711-preview", Group: "kimi-k2", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "kimi-k2.5", Name: "Kimi K2.5", Group: "Kimi K2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "kimi-k2-0905-Preview", Name: "Kimi K2 0905 Preview", Group: "Kimi K2", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "kimi-k2-turbo-preview", Name: "Kimi K2 Turbo Preview", Group: "Kimi K2", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "kimi-k2-thinking", Name: "Kimi K2 Thinking", Group: "Kimi K2 Thinking", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"baichuan": {
		{ModelID: "Baichuan4", Name: "Baichuan4", Group: "Baichuan4", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan4-Turbo", Name: "Baichuan4 Turbo", Group: "Baichuan4", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan4-Air", Name: "Baichuan4 Air", Group: "Baichuan4", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan3-Turbo", Name: "Baichuan3 Turbo", Group: "Baichuan3", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan3-Turbo-128k", Name: "Baichuan3 Turbo 128k", Group: "Baichuan3", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan-M2", Name: "Baichuan M2", Group: "Baichuan-M2", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan-M2-Plus", Name: "Baichuan M2 Plus", Group: "Baichuan-M2", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Baichuan-M3", Name: "Baichuan M3", Group: "Baichuan-M3", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"dashscope": {
		{ModelID: "qwen-max", Name: "qwen-max", Group: "qwen-max", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "qwen3-max", Name: "qwen3-max", Group: "qwen-max", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "qwen-plus", Name: "qwen-plus", Group: "qwen-plus", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "qwen3.5-plus", Name: "qwen3.5-plus", Group: "qwen-plus", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "qwen-flash", Name: "qwen-flash", Group: "qwen-flash", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen-vl-plus", Name: "qwen-vl-plus", Group: "qwen-vl", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "qwen-coder-plus", Name: "qwen-coder-plus", Group: "qwen-coder", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "text-embedding-v4", Name: "text-embedding-v4", Group: "qwen-text-embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "text-embedding-v3", Name: "text-embedding-v3", Group: "qwen-text-embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "qwen3-rerank", Name: "qwen3-rerank", Group: "qwen-rerank", Capabilities: []models.ModelCapability{models.CapabilityRerank}},
	},
	"doubao": {
		{ModelID: "doubao-1-5-pro-32k-250115", Name: "doubao-1.5-pro-32k", Group: "Doubao-1.5-pro", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "doubao-1-5-pro-256k-250115", Name: "Doubao-1.5-pro-256k", Group: "Doubao-1.5-pro", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "doubao-1-5-vision-pro-32k-250115", Name: "doubao-1.5-vision-pro", Group: "Doubao-1.5-vision-pro", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "doubao-pro-32k-241215", Name: "Doubao-pro-32k", Group: "Doubao-pro", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "doubao-pro-256k-241115", Name: "Doubao-pro-256k", Group: "Doubao-pro", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "doubao-lite-32k-240828", Name: "Doubao-lite-32k", Group: "Doubao-lite", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "doubao-lite-128k-240828", Name: "Doubao-lite-128k", Group: "Doubao-lite", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "doubao-1-5-lite-32k-250115", Name: "Doubao-1.5-lite-32k", Group: "Doubao-lite", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek-r1-250120", Name: "DeepSeek-R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3-250324", Name: "DeepSeek-V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "doubao-embedding-large-text-240915", Name: "Doubao-embedding-large", Group: "Doubao-embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "doubao-embedding-text-240715", Name: "Doubao-embedding", Group: "Doubao-embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"minimax": {
		{ModelID: "MiniMax-M2.5", Name: "MiniMax-M2.5", Group: "M2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "MiniMax-M2.5-lightning", Name: "MiniMax-M2.5-lightning", Group: "M2.5", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "MiniMax-M2.1", Name: "MiniMax-M2.1", Group: "M2.1", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "MiniMax-M2.1-lightning", Name: "MiniMax-M2.1-lightning", Group: "M2.1", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "MiniMax-M2", Name: "MiniMax-M2", Group: "M2", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"hunyuan": {
		{ModelID: "hunyuan-pro", Name: "hunyuan-pro", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "hunyuan-standard", Name: "hunyuan-standard", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "hunyuan-lite", Name: "hunyuan-lite", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "hunyuan-standard-256k", Name: "hunyuan-standard-256k", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "hunyuan-vision", Name: "hunyuan-vision", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "hunyuan-code", Name: "hunyuan-code", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "hunyuan-turbo", Name: "hunyuan-turbo", Group: "Hunyuan", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "hunyuan-embedding", Name: "hunyuan-embedding", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"baidu-cloud": {
		{ModelID: "deepseek-r1", Name: "DeepSeek R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3", Name: "DeepSeek V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "ernie-4.0-8k-latest", Name: "ERNIE-4.0", Group: "ERNIE", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityFunctionCalling}},
		{ModelID: "ernie-4.0-turbo-8k-latest", Name: "ERNIE 4.0 Turbo", Group: "ERNIE", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "ernie-speed-8k", Name: "ERNIE Speed", Group: "ERNIE", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "ernie-lite-8k", Name: "ERNIE Lite", Group: "ERNIE", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "bge-large-zh", Name: "BGE Large ZH", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "bge-large-en", Name: "BGE Large EN", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"yi": {
		{ModelID: "yi-lightning", Name: "Yi Lightning", Group: "yi-lightning", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "yi-vision-v2", Name: "Yi Vision v2", Group: "yi-vision", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
	},
	"stepfun": {
		{ModelID: "step-1-8k", Name: "Step 1 8K", Group: "Step 1", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "step-1-flash", Name: "Step 1 Flash", Group: "Step 1", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"mimo": {
		{ModelID: "mimo-v2-flash", Name: "Mimo V2 Flash", Group: "Mimo", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"aihubmix": {
		{ModelID: "gpt-5", Name: "gpt-5", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5-mini", Name: "gpt-5-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gpt-5-nano", Name: "gpt-5-nano", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "o3", Name: "o3", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "o4-mini", Name: "o4-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "gpt-4.1", Name: "gpt-4.1", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4o", Name: "gpt-4o", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "DeepSeek-V3", Name: "DeepSeek-V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "DeepSeek-R1", Name: "DeepSeek-R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "claude-sonnet-4-20250514", Name: "claude-sonnet-4-20250514", Group: "Claude", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-pro", Name: "gemini-2.5-pro", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-flash", Name: "gemini-2.5-flash", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
	},
	"silicon": {
		{ModelID: "deepseek-ai/DeepSeek-V3.2", Name: "deepseek-ai/DeepSeek-V3.2", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Qwen/Qwen3-8B", Name: "Qwen/Qwen3-8B", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "BAAI/bge-m3", Name: "BAAI/bge-m3", Group: "BAAI", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"openrouter": {
		{ModelID: "google/gemini-2.5-flash-image-preview", Name: "Google: Gemini 2.5 Flash Image", Group: "google", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "google/gemini-2.5-flash-preview", Name: "Google: Gemini 2.5 Flash Preview", Group: "google", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "qwen/qwen-2.5-7b-instruct:free", Name: "Qwen: Qwen-2.5-7B Instruct", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek/deepseek-chat", Name: "DeepSeek: V3", Group: "deepseek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "mistralai/mistral-7b-instruct:free", Name: "Mistral: Mistral 7B Instruct", Group: "mistralai", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"302ai": {
		{ModelID: "deepseek-chat", Name: "deepseek-chat", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek-reasoner", Name: "deepseek-reasoner", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "chatgpt-4o-latest", Name: "chatgpt-4o-latest", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4.1", Name: "gpt-4.1", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "o3", Name: "o3", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "o4-mini", Name: "o4-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "qwen3-235b-a22b", Name: "qwen3-235b-a22b", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gemini-2.5-flash-preview-05-20", Name: "gemini-2.5-flash-preview-05-20", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "gemini-2.5-pro-preview-06-05", Name: "gemini-2.5-pro-preview-06-05", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-sonnet-4-20250514", Name: "claude-sonnet-4-20250514", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-opus-4-20250514", Name: "claude-opus-4-20250514", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
	},
	"tokenflux": {
		{ModelID: "gpt-4.1", Name: "GPT-4.1", Group: "GPT-4.1", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4.1-mini", Name: "GPT-4.1 Mini", Group: "GPT-4.1", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "claude-sonnet-4", Name: "Claude Sonnet 4", Group: "Claude", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-3-7-sonnet", Name: "Claude 3.7 Sonnet", Group: "Claude", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "deepseek-r1", Name: "DeepSeek R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3", Name: "DeepSeek V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen-max", Name: "Qwen Max", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen-plus", Name: "Qwen Plus", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"github": {
		{ModelID: "gpt-4o", Name: "OpenAI GPT-4o", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
	},
	"copilot": {
		{ModelID: "gpt-4o-mini", Name: "OpenAI GPT-4o-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"jina": {
		{ModelID: "jina-clip-v1", Name: "jina-clip-v1", Group: "Jina Clip", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "jina-clip-v2", Name: "jina-clip-v2", Group: "Jina Clip", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "jina-embeddings-v2-base-en", Name: "jina-embeddings-v2-base-en", Group: "Jina Embeddings V2", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "jina-embeddings-v2-base-zh", Name: "jina-embeddings-v2-base-zh", Group: "Jina Embeddings V2", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "jina-embeddings-v2-base-code", Name: "jina-embeddings-v2-base-code", Group: "Jina Embeddings V2", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "jina-embeddings-v3", Name: "jina-embeddings-v3", Group: "Jina Embeddings V3", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"voyageai": {
		{ModelID: "voyage-3-large", Name: "voyage-3-large", Group: "Voyage Embeddings V3", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "voyage-3", Name: "voyage-3", Group: "Voyage Embeddings V3", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "voyage-3-lite", Name: "voyage-3-lite", Group: "Voyage Embeddings V3", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "voyage-code-3", Name: "voyage-code-3", Group: "Voyage Embeddings V3", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "rerank-2", Name: "rerank-2", Group: "Voyage Rerank V2", Capabilities: []models.ModelCapability{models.CapabilityRerank}},
		{ModelID: "rerank-2-lite", Name: "rerank-2-lite", Group: "Voyage Rerank V2", Capabilities: []models.ModelCapability{models.CapabilityRerank}},
	},
	"ocoolai": {
		{ModelID: "deepseek-chat", Name: "deepseek-chat", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek-reasoner", Name: "deepseek-reasoner", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "gpt-4o", Name: "gpt-4o", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4o-mini", Name: "gpt-4o-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "claude-3-5-sonnet-20240620", Name: "claude-3-5-sonnet-20240620", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-pro", Name: "gemini-pro", Group: "Gemini", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "text-embedding-3-large", Name: "text-embedding-3-large", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "text-embedding-3-small", Name: "text-embedding-3-small", Group: "Embedding", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"ppio": {
		{ModelID: "deepseek/deepseek-v3.2", Name: "DeepSeek V3.2", Group: "deepseek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "minimax/minimax-m2", Name: "MiniMax M2", Group: "minimaxai", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen/qwen3-235b-a22b-instruct-2507", Name: "Qwen3-235b-a22b-instruct-2507", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen/qwen3-vl-235b-a22b-instruct", Name: "Qwen3-vl-235b-a22b-instruct", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "qwen/qwen3-embedding-8b", Name: "Qwen3 Embedding 8B", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
		{ModelID: "qwen/qwen3-reranker-8b", Name: "Qwen3 Reranker 8B", Group: "qwen", Capabilities: []models.ModelCapability{models.CapabilityRerank}},
	},
	"alayanew": {},
	"qiniu": {
		{ModelID: "deepseek-r1", Name: "DeepSeek R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3", Name: "DeepSeek V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwq-32b", Name: "QWQ 32B", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "qwen2.5-72b-instruct", Name: "Qwen2.5 72B Instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"dmxapi": {
		{ModelID: "Qwen/Qwen2.5-7B-Instruct", Name: "Qwen/Qwen2.5-7B-Instruct", Group: "免费模型", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "ERNIE-Speed-128K", Name: "ERNIE-Speed-128K", Group: "免费模型", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "gpt-4o", Name: "gpt-4o", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4o-mini", Name: "gpt-4o-mini", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "DMXAPI-DeepSeek-R1", Name: "DMXAPI-DeepSeek-R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "DMXAPI-DeepSeek-V3", Name: "DMXAPI-DeepSeek-V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"burncloud": {
		{ModelID: "claude-opus-4-5-20251101", Name: "Claude 4.5 Opus", Group: "Claude 4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-sonnet-4-5-20250929", Name: "Claude 4.5 Sonnet", Group: "Claude 4.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5", Name: "GPT 5", Group: "GPT 5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Group: "Gemini 2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Group: "Gemini 2.5", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "deepseek-reasoner", Name: "DeepSeek Reasoner", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-chat", Name: "DeepSeek Chat", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"cephalon": {
		{ModelID: "DeepSeek-R1", Name: "DeepSeek-R1满血版", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"lanyun": {
		{ModelID: "/maas/deepseek-ai/DeepSeek-R1-0528", Name: "deepseek-ai/DeepSeek-R1", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "/maas/deepseek-ai/DeepSeek-V3-0324", Name: "deepseek-ai/DeepSeek-V3", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "/maas/qwen/Qwen2.5-72B-Instruct", Name: "Qwen2.5-72B-Instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "/maas/qwen/Qwen3-235B-A22B", Name: "Qwen/Qwen3-235B", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"ph8": {
		{ModelID: "deepseek-v3-241226", Name: "deepseek-v3-241226", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek-r1-250120", Name: "deepseek-r1-250120", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"sophnet": {},
	"aionly": {
		{ModelID: "claude-opus-4-5-20251101", Name: "Claude Opus 4.5", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-haiku-4-5-20251001", Name: "Claude Haiku 4.5", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "claude-sonnet-4-5-20250929", Name: "Claude Sonnet 4.5", Group: "Anthropic", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-5.1", Name: "GPT-5.1", Group: "OpenAI", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-3-pro-preview", Name: "Gemini 3 Pro Preview", Group: "Google", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Group: "Google", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Group: "Google", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
	},
	"longcat": {
		{ModelID: "LongCat-Flash-Chat", Name: "LongCat Flash Chat", Group: "LongCat", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "LongCat-Flash-Thinking", Name: "LongCat Flash Thinking", Group: "LongCat", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
	},
	"infini": {
		{ModelID: "deepseek-r1", Name: "deepseek-r1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3", Name: "deepseek-v3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "qwen2.5-72b-instruct", Name: "qwen2.5-72b-instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "llama-3.3-70b-instruct", Name: "llama-3.3-70b-instruct", Group: "Llama", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "bge-m3", Name: "bge-m3", Group: "BAAI", Capabilities: []models.ModelCapability{models.CapabilityEmbedding}},
	},
	"modelscope": {
		{ModelID: "Qwen/Qwen2.5-72B-Instruct", Name: "Qwen/Qwen2.5-72B-Instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "Qwen/Qwen2.5-VL-72B-Instruct", Name: "Qwen/Qwen2.5-VL-72B-Instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision}},
		{ModelID: "Qwen/Qwen2.5-Coder-32B-Instruct", Name: "Qwen/Qwen2.5-Coder-32B-Instruct", Group: "Qwen", Capabilities: []models.ModelCapability{models.CapabilityText}},
		{ModelID: "deepseek-ai/DeepSeek-R1", Name: "deepseek-ai/DeepSeek-R1", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-ai/DeepSeek-V3", Name: "deepseek-ai/DeepSeek-V3", Group: "deepseek-ai", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"xirang": {},
	"tencent-cloud-ti": {
		{ModelID: "deepseek-r1", Name: "DeepSeek R1", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityReasoning}},
		{ModelID: "deepseek-v3", Name: "DeepSeek V3", Group: "DeepSeek", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"poe": {
		{ModelID: "gpt-4o", Name: "GPT-4o", Group: "poe", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
	},
	"azure-openai": {
		{ModelID: "gpt-4o", Name: "GPT-4o", Group: "GPT 4o", Capabilities: []models.ModelCapability{models.CapabilityText, models.CapabilityVision, models.CapabilityFunctionCalling}},
		{ModelID: "gpt-4o-mini", Name: "GPT-4o-mini", Group: "GPT 4o", Capabilities: []models.ModelCapability{models.CapabilityText}},
	},
	"vertexai":    {},
	"aws-bedrock": {},
	"ollama":      {},
	"lmstudio":    {},
	"gpustack":    {},
	"ovms":        {},
	"new-api":     {},
	"huggingface": {},
	"gateway":     {},
}

func GetModelsByProviderID(providerID string) []ModelDefinition {
	if models, ok := SystemModels[providerID]; ok {
		return models
	}
	return []ModelDefinition{}
}

func GetAllProviderModels() []ProviderModels {
	var result []ProviderModels
	for providerID, modelList := range SystemModels {
		result = append(result, ProviderModels{
			ProviderID: providerID,
			Models:     modelList,
		})
	}
	return result
}
