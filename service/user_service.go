package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/pareraamas/user_mgmt/model"
	"github.com/pareraamas/user_mgmt/repository"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// Service utama untuk user management
// Komentar dalam Bahasa Indonesia

type UserService struct {
	Repo *repository.UserRepository
	AuditLog *AuditLogService
	WebhookURLRegister    string
	WebhookURLLogin       string
	WebhookURLRecovery    string
	WebhookURLVerification string
	// Map in-memory untuk flow_id TOTP
	totpFlow map[string]string
}

// Tambahkan inisialisasi map pada NewUserService


// Simpan flow_id <-> user_id untuk login TOTP
func (s *UserService) SaveTOTPFlow(flowID, userID string) {
	s.totpFlow[flowID] = userID
}

// Ambil user_id dari flow_id
func (s *UserService) GetUserIDByTOTPFlow(flowID string) (string, error) {
	userID, ok := s.totpFlow[flowID]
	if !ok {
		return "", errors.New("flow_id tidak ditemukan")
	}
	delete(s.totpFlow, flowID) // hapus setelah sekali pakai
	return userID, nil
}

// Ambil user berdasarkan ID
func (s *UserService) GetByID(userID string) (*model.User, error) {
	return s.Repo.FindByID(userID)
}

// Update user
func (s *UserService) Update(user *model.User) error {
	return s.Repo.Update(user)
}

// Fungsi untuk trigger webhook/callback pada event penting user
// Semua event utama kini sudah otomatis trigger webhook jika URL di-set via ENV.
// Contoh ENV:
//   WEBHOOK_URL_REGISTER, WEBHOOK_URL_LOGIN, WEBHOOK_URL_RECOVERY, WEBHOOK_URL_VERIFICATION
// Payload dikirim secara async POST JSON: {event, data, timestamp}
//
// Dokumentasi lengkap di API_DOKUMENTASI.md

// Fungsi untuk trigger webhook/callback pada event penting user
// event: nama event (misal: register, login, recovery_success, dsb)
// payload: data event (bisa user, dsb)
// url: endpoint webhook/callback
func (s *UserService) triggerWebhookEvent(event string, payload interface{}, url string) {
	body := map[string]interface{}{
		"event": event,
		"data":  payload,
		"timestamp": time.Now().UnixMilli(),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return // log error jika perlu
	}
	go func() {
		_, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	}()
}

func NewUserService(repo *repository.UserRepository, auditLog *AuditLogService, webhookRegister, webhookLogin, webhookRecovery, webhookVerification string) *UserService {
	return &UserService{
		Repo: repo,
		AuditLog: auditLog,
		WebhookURLRegister: webhookRegister,
		WebhookURLLogin: webhookLogin,
		WebhookURLRecovery: webhookRecovery,
		WebhookURLVerification: webhookVerification,
		totpFlow: make(map[string]string),
	}
}

func (s *UserService) Register(email, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
		Status:   "pending",
	}
	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}
	// Audit log event register
	if s.AuditLog != nil {
		s.AuditLog.Log(user.ID, "register", map[string]interface{}{"email": user.Email})
	}
	// Trigger webhook jika URL di-set
	if s.WebhookURLRegister != "" {
		s.triggerWebhookEvent("register", user, s.WebhookURLRegister)
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	// Audit log event login
	if s.AuditLog != nil {
		s.AuditLog.Log(user.ID, "login", map[string]interface{}{"email": user.Email})
	}
	// Trigger webhook jika URL di-set
	if s.WebhookURLLogin != "" {
		s.triggerWebhookEvent("login", user, s.WebhookURLLogin)
	}
	return user, nil
}

// Generate dan simpan kode recovery password
func (s *UserService) GenerateAndSaveRecoveryCode(email string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("Email tidak ditemukan")
	}
	code := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(900000) + 100000) // 6 digit
	user.RecoveryCode = code
	user.RecoveryRequested = time.Now().UnixMilli()
	user.UpdatedAt = time.Now().UnixMilli()
	if err := s.Repo.Update(user); err != nil {
		return "", errors.New("Gagal menyimpan kode recovery")
	}
	return code, nil
}

// Update profile user (email dan status)
func (s *UserService) UpdateProfile(userID, email, status string) (*model.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}
	if email != "" {
		user.Email = email
	}
	if status != "" {
		user.Status = status
	}
	user.UpdatedAt = time.Now().UnixMilli()
	if err := s.Repo.Update(user); err != nil {
		return nil, errors.New("Gagal update profil")
	}
	return user, nil
}

// Reset password dengan kode recovery
func (s *UserService) ResetPasswordWithCode(email, code, newPassword string) (bool, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return false, errors.New("Email tidak ditemukan")
	}
	if user.RecoveryCode != code {
		return false, nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	user.Password = string(hash)
	user.RecoveryCode = ""
	user.UpdatedAt = time.Now().UnixMilli()
	if err := s.Repo.Update(user); err != nil {
		return false, errors.New("Gagal update password")
	}
	return true, nil
}

// Generate dan simpan kode verifikasi email
func (s *UserService) GenerateAndSaveVerificationCode(email string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("Email tidak ditemukan")
	}
	code := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(900000) + 100000) // 6 digit
	user.VerificationCode = code
	user.UpdatedAt = time.Now().UnixMilli()
	if err := s.Repo.Update(user); err != nil {
		return "", errors.New("Gagal menyimpan kode verifikasi")
	}
	return code, nil
}

// Verifikasi kode email
func (s *UserService) VerifyEmail(email, code string) error {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return errors.New("Email tidak ditemukan")
	}
	if user.VerificationCode != code {
		return errors.New("Kode verifikasi salah atau kadaluarsa")
	}
	user.IsEmailVerified = true
	user.VerificationCode = ""
	user.Status = "active"
	if err := s.Repo.Update(user); err != nil {
		return err
	}
	// Audit log event verifikasi email
	if s.AuditLog != nil {
		s.AuditLog.Log(user.ID, "verification_success", map[string]interface{}{"email": user.Email})
	}
	// Trigger webhook jika URL di-set
	if s.WebhookURLVerification != "" {
		s.triggerWebhookEvent("verification_success", user, s.WebhookURLVerification)
	}
	return nil
}
