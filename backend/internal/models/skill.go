package models

import (
	"time"

	"gorm.io/datatypes"
)

// Skill represents a capability that the agent can use
type Skill struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	UserID      uint           `gorm:"index" json:"user_id"`
	Priority    int            `gorm:"default:0" json:"priority"`
	Config      datatypes.JSON `json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
}
