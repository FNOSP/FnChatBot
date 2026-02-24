package models

import (
	"time"

	"gorm.io/datatypes"
)

// ModelConfig represents the AI model configuration
type ModelConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Provider    string    `gorm:"type:varchar(50);not null;index" json:"provider"` // openai, anthropic, google, custom
	ProviderID  uint      `gorm:"index" json:"provider_id"`
	ProviderRef *Provider `gorm:"foreignKey:ProviderID" json:"provider_ref,omitempty"`
	Model       string    `gorm:"type:varchar(100);not null" json:"model"`
	Temperature float32   `gorm:"default:0.7" json:"temperature"`
	MaxTokens   int       `gorm:"default:2048" json:"max_tokens"`
	IsDefault   bool      `gorm:"default:false;index" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Conversation represents a chat session
type Conversation struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Title     string      `gorm:"type:varchar(200);not null" json:"title"`
	ModelID   uint        `json:"model_id"`
	Model     ModelConfig `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	CreatedAt time.Time   `gorm:"index" json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// MessageRole defines the role of the message sender
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// Message represents a single message in a conversation
type Message struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID uint           `gorm:"index;not null" json:"conversation_id"`
	Role           MessageRole    `gorm:"type:varchar(20);not null" json:"role"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Meta           datatypes.JSON `json:"meta"` // Stores thinking process, task updates, etc.
	CreatedAt      time.Time      `gorm:"index" json:"created_at"`
}

// MCPConfig and Skill moved to separate files

// AgentTask represents a task executed by the agent
type AgentTask struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID uint           `gorm:"index;not null" json:"conversation_id"`
	SkillName      string         `gorm:"type:varchar(100);not null" json:"skill_name"`
	Status         string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"` // pending, running, completed, failed
	Parameters     datatypes.JSON `json:"parameters"`
	Result         datatypes.JSON `json:"result"`
	ErrorMessage   string         `gorm:"type:text" json:"error_message"`
	CreatedAt      time.Time      `json:"created_at"`
	StartedAt      *time.Time     `json:"started_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
}

type SandboxPathInfo struct {
	ID          uint   `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
