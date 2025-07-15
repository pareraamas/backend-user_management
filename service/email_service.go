package service

import (
	"fmt"
	"log"
	"github.com/pareraamas/user_mgmt/config"
)

// Service untuk simulasi pengiriman email
// Komentar dalam Bahasa Indonesia

type EmailService struct {
	Config config.SMTPConfig
}

func NewEmailService(cfg config.SMTPConfig) *EmailService {
	return &EmailService{Config: cfg}
}

// Fungsi simulasi kirim email (hanya log)
func (s *EmailService) SendMail(to, subject, body string) error {
	log.Printf("[SIMULASI EMAIL] Ke: %s | Subjek: %s | Isi: %s", to, subject, body)
	fmt.Printf("[SIMULASI EMAIL] Ke: %s | Subjek: %s | Isi: %s\n", to, subject, body)
	return nil
}
