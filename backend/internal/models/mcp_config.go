package models

import (
	"time"
)

// MCPConfig represents the Model Context Protocol configuration
type MCPConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	BaseURL   string    `gorm:"type:varchar(500);not null" json:"base_url"`
	ApiKey    string    `gorm:"type:varchar(500)" json:"api_key"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
