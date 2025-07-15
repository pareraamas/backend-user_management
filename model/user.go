package model

import "gorm.io/datatypes"

// Struktur User diperluas agar mirip ORY Kratos
// Komentar dalam Bahasa Indonesia

type User struct {
	ID                string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	Password          string         `gorm:"not null" json:"-"`
	Status            string         `gorm:"default:'pending'" json:"status"` // active, pending, banned
	IsEmailVerified   bool           `gorm:"default:false" json:"is_email_verified"`
	VerificationCode  string         `gorm:"size:6" json:"-"`
	RecoveryCode      string         `gorm:"size:6" json:"-"`
	RecoveryRequested int64          `json:"recovery_requested"`
	LastLogin         int64          `json:"last_login"`
	CreatedAt         int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	// Field untuk custom identity schema (atribut dinamis)
	CustomAttributes  datatypes.JSON `json:"custom_attributes" gorm:"type:json"`
	// Field untuk TOTP (2FA)
	TOTPSecret        string         `json:"-"`
	TOTPEnabled       bool           `json:"totp_enabled"`
}
