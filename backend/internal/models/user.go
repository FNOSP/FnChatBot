package models

import "time"

// User represents an application user with authentication and role info.
type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Username           string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash       string    `gorm:"type:varchar(255);not null" json:"-"`
	Description        string    `gorm:"type:text" json:"description"`
	IsAdmin            bool      `gorm:"default:false" json:"is_admin"`
	Enabled            bool      `gorm:"default:true;index" json:"enabled"`
	MustChangePassword bool      `gorm:"default:false" json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
