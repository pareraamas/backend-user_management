package repository

import (
	"github.com/pareraamas/user_mgmt/model"
	"gorm.io/gorm"
)

// Repository untuk audit log user
// Komentar Bahasa Indonesia
type AuditLogRepository struct {
	DB *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.DB.Create(log).Error
}

// (Opsional) Query log berdasarkan user/event
func (r *AuditLogRepository) FindByUser(userID string, limit int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&logs).Error
	return logs, err
}
