package models

import (
	"time"

	"gorm.io/gorm"
)

type SandboxConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SandboxPath struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Path            string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"path"`
	Description     string    `gorm:"type:text" json:"description"`
	SandboxConfigID uint      `gorm:"index;not null" json:"sandbox_config_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func GetSandboxConfig(db *gorm.DB) (*SandboxConfig, error) {
	var config SandboxConfig
	result := db.First(&config)
	if result.Error == gorm.ErrRecordNotFound {
		config = SandboxConfig{
			Enabled: true,
		}
		if err := db.Create(&config).Error; err != nil {
			return nil, err
		}
		return &config, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}
