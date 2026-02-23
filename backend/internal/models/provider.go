package models

import (
	"time"

	"gorm.io/gorm"
)

// ProviderType 定义提供商类型
type ProviderType string

const (
	ProviderTypeOpenAI          ProviderType = "openai"
	ProviderTypeOpenAIResponse  ProviderType = "openai-response"
	ProviderTypeAnthropic       ProviderType = "anthropic"
	ProviderTypeGemini          ProviderType = "gemini"
	ProviderTypeAzureOpenAI     ProviderType = "azure-openai"
	ProviderTypeVertexAI        ProviderType = "vertexai"
	ProviderTypeMistral         ProviderType = "mistral"
	ProviderTypeAWBedrock       ProviderType = "aws-bedrock"
	ProviderTypeVertexAnthropic ProviderType = "vertex-anthropic"
	ProviderTypeNewAPI          ProviderType = "new-api"
	ProviderTypeGateway         ProviderType = "gateway"
	ProviderTypeOllama          ProviderType = "ollama"
)

// ModelCapability 定义模型能力标签类型
type ModelCapability string

const (
	// CapabilityText 文本生成能力
	CapabilityText ModelCapability = "text"
	// CapabilityVision 视觉理解能力
	CapabilityVision ModelCapability = "vision"
	// CapabilityEmbedding 向量嵌入能力
	CapabilityEmbedding ModelCapability = "embedding"
	// CapabilityReasoning 推理能力
	CapabilityReasoning ModelCapability = "reasoning"
	// CapabilityFunctionCalling 函数调用能力
	CapabilityFunctionCalling ModelCapability = "function_calling"
	// CapabilityWebSearch 网络搜索能力
	CapabilityWebSearch ModelCapability = "web_search"
	// CapabilityRerank 重排序能力
	CapabilityRerank ModelCapability = "rerank"
	// CapabilityImageGeneration 图像生成能力
	CapabilityImageGeneration ModelCapability = "image_generation"
)

// EndpointType 定义端点类型
type EndpointType string

const (
	// EndpointTypeOpenAI OpenAI 兼容端点
	EndpointTypeOpenAI EndpointType = "openai"
	// EndpointTypeOpenAIResponse OpenAI 响应格式端点
	EndpointTypeOpenAIResponse EndpointType = "openai-response"
	// EndpointTypeAnthropic Anthropic 端点
	EndpointTypeAnthropic EndpointType = "anthropic"
	// EndpointTypeGemini Gemini 端点
	EndpointTypeGemini EndpointType = "gemini"
	// EndpointTypeImageGeneration 图像生成端点
	EndpointTypeImageGeneration EndpointType = "image-generation"
	// EndpointTypeJinaRerank Jina 重排序端点
	EndpointTypeJinaRerank EndpointType = "jina-rerank"
)

// ProviderApiOptions 定义提供商 API 兼容性选项
type ProviderApiOptions struct {
	IsNotSupportArrayContent   bool `json:"is_not_support_array_content"`
	IsNotSupportStreamOptions  bool `json:"is_not_support_stream_options"`
	IsSupportDeveloperRole     bool `json:"is_support_developer_role"`
	IsSupportServiceTier       bool `json:"is_support_service_tier"`
	IsNotSupportEnableThinking bool `json:"is_not_support_enable_thinking"`
	IsNotSupportAPIVersion     bool `json:"is_not_support_api_version"`
	IsNotSupportVerbosity      bool `json:"is_not_support_verbosity"`
}

// Provider 定义 AI 服务提供商
type Provider struct {
	ID         uint               `gorm:"primaryKey" json:"id"`
	ProviderID string             `gorm:"type:varchar(50);uniqueIndex;not null" json:"provider_id"`
	Name       string             `gorm:"type:varchar(100);not null" json:"name"`
	Type       ProviderType       `gorm:"type:varchar(50);not null;index" json:"type"`
	BaseURL    string             `gorm:"type:varchar(500);not null" json:"base_url"`
	APIKey     string             `gorm:"type:varchar(500)" json:"api_key,omitempty"`
	Enabled    bool               `gorm:"default:false;index" json:"enabled"`
	IsSystem   bool               `gorm:"default:false" json:"is_system"`
	ApiOptions ProviderApiOptions `gorm:"type:text;serializer:json" json:"api_options"`
	Models     []Model            `gorm:"foreignKey:ProviderID;references:ID" json:"models,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// Model 定义 AI 模型信息
type Model struct {
	ID                     uint              `gorm:"primaryKey" json:"id"`
	ProviderID             uint              `gorm:"index;not null" json:"provider_id"`
	Provider               *Provider         `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	ModelID                string            `gorm:"type:varchar(100);not null;index" json:"model_id"`
	Name                   string            `gorm:"type:varchar(100);not null" json:"name"`
	Group                  string            `gorm:"type:varchar(100)" json:"group"`
	Description            string            `gorm:"type:text" json:"description"`
	OwnedBy                string            `gorm:"type:varchar(100)" json:"owned_by"`
	Capabilities           []ModelCapability `gorm:"type:text;serializer:json" json:"capabilities"`
	SupportedEndpointTypes []EndpointType    `gorm:"type:text;serializer:json" json:"supported_endpoint_types"`
	EndpointType           EndpointType      `gorm:"type:varchar(50)" json:"endpoint_type"`
	MaxTokens              int               `gorm:"default:4096" json:"max_tokens"`
	InputPrice             float64           `gorm:"type:decimal(10,6);default:0" json:"input_price"`
	OutputPrice            float64           `gorm:"type:decimal(10,6);default:0" json:"output_price"`
	SupportedTextDelta     bool              `gorm:"default:true" json:"supported_text_delta"`
	Enabled                bool              `gorm:"default:true;index" json:"enabled"`
	IsDefault              bool              `gorm:"default:false" json:"is_default"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
}

// TableName 指定 Provider 表名
func (Provider) TableName() string {
	return "providers"
}

// TableName 指定 Model 表名
func (Model) TableName() string {
	return "models"
}

// BeforeCreate GORM 钩子：创建前自动设置时间
func (p *Provider) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前自动更新时间
func (p *Provider) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate GORM 钩子：创建前自动设置时间
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前自动更新时间
func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// HasCapability 检查模型是否具有指定能力
func (m *Model) HasCapability(capability ModelCapability) bool {
	for _, c := range m.Capabilities {
		if c == capability {
			return true
		}
	}
	return false
}

// SupportsEndpointType 检查模型是否支持指定端点类型
func (m *Model) SupportsEndpointType(endpointType EndpointType) bool {
	for _, et := range m.SupportedEndpointTypes {
		if et == endpointType {
			return true
		}
	}
	return false
}
