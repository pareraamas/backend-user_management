package migration

import (
	"gorm.io/gorm"
)

// Migrasi tabel user sesuai kebutuhan ORY Kratos (sederhana)
func CreateUsersTable(db *gorm.DB) error {
	type User struct {
		ID                string `gorm:"primaryKey;type:varchar(36)"`
		Email             string `gorm:"uniqueIndex;not null"`
		Password          string `gorm:"not null"`
		Status            string `gorm:"default:'pending'"`
		IsEmailVerified   bool   `gorm:"default:false"`
		VerificationCode  string `gorm:"size:6"`
		RecoveryCode      string `gorm:"size:6"`
		RecoveryRequested int64  
		LastLogin         int64  
		CreatedAt         int64  `gorm:"autoCreateTime:milli"`
		UpdatedAt         int64  `gorm:"autoUpdateTime:milli"`
	}
	return db.AutoMigrate(&User{})
}
