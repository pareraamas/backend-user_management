package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pareraamas/user_mgmt/model"
	"github.com/pareraamas/user_mgmt/repository"
)

// Service untuk audit log event penting user
// Komentar Bahasa Indonesia
type AuditLogService struct {
	Repo *repository.AuditLogRepository
}

func NewAuditLogService(repo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{Repo: repo}
}

// Simpan log event user
type AuditLogData map[string]interface{}

func (s *AuditLogService) Log(userID, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	log := &model.AuditLog{
		ID:     uuid.NewString(),
		UserID: userID,
		Event:  event,
		Data:   jsonData,
	}
	_ = s.Repo.Create(log) // abaikan error untuk log async
}
