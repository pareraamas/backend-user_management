package model

import "gorm.io/datatypes"

// Model untuk audit log event penting user
// Komentar dalam Bahasa Indonesia
type AuditLog struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"index" json:"user_id"`
	Event     string         `json:"event"`
	Data      datatypes.JSON `json:"data" gorm:"type:json"`
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
}
